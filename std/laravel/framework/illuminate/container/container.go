package container

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const containerClassName = "Illuminate\\Container\\Container"

type ctnState struct {
	mu        sync.Mutex
	bindings  map[string]binding
	instances map[string]data.Value
	aliases   map[string]string
	resolved  map[string]bool
	buildStack []string
}

type binding struct {
	concrete data.Value
	shared   bool
}

// ctnStateProp：Go 侧状态挂在对象隐蔽属性上，由 GC 随对象回收（见 kit.CachedState）。
// 早前用「全局表 + __origami_ctn_id」的写法：每 new 一个 Container 就多一条永不释放的
// 记录，state 里还攥着 requests/instances 整棵对象图。
const ctnStateProp = "__origami_ctn_state"

var globalCtn atomic.Value // *data.ClassValue

func newCtnState() *ctnState {
	return &ctnState{bindings: map[string]binding{}, instances: map[string]data.Value{}, aliases: map[string]string{}, resolved: map[string]bool{}}
}

func ctnOf(cv *data.ClassValue) *ctnState {
	if cv == nil {
		return newCtnState()
	}
	return kit.CachedState(cv, ctnStateProp, newCtnState)
}

// ContainerClass 核心 IoC：bind / singleton / instance / make / bound（默认不开，见 Load）。
type ContainerClass struct {
	node.Node
	methods map[string]data.Method
}

func NewContainerClass() data.ClassStmt {
	c := &ContainerClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *ContainerClass) GetName() string    { return containerClassName }
func (c *ContainerClass) GetExtend() *string { return nil }
func (c *ContainerClass) GetImplements() []string {
	return []string{"ArrayAccess", "Illuminate\\Contracts\\Container\\Container"}
}
func (c *ContainerClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "bindings", "instances", "aliases", "resolved", "scopedInstances", "methodBindings", "extenders", "tags", "buildStack", "with", "contextual":
		return node.NewProperty(nil, name, "protected", false, data.NewArrayValue(nil)), true
	case "__origami_ctn_id":
		return node.NewProperty(nil, name, "protected", false, data.NewIntValue(0)), true
	}
	return nil, false
}
func (c *ContainerClass) GetPropertyList() []data.Property {
	out := []data.Property{}
	for _, n := range []string{"bindings", "instances", "aliases", "resolved", "__origami_ctn_id"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *ContainerClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *ContainerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ContainerClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *ContainerClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *ContainerClass) GetStaticMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "getinstance", "setinstance":
		return c.GetMethod(name)
	}
	return c.GetMethod(name)
}

func (c *ContainerClass) register() {
	c.methods["__construct"] = kit.InstanceMethod("__construct", nil, ctnConstruct)
	c.methods["bound"] = kit.InstanceMethod("bound", []string{"abstract"}, ctnBound)
	c.methods["has"] = kit.InstanceMethod("has", []string{"id"}, ctnBound)
	c.methods["bind"] = kit.InstanceMethodOpt("bind", []string{"abstract", "concrete", "shared"}, 1, ctnBind)
	c.methods["bindif"] = kit.InstanceMethodOpt("bindIf", []string{"abstract", "concrete", "shared"}, 1, ctnBindIf)
	c.methods["singleton"] = kit.InstanceMethodOpt("singleton", []string{"abstract", "concrete"}, 1, ctnSingleton)
	c.methods["singletonif"] = kit.InstanceMethodOpt("singletonIf", []string{"abstract", "concrete"}, 1, ctnSingletonIf)
	c.methods["instance"] = kit.InstanceMethod("instance", []string{"abstract", "instance"}, ctnInstance)
	c.methods["alias"] = kit.InstanceMethod("alias", []string{"abstract", "alias"}, ctnAlias)
	c.methods["make"] = kit.InstanceMethodOpt("make", []string{"abstract", "parameters"}, 1, ctnMake)
	c.methods["makewith"] = kit.InstanceMethod("makeWith", []string{"abstract", "parameters"}, ctnMake)
	c.methods["get"] = kit.InstanceMethod("get", []string{"id"}, ctnGet)
	c.methods["build"] = kit.InstanceMethod("build", []string{"concrete"}, ctnBuild)
	c.methods["getalias"] = kit.InstanceMethod("getAlias", []string{"abstract"}, ctnGetAlias)
	c.methods["flush"] = kit.InstanceMethod("flush", nil, ctnFlush)
	c.methods["forgetinstance"] = kit.InstanceMethod("forgetInstance", []string{"abstract"}, ctnForgetInstance)
	c.methods["forgetinstances"] = kit.InstanceMethod("forgetInstances", nil, ctnForgetInstances)
	c.methods["getbindings"] = kit.InstanceMethod("getBindings", nil, ctnGetBindings)
	c.methods["offsetexists"] = kit.InstanceMethod("offsetExists", []string{"offset"}, ctnBound)
	c.methods["offsetget"] = kit.InstanceMethod("offsetGet", []string{"offset"}, ctnGet)
	c.methods["offsetset"] = kit.InstanceMethod("offsetSet", []string{"offset", "value"}, ctnOffsetSet)
	c.methods["offsetunset"] = kit.InstanceMethod("offsetUnset", []string{"offset"}, ctnForgetInstance)
	c.methods["getinstance"] = kit.StaticMethod("getInstance", nil, -1, ctnGetInstance)
	c.methods["setinstance"] = kit.StaticMethod("setInstance", []string{"container"}, -1, ctnSetInstance)
	c.methods["resolved"] = kit.InstanceMethod("resolved", []string{"abstract"}, ctnResolved)
	c.methods["isshared"] = kit.InstanceMethod("isShared", []string{"abstract"}, ctnIsShared)
}

func ctnRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Container missing $this"))
}

func ctnConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = ctnOf(cv)
	return cv, nil
}

func absName(v data.Value) string {
	v = kit.Unwrap(v)
	if v == nil {
		return ""
	}
	return v.AsString()
}

func (s *ctnState) getAlias(name string) string {
	for {
		if a, ok := s.aliases[name]; ok && a != name {
			name = a
			continue
		}
		return name
	}
}

func ctnBound(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name := absName(kit.Arg(ctx, 0))
	st := ctnOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	name = st.getAlias(name)
	_, okB := st.bindings[name]
	_, okI := st.instances[name]
	return data.NewBoolValue(okB || okI), nil
}

func ctnBind(ctx data.Context) (data.GetValue, data.Control) {
	return ctnBindShared(ctx, false)
}

func ctnBindIf(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	name := absName(kit.Arg(ctx, 0))
	st.mu.Lock()
	name = st.getAlias(name)
	_, ok := st.bindings[name]
	_, okI := st.instances[name]
	st.mu.Unlock()
	if ok || okI {
		return data.NewNullValue(), nil
	}
	return ctnBindShared(ctx, false)
}

func ctnSingleton(ctx data.Context) (data.GetValue, data.Control) {
	return ctnBindShared(ctx, true)
}

func ctnSingletonIf(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	name := absName(kit.Arg(ctx, 0))
	st.mu.Lock()
	name = st.getAlias(name)
	_, ok := st.bindings[name]
	_, okI := st.instances[name]
	st.mu.Unlock()
	if ok || okI {
		return data.NewNullValue(), nil
	}
	return ctnBindShared(ctx, true)
}

func ctnBindShared(ctx data.Context, forceShared bool) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	abstract := absName(kit.Arg(ctx, 0))
	concrete := kit.Arg(ctx, 1)
	shared := forceShared
	if !forceShared {
		if s := kit.Arg(ctx, 2); s != nil {
			if b, ok := s.(data.AsBool); ok {
				shared, _ = b.AsBool()
			}
		}
	}
	if concrete == nil || kit.IsNull(concrete) {
		concrete = data.NewStringValue(abstract)
	}
	st := ctnOf(cv)
	st.mu.Lock()
	st.bindings[abstract] = binding{concrete: concrete, shared: shared}
	delete(st.instances, abstract)
	st.mu.Unlock()
	return data.NewNullValue(), nil
}

func ctnInstance(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	abstract := absName(kit.Arg(ctx, 0))
	inst := kit.Arg(ctx, 1)
	st := ctnOf(cv)
	st.mu.Lock()
	st.instances[abstract] = inst
	delete(st.bindings, abstract)
	st.mu.Unlock()
	return inst, nil
}

func ctnAlias(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	abstract := absName(kit.Arg(ctx, 0))
	alias := absName(kit.Arg(ctx, 1))
	st := ctnOf(cv)
	st.mu.Lock()
	st.aliases[alias] = abstract
	st.mu.Unlock()
	return data.NewNullValue(), nil
}

func ctnGetAlias(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	st.mu.Lock()
	name := st.getAlias(absName(kit.Arg(ctx, 0)))
	st.mu.Unlock()
	return data.NewStringValue(name), nil
}

func ctnMake(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return ctnResolve(ctx, cv, absName(kit.Arg(ctx, 0)))
}

func ctnGet(ctx data.Context) (data.GetValue, data.Control) {
	return ctnMake(ctx)
}

func ctnResolve(ctx data.Context, cv *data.ClassValue, abstract string) (data.GetValue, data.Control) {
	st := ctnOf(cv)
	st.mu.Lock()
	abstract = st.getAlias(abstract)
	if inst, ok := st.instances[abstract]; ok {
		st.mu.Unlock()
		return inst, nil
	}
	b, hasB := st.bindings[abstract]
	st.mu.Unlock()

	var concrete data.Value = data.NewStringValue(abstract)
	shared := false
	if hasB {
		concrete = b.concrete
		shared = b.shared
	}

	var result data.Value
	var ctl data.Control
	switch t := kit.Unwrap(concrete).(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		ret, c := kit.Call(ctx, t, cv)
		ctl = c
		if ctl == nil {
			result = asVal(ret)
		}
	case *data.StringValue:
		ret, c := ctnBuildClass(ctx, cv, t.AsString())
		ctl = c
		if ctl == nil {
			result = asVal(ret)
		}
	default:
		result = concrete
	}
	if ctl != nil {
		return nil, ctl
	}
	if shared {
		st.mu.Lock()
		st.instances[abstract] = result
		st.resolved[abstract] = true
		st.mu.Unlock()
	}
	return result, nil
}

func ctnBuild(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	concrete := kit.Arg(ctx, 0)
	if s, ok := kit.Unwrap(concrete).(*data.StringValue); ok {
		return ctnBuildClass(ctx, cv, s.AsString())
	}
	if f, ok := kit.Unwrap(concrete).(*data.FuncValue); ok {
		return kit.Call(ctx, f, cv)
	}
	return concrete, nil
}

func ctnBuildClass(ctx data.Context, ctn *data.ClassValue, className string) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("no VM"))
	}
	cls, ctl := vm.GetOrLoadClass(className)
	if ctl != nil {
		return nil, ctl
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	ctor := cls.GetConstruct()
	if ctor == nil {
		if m, ok := cv.GetMethod("__construct"); ok {
			ctor = m
		}
	}
	if ctor == nil {
		return cv, nil
	}
	args := make([]data.Value, 0)
	for _, p := range ctor.GetParams() {
		// 尝试从参数类型名 autowire
		typeName := paramTypeName(p)
		if typeName == "" {
			args = append(args, data.NewNullValue())
			continue
		}
		dep, ctl := ctnResolve(ctx, ctn, typeName)
		if ctl != nil {
			// 可选参数则 null
			args = append(args, data.NewNullValue())
			continue
		}
		args = append(args, asVal(dep))
	}
	nctx := cv.CreateContext(ctor.GetVariables())
	data.BindDeclaredArgs(nctx, ctor, args)
	if _, ctl := ctor.Call(nctx); ctl != nil {
		return nil, ctl
	}
	return cv, nil
}

func paramTypeName(p data.GetValue) string {
	type getter interface{ GetType() data.Types }
	if g, ok := p.(getter); ok && g.GetType() != nil {
		return g.GetType().String()
	}
	return ""
}

func asVal(v data.GetValue) data.Value {
	if v == nil {
		return data.NewNullValue()
	}
	if val, ok := v.(data.Value); ok {
		return val
	}
	return data.NewNullValue()
}

func ctnFlush(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	st.mu.Lock()
	st.bindings = map[string]binding{}
	st.instances = map[string]data.Value{}
	st.aliases = map[string]string{}
	st.resolved = map[string]bool{}
	st.mu.Unlock()
	return data.NewNullValue(), nil
}

func ctnForgetInstance(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	st.mu.Lock()
	delete(st.instances, absName(kit.Arg(ctx, 0)))
	st.mu.Unlock()
	return data.NewNullValue(), nil
}

func ctnForgetInstances(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	st.mu.Lock()
	st.instances = map[string]data.Value{}
	st.mu.Unlock()
	return data.NewNullValue(), nil
}

func ctnGetBindings(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for k := range st.bindings {
		zv := data.NewZVal(data.NewBoolValue(true))
		zv.Name = k
		out.List = append(out.List, zv)
	}
	return out, nil
}

func ctnOffsetSet(ctx data.Context) (data.GetValue, data.Control) {
	// $c[$k] = $v  → bind if Closure else instance
	v := kit.Arg(ctx, 1)
	switch kit.Unwrap(v).(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return ctnBindShared(ctx, false)
	default:
		return ctnInstance(ctx)
	}
}

func ctnGetInstance(ctx data.Context) (data.GetValue, data.Control) {
	if v := globalCtn.Load(); v != nil {
		return v.(*data.ClassValue), nil
	}
	// create new
	vm := ctx.GetVM()
	cls, _ := vm.GetClass(containerClassName)
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	_ = ctnOf(cv)
	globalCtn.Store(cv)
	return cv, nil
}

func ctnSetInstance(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Arg(ctx, 0)
	if cv, ok := kit.Unwrap(v).(*data.ClassValue); ok {
		globalCtn.Store(cv)
		return cv, nil
	}
	globalCtn.Store((*data.ClassValue)(nil))
	return data.NewNullValue(), nil
}

func ctnResolved(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	name := st.getAlias(absName(kit.Arg(ctx, 0)))
	_, ok := st.resolved[name]
	_, okI := st.instances[name]
	return data.NewBoolValue(ok || okI), nil
}

func ctnIsShared(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := ctnRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st := ctnOf(cv)
	st.mu.Lock()
	defer st.mu.Unlock()
	name := st.getAlias(absName(kit.Arg(ctx, 0)))
	if _, ok := st.instances[name]; ok {
		return data.NewBoolValue(true), nil
	}
	if b, ok := st.bindings[name]; ok {
		return data.NewBoolValue(b.shared), nil
	}
	return data.NewBoolValue(false), nil
}
