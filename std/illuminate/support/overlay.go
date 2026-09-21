package support

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// VendorOverlay 先用 Go 实现热路径方法，再把 vendor 官方 PHP 类挂到同名上。
// 未覆盖的方法走 PHP（含 Macroable），避免半成品 AddClass 挡住生态。
type VendorOverlay struct {
	node.Node
	name       string
	vendorRel  string
	implements []string
	props      []data.Property
	native     map[string]data.Method
	ctor       data.Method
	vm         data.VM

	mu  sync.Mutex
	php data.ClassStmt
}

func newVendorOverlay(name, vendorRel string, implements []string, props []data.Property, native map[string]data.Method, ctor data.Method) *VendorOverlay {
	if native == nil {
		native = map[string]data.Method{}
	}
	return &VendorOverlay{
		name:       name,
		vendorRel:  vendorRel,
		implements: implements,
		props:      props,
		native:     native,
		ctor:       ctor,
	}
}

func (o *VendorOverlay) SetVM(vm data.VM) { o.vm = vm }

func (o *VendorOverlay) AttachVendorClass(c data.ClassStmt) {
	if c == nil || c == o {
		return
	}
	o.mu.Lock()
	if o.php == nil {
		o.php = c
	}
	o.mu.Unlock()
}

func (o *VendorOverlay) VendorClass() data.ClassStmt {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.php
}

func (o *VendorOverlay) GetName() string { return o.name }
func (o *VendorOverlay) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(o, ctx.CreateBaseContext()), nil
}

func (o *VendorOverlay) GetExtend() *string {
	if php := o.VendorClass(); php != nil {
		return php.GetExtend()
	}
	return nil
}

func (o *VendorOverlay) GetImplements() []string {
	if php := o.VendorClass(); php != nil {
		if impl := php.GetImplements(); len(impl) > 0 {
			return impl
		}
	}
	return o.implements
}

func (o *VendorOverlay) GetProperty(name string) (data.Property, bool) {
	for _, p := range o.props {
		if p != nil && p.GetName() == name {
			return p, true
		}
	}
	if php := o.VendorClass(); php != nil {
		return php.GetProperty(name)
	}
	return nil, false
}

func (o *VendorOverlay) GetPropertyList() []data.Property {
	if php := o.VendorClass(); php != nil {
		list := php.GetPropertyList()
		if len(list) > 0 {
			return list
		}
	}
	return o.props
}

func (o *VendorOverlay) GetConstruct() data.Method {
	if o.ctor != nil {
		return o.ctor
	}
	_ = o.ensurePHP(o.vm)
	if php := o.VendorClass(); php != nil {
		return php.GetConstruct()
	}
	return nil
}

func (o *VendorOverlay) GetMethod(name string) (data.Method, bool) {
	if m, ok := o.native[strings.ToLower(name)]; ok {
		return m, true
	}
	_ = o.ensurePHP(o.vm)
	if php := o.VendorClass(); php != nil {
		return php.GetMethod(name)
	}
	if strings.EqualFold(name, "__call") {
		return &overlayMagicMethod{overlay: o, name: "__call", static: false}, true
	}
	return nil, false
}

func (o *VendorOverlay) GetMethods() []data.Method {
	seen := map[string]struct{}{}
	out := make([]data.Method, 0, len(o.native)+8)
	for k, m := range o.native {
		seen[k] = struct{}{}
		out = append(out, m)
	}
	if php := o.VendorClass(); php != nil {
		for _, m := range php.GetMethods() {
			key := strings.ToLower(m.GetName())
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, m)
		}
	}
	return out
}

func (o *VendorOverlay) GetStaticMethod(name string) (data.Method, bool) {
	if m, ok := o.native[strings.ToLower(name)]; ok {
		return m, true
	}
	_ = o.ensurePHP(o.vm)
	if php := o.VendorClass(); php != nil {
		if gsm, ok := php.(data.GetStaticMethod); ok {
			if m, ok := gsm.GetStaticMethod(name); ok {
				return m, true
			}
		}
		if m, ok := php.GetMethod(name); ok {
			return m, true
		}
		for _, m := range php.GetMethods() {
			if m != nil && strings.EqualFold(m.GetName(), name) {
				return m, true
			}
		}
	}
	if strings.EqualFold(name, "__callStatic") {
		return &overlayMagicMethod{overlay: o, name: "__callStatic", static: true}, true
	}
	return nil, false
}

func (o *VendorOverlay) GetStaticProperty(name string) (data.Value, bool) {
	if php := o.VendorClass(); php != nil {
		if gsp, ok := php.(data.GetStaticProperty); ok {
			return gsp.GetStaticProperty(name)
		}
	}
	return nil, false
}

func (o *VendorOverlay) ensurePHP(vm data.VM) data.Control {
	if o.VendorClass() != nil || vm == nil || o.vendorRel == "" {
		return nil
	}
	path := findVendorRel(o.vendorRel)
	if path == "" {
		return nil
	}
	_, ctl := vm.LoadAndRun(path)
	return ctl
}

func findVendorRel(rel string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	rel = filepath.FromSlash(rel)
	for d := wd; ; d = filepath.Dir(d) {
		p := filepath.Join(d, rel)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
	}
	return ""
}

type overlayMagicMethod struct {
	overlay *VendorOverlay
	name    string
	static  bool
}

func (m *overlayMagicMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if ctl := m.overlay.ensurePHP(ctx.GetVM()); ctl != nil {
		return nil, ctl
	}
	if m.name == "__callStatic" || m.name == "__call" {
		if nameVal, ok := ctx.GetIndexValue(0); ok && nameVal != nil {
			orig := nameVal.AsString()
			params, _ := ctx.GetIndexValue(1)
			var fn data.Method
			var ok bool
			if m.static {
				fn, ok = m.overlay.GetStaticMethod(orig)
			} else {
				fn, ok = m.overlay.GetMethod(orig)
			}
			if ok && fn != nil && fn != m {
				if _, magic := fn.(*overlayMagicMethod); !magic {
					return callOverlayMethod(ctx, fn, params, !m.static)
				}
			}
		}
	}
	php := m.overlay.VendorClass()
	if php == nil {
		return data.NewNullValue(), nil
	}
	if m.static {
		if gsm, ok := php.(data.GetStaticMethod); ok {
			if fn, ok := gsm.GetStaticMethod(m.name); ok {
				return fn.Call(ctx)
			}
		}
	}
	if fn, ok := php.GetMethod(m.name); ok {
		return fn.Call(ctx)
	}
	return data.NewNullValue(), nil
}

func callOverlayMethod(ctx data.Context, fn data.Method, params data.Value, inst bool) (data.GetValue, data.Control) {
	vars := fn.GetVariables()
	var nctx data.Context
	if inst {
		if cv := overlayReceiver(ctx); cv != nil {
			nctx = cv.CreateContext(vars)
		}
	}
	if nctx == nil {
		if vm := ctx.GetVM(); vm != nil {
			nctx = vm.CreateContext(vars)
		} else {
			return fn.Call(ctx)
		}
	}
	i := 0
	for _, e := range toEntries(params) {
		if i >= len(vars) {
			break
		}
		if ctl := vars[i].SetValue(nctx, e.value); ctl != nil {
			return nil, ctl
		}
		i++
	}
	return fn.Call(nctx)
}

func (m *overlayMagicMethod) GetName() string { return m.name }
func (m *overlayMagicMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *overlayMagicMethod) GetIsStatic() bool        { return m.static }
func (m *overlayMagicMethod) GetParams() []data.GetValue {
	if m.name == "__call" || m.name == "__callStatic" {
		return []data.GetValue{
			node.NewParameter(nil, "method", 0, nil, nil),
			node.NewParameter(nil, "parameters", 1, nil, nil),
		}
	}
	return nil
}
func (m *overlayMagicMethod) GetVariables() []data.Variable {
	if m.name == "__call" || m.name == "__callStatic" {
		return []data.Variable{
			node.NewVariable(nil, "method", 0, nil),
			node.NewVariable(nil, "parameters", 1, nil),
		}
	}
	return nil
}
func (m *overlayMagicMethod) GetReturnType() data.Types { return nil }

func overlayReceiver(ctx data.Context) *data.ClassValue {
	if methodCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return methodCtx.ClassValue
	}
	if value, ok := ctx.(*data.ClassValue); ok {
		return value
	}
	return nil
}
