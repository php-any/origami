package reflection

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Origami 没有 PHP 8.4 ghost/proxy 懒对象。这些方法按「立即实体化」实现，
// 让 Livewire IsLazy / Filament 编辑页能过 isUninitializedLazyObject / newLazyProxy。

type ReflectionClassIsUninitializedLazyObjectMethod struct{}

func (m *ReflectionClassIsUninitializedLazyObjectMethod) GetName() string {
	return "isUninitializedLazyObject"
}
func (m *ReflectionClassIsUninitializedLazyObjectMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionClassIsUninitializedLazyObjectMethod) GetIsStatic() bool { return false }
func (m *ReflectionClassIsUninitializedLazyObjectMethod) GetReturnType() data.Types {
	return data.Bool{}
}
var reflectionClassIsUninitializedLazyObjectMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
}

func (m *ReflectionClassIsUninitializedLazyObjectMethod) GetParams() []data.GetValue {
	return reflectionClassIsUninitializedLazyObjectMethodGetParams
}
var reflectionClassIsUninitializedLazyObjectMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
}

func (m *ReflectionClassIsUninitializedLazyObjectMethod) GetVariables() []data.Variable {
	return reflectionClassIsUninitializedLazyObjectMethodGetVariables
}
func (m *ReflectionClassIsUninitializedLazyObjectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

type ReflectionClassNewLazyProxyMethod struct{}

func (m *ReflectionClassNewLazyProxyMethod) GetName() string { return "newLazyProxy" }
func (m *ReflectionClassNewLazyProxyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionClassNewLazyProxyMethod) GetIsStatic() bool { return false }
func (m *ReflectionClassNewLazyProxyMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
var reflectionClassNewLazyProxyMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "factory", 0, nil, nil),
	node.NewParameter(nil, "options", 1, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
}

func (m *ReflectionClassNewLazyProxyMethod) GetParams() []data.GetValue {
	return reflectionClassNewLazyProxyMethodGetParams
}
var reflectionClassNewLazyProxyMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "factory", 0, nil),
	node.NewVariable(nil, "options", 1, data.NewBaseType("int")),
}

func (m *ReflectionClassNewLazyProxyMethod) GetVariables() []data.Variable {
	return reflectionClassNewLazyProxyMethodGetVariables
}
func (m *ReflectionClassNewLazyProxyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	factory, _ := ctx.GetIndexValue(0)
	return callPhpCallable(ctx, factory)
}

type ReflectionClassNewLazyGhostMethod struct{}

func (m *ReflectionClassNewLazyGhostMethod) GetName() string { return "newLazyGhost" }
func (m *ReflectionClassNewLazyGhostMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionClassNewLazyGhostMethod) GetIsStatic() bool { return false }
func (m *ReflectionClassNewLazyGhostMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
var reflectionClassNewLazyGhostMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "initializer", 0, nil, nil),
	node.NewParameter(nil, "options", 1, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
}

func (m *ReflectionClassNewLazyGhostMethod) GetParams() []data.GetValue {
	return reflectionClassNewLazyGhostMethodGetParams
}
var reflectionClassNewLazyGhostMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "initializer", 0, nil),
	node.NewVariable(nil, "options", 1, data.NewBaseType("int")),
}

func (m *ReflectionClassNewLazyGhostMethod) GetVariables() []data.Variable {
	return reflectionClassNewLazyGhostMethodGetVariables
}
func (m *ReflectionClassNewLazyGhostMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Class does not exist"))
	}
	object := data.NewClassValue(classStmt, ctx.CreateBaseContext())
	init, _ := ctx.GetIndexValue(0)
	if init == nil {
		return object, nil
	}
	if _, ctl := callPhpCallable(ctx, init, object); ctl != nil {
		return nil, ctl
	}
	return object, nil
}

type ReflectionClassInitializeLazyObjectMethod struct{}

func (m *ReflectionClassInitializeLazyObjectMethod) GetName() string {
	return "initializeLazyObject"
}
func (m *ReflectionClassInitializeLazyObjectMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionClassInitializeLazyObjectMethod) GetIsStatic() bool { return false }
func (m *ReflectionClassInitializeLazyObjectMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
var reflectionClassInitializeLazyObjectMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
}

func (m *ReflectionClassInitializeLazyObjectMethod) GetParams() []data.GetValue {
	return reflectionClassInitializeLazyObjectMethodGetParams
}
var reflectionClassInitializeLazyObjectMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
}

func (m *ReflectionClassInitializeLazyObjectMethod) GetVariables() []data.Variable {
	return reflectionClassInitializeLazyObjectMethodGetVariables
}
func (m *ReflectionClassInitializeLazyObjectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj, _ := ctx.GetIndexValue(0)
	if obj == nil {
		return data.NewNullValue(), nil
	}
	return obj, nil
}

type ReflectionClassMarkLazyObjectAsInitializedMethod struct{}

func (m *ReflectionClassMarkLazyObjectAsInitializedMethod) GetName() string {
	return "markLazyObjectAsInitialized"
}
func (m *ReflectionClassMarkLazyObjectAsInitializedMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionClassMarkLazyObjectAsInitializedMethod) GetIsStatic() bool { return false }
func (m *ReflectionClassMarkLazyObjectAsInitializedMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
var reflectionClassMarkLazyObjectAsInitializedMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
}

func (m *ReflectionClassMarkLazyObjectAsInitializedMethod) GetParams() []data.GetValue {
	return reflectionClassMarkLazyObjectAsInitializedMethodGetParams
}
var reflectionClassMarkLazyObjectAsInitializedMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
}

func (m *ReflectionClassMarkLazyObjectAsInitializedMethod) GetVariables() []data.Variable {
	return reflectionClassMarkLazyObjectAsInitializedMethodGetVariables
}
func (m *ReflectionClassMarkLazyObjectAsInitializedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj, _ := ctx.GetIndexValue(0)
	if obj == nil {
		return data.NewNullValue(), nil
	}
	return obj, nil
}

func callPhpCallable(ctx data.Context, cb data.Value, args ...data.Value) (data.GetValue, data.Control) {
	if cb == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("callable expected"))
	}
	switch f := cb.(type) {
	case *data.FuncValue:
		if f.Value == nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("callable expected"))
		}
		callCtx := ctx.CreateContext(f.Value.GetVariables())
		data.BindDeclaredArgs(callCtx, f.Value, args)
		return f.Call(callCtx)
	case *data.BoundFuncValue:
		callCtx := ctx.CreateContext(f.Value.GetVariables())
		data.BindDeclaredArgs(callCtx, f.Value, args)
		return f.Call(callCtx)
	default:
		return nil, data.NewErrorThrow(nil, fmt.Errorf("callable expected"))
	}
}
