package view

import (
	"fmt"
	"html"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const attributeBagClassName = "Illuminate\\View\\ComponentAttributeBag"
const appendableClassName = "Illuminate\\View\\AppendableAttributeValue"

// AttributeBagClass 对齐 ComponentAttributeBag 热路径。
type AttributeBagClass struct {
	node.Node
	methods map[string]data.Method
}

func NewAttributeBagClass() data.ClassStmt {
	c := &AttributeBagClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *AttributeBagClass) GetName() string    { return attributeBagClassName }
func (c *AttributeBagClass) GetExtend() *string { return nil }
func (c *AttributeBagClass) GetImplements() []string {
	return []string{"ArrayAccess", "IteratorAggregate", "Illuminate\\Contracts\\Support\\Htmlable", "Illuminate\\Contracts\\Support\\Arrayable"}
}
func (c *AttributeBagClass) GetProperty(name string) (data.Property, bool) {
	if name == "attributes" {
		return node.NewProperty(nil, "attributes", "protected", false, data.NewArrayValue(nil)), true
	}
	return nil, false
}
func (c *AttributeBagClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "attributes", "protected", false, data.NewArrayValue(nil))}
}
func (c *AttributeBagClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *AttributeBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *AttributeBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *AttributeBagClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *AttributeBagClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *AttributeBagClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"attributes"}, 0, bagConstruct)
	c.methods["get"] = kit.InstanceMethodOpt("get", []string{"key", "default"}, 1, bagGet)
	c.methods["has"] = kit.InstanceMethod("has", []string{"key"}, bagHas)
	c.methods["missing"] = kit.InstanceMethod("missing", []string{"key"}, bagMissing)
	c.methods["all"] = kit.InstanceMethod("all", nil, bagAll)
	c.methods["getattributes"] = kit.InstanceMethod("getAttributes", nil, bagAll)
	c.methods["setattributes"] = kit.InstanceMethod("setAttributes", []string{"attributes"}, bagSetAttributes)
	c.methods["merge"] = kit.InstanceMethodOpt("merge", []string{"attributeDefaults", "escape"}, 0, bagMerge)
	c.methods["class"] = kit.InstanceMethod("class", []string{"classList"}, bagClass)
	c.methods["style"] = kit.InstanceMethod("style", []string{"styleList"}, bagStyle)
	c.methods["only"] = kit.InstanceMethod("only", []string{"keys"}, bagOnly)
	c.methods["except"] = kit.InstanceMethod("except", []string{"keys"}, bagExcept)
	c.methods["exceptprops"] = kit.InstanceMethod("exceptProps", []string{"props"}, bagExceptProps)
	c.methods["onlyprops"] = kit.InstanceMethod("onlyProps", []string{"props"}, bagOnlyProps)
	c.methods["wherestartswith"] = kit.InstanceMethod("whereStartsWith", []string{"needles"}, bagWhereStartsWith)
	c.methods["wheredoesntstartswith"] = kit.InstanceMethod("whereDoesntStartWith", []string{"needles"}, bagWhereDoesntStartWith)
	c.methods["__tostring"] = kit.InstanceMethod("__toString", nil, bagToHtml)
	c.methods["tohtml"] = kit.InstanceMethod("toHtml", nil, bagToHtml)
	c.methods["toarray"] = kit.InstanceMethod("toArray", nil, bagAll)
	c.methods["getiterator"] = kit.InstanceMethod("getIterator", nil, bagGetIterator)
	c.methods["offsetexists"] = kit.InstanceMethod("offsetExists", []string{"offset"}, bagHas)
	c.methods["offsetget"] = kit.InstanceMethod("offsetGet", []string{"offset"}, bagOffsetGet)
	c.methods["offsetset"] = kit.InstanceMethod("offsetSet", []string{"offset", "value"}, bagOffsetSet)
	c.methods["offsetunset"] = kit.InstanceMethod("offsetUnset", []string{"offset"}, bagOffsetUnset)
	c.methods["shouldescapeattributevalue"] = kit.InstanceMethod("shouldEscapeAttributeValue", []string{"key", "value"}, bagShouldEscape)
	kit.RegisterMacroable(c.methods, attributeBagClassName)
}

func bagRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("ComponentAttributeBag missing $this"))
}

func bagAttrs(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("attributes")
	if av, ok := kit.Unwrap(v).(*data.ArrayValue); ok && av != nil {
		return av
	}
	empty := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("attributes", empty)
	return empty
}

func bagConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	attrs := kit.Arg(ctx, 0)
	if attrs == nil || kit.IsNull(attrs) {
		attrs = data.NewArrayValue(nil)
	}
	if av, ok := kit.Unwrap(attrs).(*data.ArrayValue); ok {
		_ = cv.SetProperty("attributes", av)
	} else {
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, e := range kit.Entries(attrs) {
			zv := data.NewZVal(e.Value)
			zv.Name = e.KeyStr
			out.List = append(out.List, zv)
		}
		_ = cv.SetProperty("attributes", out)
	}
	return cv, nil
}

func newBag(ctx data.Context, attrs *data.ArrayValue) *data.ClassValue {
	vm := ctx.GetVM()
	cls, _ := vm.GetClass(attributeBagClassName)
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	_ = cv.SetProperty("attributes", attrs)
	return cv
}

func bagGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if e.KeyStr == key {
			return e.Value, nil
		}
	}
	def := kit.Arg(ctx, 1)
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func bagHas(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if e.KeyStr == key {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func bagMissing(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := bagHas(ctx)
	if ctl != nil {
		return nil, ctl
	}
	b, _ := v.(data.AsBool).AsBool()
	return data.NewBoolValue(!b), nil
}

func bagAll(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return bagAttrs(cv), nil
}

func bagSetAttributes(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	attrs := kit.Arg(ctx, 0)
	if attrs == nil || kit.IsNull(attrs) {
		attrs = data.NewArrayValue(nil)
	}
	if av, ok := kit.Unwrap(attrs).(*data.ArrayValue); ok {
		_ = cv.SetProperty("attributes", av)
	} else {
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, e := range kit.Entries(attrs) {
			zv := data.NewZVal(e.Value)
			zv.Name = e.KeyStr
			out.List = append(out.List, zv)
		}
		_ = cv.SetProperty("attributes", out)
	}
	return cv, nil
}

func bagMerge(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	defaults := kit.Arg(ctx, 0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	seen := map[string]bool{}
	for _, e := range kit.Entries(bagAttrs(cv)) {
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
		seen[e.KeyStr] = true
	}
	for _, e := range kit.Entries(defaults) {
		if seen[e.KeyStr] {
			continue
		}
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
	}
	return newBag(ctx, out), nil
}

func bagClass(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	list := classListString(kit.Arg(ctx, 0))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(bagAttrs(cv)) {
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
	}
	if list != "" {
		existing := ""
		for _, e := range kit.Entries(out) {
			if e.KeyStr == "class" {
				existing = e.Value.AsString()
			}
		}
		merged := strings.TrimSpace(existing + " " + list)
		found := false
		for i, z := range out.List {
			if z != nil && z.Name == "class" {
				out.List[i] = data.NewZVal(data.NewStringValue(merged))
				out.List[i].Name = "class"
				found = true
				break
			}
		}
		if !found {
			zv := data.NewZVal(data.NewStringValue(merged))
			zv.Name = "class"
			out.List = append(out.List, zv)
		}
	}
	return newBag(ctx, out), nil
}

func classListString(v data.Value) string {
	v = kit.Unwrap(v)
	if v == nil {
		return ""
	}
	if s, ok := v.(*data.StringValue); ok {
		return s.AsString()
	}
	parts := []string{}
	for _, e := range kit.Entries(v) {
		if b, ok := e.Value.(data.AsBool); ok {
			if okv, err := b.AsBool(); err == nil && okv {
				parts = append(parts, e.KeyStr)
			}
			continue
		}
		if kit.Truthy(e.Value) {
			if e.KeyStr != "" {
				parts = append(parts, e.KeyStr)
			} else {
				parts = append(parts, e.Value.AsString())
			}
		}
	}
	return strings.Join(parts, " ")
}

func bagStyle(ctx data.Context) (data.GetValue, data.Control) {
	// 简化：同 class，写入 style 属性
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	style := kit.Arg(ctx, 0).AsString()
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(bagAttrs(cv)) {
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
	}
	zv := data.NewZVal(data.NewStringValue(style))
	zv.Name = "style"
	out.List = append(out.List, zv)
	return newBag(ctx, out), nil
}

func bagFilterKeys(ctx data.Context, keep func(string) bool) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if keep(e.KeyStr) {
			zv := data.NewZVal(e.Value)
			zv.Name = e.KeyStr
			out.List = append(out.List, zv)
		}
	}
	return newBag(ctx, out), nil
}

func keySet(v data.Value) map[string]bool {
	m := map[string]bool{}
	for _, e := range kit.Entries(v) {
		if _, isInt := e.Key.(*data.IntValue); isInt || e.KeyStr == "" {
			m[e.Value.AsString()] = true
		} else {
			m[e.KeyStr] = true
		}
	}
	if len(m) == 0 && v != nil {
		m[v.AsString()] = true
	}
	return m
}

func bagOnly(ctx data.Context) (data.GetValue, data.Control) {
	keys := keySet(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return keys[k] })
}
func bagExcept(ctx data.Context) (data.GetValue, data.Control) {
	keys := keySet(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return !keys[k] })
}

func extractPropNames(props data.Value) map[string]bool {
	m := map[string]bool{}
	for _, e := range kit.Entries(props) {
		name := e.KeyStr
		if _, isInt := e.Key.(*data.IntValue); isInt || name == "" {
			name = e.Value.AsString()
		}
		m[name] = true
		// kebab
		kebab := strings.ReplaceAll(name, "_", "-")
		var b strings.Builder
		for i, r := range kebab {
			if r >= 'A' && r <= 'Z' {
				if i > 0 {
					b.WriteByte('-')
				}
				b.WriteRune(r + 32)
			} else {
				b.WriteRune(r)
			}
		}
		m[b.String()] = true
	}
	return m
}

func bagExceptProps(ctx data.Context) (data.GetValue, data.Control) {
	props := extractPropNames(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return !props[k] })
}
func bagOnlyProps(ctx data.Context) (data.GetValue, data.Control) {
	props := extractPropNames(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return props[k] })
}

func bagWhereStartsWith(ctx data.Context) (data.GetValue, data.Control) {
	prefix := kit.Arg(ctx, 0).AsString()
	return bagFilterKeys(ctx, func(k string) bool { return strings.HasPrefix(k, prefix) })
}
func bagWhereDoesntStartWith(ctx data.Context) (data.GetValue, data.Control) {
	prefix := kit.Arg(ctx, 0).AsString()
	return bagFilterKeys(ctx, func(k string) bool { return !strings.HasPrefix(k, prefix) })
}

func bagToHtml(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	var b strings.Builder
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if e.Value == nil || kit.IsNull(e.Value) {
			continue
		}
		if bv, ok := e.Value.(*data.BoolValue); ok {
			okv, _ := bv.AsBool()
			if !okv {
				continue
			}
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(e.KeyStr)
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(e.KeyStr)
		b.WriteString(`="`)
		b.WriteString(html.EscapeString(e.Value.AsString()))
		b.WriteByte('"')
	}
	return data.NewStringValue(b.String()), nil
}

func bagGetIterator(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	vm := ctx.GetVM()
	cls, ok := vm.GetClass("ArrayIterator")
	if !ok {
		var c2 data.Control
		cls, c2 = vm.GetOrLoadClass("ArrayIterator")
		if c2 != nil {
			return bagAttrs(cv), nil
		}
	}
	it := data.NewClassValue(cls, ctx.CreateBaseContext())
	if ctor := cls.GetConstruct(); ctor != nil {
		nctx := it.CreateContext(ctor.GetVariables())
		data.BindDeclaredArgs(nctx, ctor, []data.Value{bagAttrs(cv)})
		if _, ctl := ctor.Call(nctx); ctl != nil {
			return nil, ctl
		}
	}
	return it, nil
}

func bagOffsetGet(ctx data.Context) (data.GetValue, data.Control) { return bagGet(ctx) }

func bagOffsetSet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	val := kit.Arg(ctx, 1)
	attrs := bagAttrs(cv)
	for i, z := range attrs.List {
		if z != nil && z.Name == key {
			attrs.List[i] = data.NewZVal(val)
			attrs.List[i].Name = key
			return data.NewNullValue(), nil
		}
	}
	zv := data.NewZVal(val)
	zv.Name = key
	attrs.List = append(attrs.List, zv)
	return data.NewNullValue(), nil
}

func bagOffsetUnset(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	attrs := bagAttrs(cv)
	out := make([]*data.ZVal, 0, len(attrs.List))
	for _, z := range attrs.List {
		if z != nil && z.Name == key {
			continue
		}
		out = append(out, z)
	}
	attrs.List = out
	return data.NewNullValue(), nil
}

func bagShouldEscape(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Arg(ctx, 1)
	if cv, ok := kit.Unwrap(v).(*data.ClassValue); ok && cv.Class != nil {
		if cv.Class.GetName() == appendableClassName {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

// AppendableClass 对齐 AppendableAttributeValue。
type AppendableClass struct {
	node.Node
	methods map[string]data.Method
}

func NewAppendableClass() data.ClassStmt {
	c := &AppendableClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = kit.InstanceMethod("__construct", []string{"value"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := kit.Receiver(ctx)
		if cv == nil {
			return data.NewNullValue(), nil
		}
		_ = cv.SetProperty("value", kit.Arg(ctx, 0))
		return cv, nil
	})
	c.methods["__tostring"] = kit.InstanceMethod("__toString", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := kit.Receiver(ctx)
		if cv == nil {
			return data.NewStringValue(""), nil
		}
		v, _ := cv.GetProperty("value")
		if v == nil {
			return data.NewStringValue(""), nil
		}
		return data.NewStringValue(v.AsString()), nil
	})
	return c
}

func (c *AppendableClass) GetName() string    { return appendableClassName }
func (c *AppendableClass) GetExtend() *string { return nil }
func (c *AppendableClass) GetImplements() []string {
	return []string{"Stringable"}
}
func (c *AppendableClass) GetProperty(name string) (data.Property, bool) {
	if name == "value" {
		return node.NewProperty(nil, "value", "public", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *AppendableClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "value", "public", false, data.NewNullValue())}
}
func (c *AppendableClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *AppendableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *AppendableClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *AppendableClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *AppendableClass) GetStaticMethod(string) (data.Method, bool) { return nil, false }
