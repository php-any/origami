package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type containerMethod struct{}

func (containerMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (containerMethod) GetIsStatic() bool          { return false }

type containerStaticMethod struct{}

func (containerStaticMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (containerStaticMethod) GetIsStatic() bool          { return true }

// --- bind ---

type ContainerBindMethod struct{ containerMethod }

func (m *ContainerBindMethod) GetName() string { return "bind" }
func (m *ContainerBindMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("abstract", 0),
		data.NewParameterDefault("concrete", 1, data.NewNullValue(), nil),
	}
}
func (m *ContainerBindMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("abstract", 0, data.NewBaseType("string")),
		data.NewVariable("concrete", 1, nil),
	}
}
func (m *ContainerBindMethod) GetReturnType() data.Types { return nil }
func (m *ContainerBindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	abstract, concrete, factory, acl := resolveConcreteArg(ctx)
	if acl != nil {
		return nil, acl
	}
	if factory != nil {
		engine.BindFactory(abstract, factory)
	} else {
		engine.Bind(abstract, concrete)
	}
	return data.NewNullValue(), nil
}

// --- singleton ---

type ContainerSingletonMethod struct{ containerMethod }

func (m *ContainerSingletonMethod) GetName() string { return "singleton" }
func (m *ContainerSingletonMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("abstract", 0),
		data.NewParameterDefault("concrete", 1, data.NewNullValue(), nil),
	}
}
func (m *ContainerSingletonMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("abstract", 0, data.NewBaseType("string")),
		data.NewVariable("concrete", 1, nil),
	}
}
func (m *ContainerSingletonMethod) GetReturnType() data.Types { return nil }
func (m *ContainerSingletonMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	abstract, concrete, factory, acl := resolveConcreteArg(ctx)
	if acl != nil {
		return nil, acl
	}
	if factory != nil {
		engine.SingletonFactory(abstract, factory)
	} else {
		engine.Singleton(abstract, concrete)
	}
	return data.NewNullValue(), nil
}

// --- scoped ---

type ContainerScopedMethod struct{ containerMethod }

func (m *ContainerScopedMethod) GetName() string { return "scoped" }
func (m *ContainerScopedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("abstract", 0),
		data.NewParameterDefault("concrete", 1, data.NewNullValue(), nil),
	}
}
func (m *ContainerScopedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("abstract", 0, data.NewBaseType("string")),
		data.NewVariable("concrete", 1, nil),
	}
}
func (m *ContainerScopedMethod) GetReturnType() data.Types { return nil }
func (m *ContainerScopedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	abstract, concrete, factory, acl := resolveConcreteArg(ctx)
	if acl != nil {
		return nil, acl
	}
	if factory != nil {
		engine.ScopedFactory(abstract, factory)
	} else {
		engine.Scoped(abstract, concrete)
	}
	return data.NewNullValue(), nil
}

// --- instance ---

type ContainerInstanceMethod struct{ containerMethod }

func (m *ContainerInstanceMethod) GetName() string { return "instance" }
func (m *ContainerInstanceMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("abstract", 0),
		data.NewParameter("instance", 1),
	}
}
func (m *ContainerInstanceMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("abstract", 0, data.NewBaseType("string")),
		data.NewVariable("instance", 1, data.NewBaseType("object")),
	}
}
func (m *ContainerInstanceMethod) GetReturnType() data.Types { return nil }
func (m *ContainerInstanceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	abstract, acl := stringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	inst, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 instance 参数"))
	}
	engine.Instance(abstract, inst)
	return data.NewNullValue(), nil
}

// --- make ---

type ContainerMakeMethod struct{ containerMethod }

func (m *ContainerMakeMethod) GetName() string { return "make" }
func (m *ContainerMakeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("abstract", 0),
		data.NewParameterDefault("parameters", 1, data.NewNullValue(), nil),
	}
}
func (m *ContainerMakeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("abstract", 0, data.NewBaseType("string")),
		data.NewVariable("parameters", 1, nil),
	}
}
func (m *ContainerMakeMethod) GetReturnType() data.Types { return data.NewBaseType("object") }
func (m *ContainerMakeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	abstract, acl := stringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	var params []data.GetValue
	if v, ok := ctx.GetIndexValue(1); ok {
		if arr, isArr := v.(*data.ArrayValue); isArr {
			for _, item := range arr.List {
				params = append(params, item.Value)
			}
		}
	}
	return engine.Make(ctx, abstract, params)
}

// --- has ---

type ContainerHasMethod struct{ containerMethod }

func (m *ContainerHasMethod) GetName() string { return "has" }
func (m *ContainerHasMethod) GetParams() []data.GetValue {
	return []data.GetValue{data.NewParameter("abstract", 0)}
}
func (m *ContainerHasMethod) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("abstract", 0, data.NewBaseType("string"))}
}
func (m *ContainerHasMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
func (m *ContainerHasMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	abstract, acl := stringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	return data.NewBoolValue(engine.Has(abstract)), nil
}

// --- isShared ---

type ContainerIsSharedMethod struct{ containerMethod }

func (m *ContainerIsSharedMethod) GetName() string { return "isShared" }
func (m *ContainerIsSharedMethod) GetParams() []data.GetValue {
	return []data.GetValue{data.NewParameter("abstract", 0)}
}
func (m *ContainerIsSharedMethod) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("abstract", 0, data.NewBaseType("string"))}
}
func (m *ContainerIsSharedMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
func (m *ContainerIsSharedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	abstract, acl := stringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	return data.NewBoolValue(engine.IsShared(abstract)), nil
}

// --- alias ---

type ContainerAliasMethod struct{ containerMethod }

func (m *ContainerAliasMethod) GetName() string { return "alias" }
func (m *ContainerAliasMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("alias", 0),
		data.NewParameter("abstract", 1),
	}
}
func (m *ContainerAliasMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("alias", 0, data.NewBaseType("string")),
		data.NewVariable("abstract", 1, data.NewBaseType("string")),
	}
}
func (m *ContainerAliasMethod) GetReturnType() data.Types { return nil }
func (m *ContainerAliasMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	alias, acl := stringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	abstract, acl := stringArg(ctx, 1)
	if acl != nil {
		return nil, acl
	}
	engine.Alias(alias, abstract)
	return data.NewNullValue(), nil
}

// --- getInstance ---

type ContainerGetInstanceMethod struct{ containerStaticMethod }

func (m *ContainerGetInstanceMethod) GetName() string               { return "getInstance" }
func (m *ContainerGetInstanceMethod) GetParams() []data.GetValue    { return nil }
func (m *ContainerGetInstanceMethod) GetVariables() []data.Variable { return nil }
func (m *ContainerGetInstanceMethod) GetReturnType() data.Types {
	return data.NewBaseType("Container\\Container")
}
func (m *ContainerGetInstanceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return ensureDefaultInstance(ctx)
}

// --- registerProviders ---

type ContainerRegisterProvidersMethod struct{ containerMethod }

func (m *ContainerRegisterProvidersMethod) GetName() string { return "registerProviders" }
func (m *ContainerRegisterProvidersMethod) GetParams() []data.GetValue {
	return []data.GetValue{data.NewParameter("providers", 0)}
}
func (m *ContainerRegisterProvidersMethod) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("providers", 0, nil)}
}
func (m *ContainerRegisterProvidersMethod) GetReturnType() data.Types { return nil }
func (m *ContainerRegisterProvidersMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	containerVal, ok := classValueFromCtx(ctx)
	if !ok {
		return nil, utils.NewThrow(errors.New("registerProviders 必须在 Container 实例上调用"))
	}
	arrVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 providers 参数"))
	}
	arr, ok := arrVal.(*data.ArrayValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("providers 必须是数组"))
	}

	var providers []data.GetValue
	for _, item := range arr.List {
		providers = append(providers, item.Value)
	}

	registered := make([]*data.ClassValue, 0, len(providers))
	for _, p := range providers {
		className, acl := providerClassName(p)
		if acl != nil {
			return nil, acl
		}
		inst, acl := instantiateProvider(ctx, className)
		if acl != nil {
			return nil, acl
		}
		cv, ok := inst.(*data.ClassValue)
		if !ok {
			return nil, utils.NewThrow(errors.New("ServiceProvider 实例化失败"))
		}
		if acl := cv.SetProperty("container", containerVal); acl != nil {
			return nil, acl
		}
		registered = append(registered, cv)
	}

	for _, cv := range registered {
		if method, ok := cv.GetMethod("register"); ok {
			fnCtx := cv.CreateContext(method.GetVariables())
			if _, acl := method.Call(fnCtx); acl != nil {
				return nil, acl
			}
		}
	}
	for _, cv := range registered {
		if method, ok := cv.GetMethod("boot"); ok {
			fnCtx := cv.CreateContext(method.GetVariables())
			if _, acl := method.Call(fnCtx); acl != nil {
				return nil, acl
			}
		}
	}
	return data.NewNullValue(), nil
}

func providerClassName(v data.GetValue) (string, data.Control) {
	if s, ok := v.(data.AsString); ok {
		return s.AsString(), nil
	}
	return "", utils.NewThrow(errors.New("provider 必须是类名字符串"))
}

func instantiateProvider(ctx data.Context, className string) (data.GetValue, data.Control) {
	stmt, acl := ctx.GetVM().GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}
	return instantiateClass(stmt, nil, ctx)
}

// --- createScope ---

type ContainerCreateScopeMethod struct{ containerMethod }

func (m *ContainerCreateScopeMethod) GetName() string               { return "createScope" }
func (m *ContainerCreateScopeMethod) GetParams() []data.GetValue    { return nil }
func (m *ContainerCreateScopeMethod) GetVariables() []data.Variable { return nil }
func (m *ContainerCreateScopeMethod) GetReturnType() data.Types {
	return data.NewBaseType("Container\\Scope")
}
func (m *ContainerCreateScopeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	return newScopeValue(engine, ctx)
}

// --- scan ---

type ContainerScanMethod struct{ containerMethod }

func (m *ContainerScanMethod) GetName() string { return "scan" }
func (m *ContainerScanMethod) GetParams() []data.GetValue {
	return []data.GetValue{data.NewParameter("directory", 0)}
}
func (m *ContainerScanMethod) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("directory", 0, data.NewBaseType("string"))}
}
func (m *ContainerScanMethod) GetReturnType() data.Types { return nil }
func (m *ContainerScanMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	engine, acl := engineFromCtx(ctx)
	if acl != nil {
		return nil, acl
	}
	dir, acl := stringArg(ctx, 0)
	if acl != nil {
		return nil, acl
	}
	restore := setRegisteringEngine(engine)
	defer restore()
	if acl := scanDirectory(ctx.GetVM(), dir); acl != nil {
		return nil, acl
	}
	return data.NewNullValue(), nil
}
