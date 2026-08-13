package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

const (
	cfiCallbackKey = "__cfi_callback__"
	cfiFlagKey     = "__cfi_flag__"
)

// CallbackFilterIteratorClass ?? PHP �?CallbackFilterIterator
type CallbackFilterIteratorClass struct {
	node.Node
}

func NewCallbackFilterIteratorClass() *CallbackFilterIteratorClass {
	return &CallbackFilterIteratorClass{}
}

func (c *CallbackFilterIteratorClass) GetName() string { return "CallbackFilterIterator" }

func (c *CallbackFilterIteratorClass) GetExtend() *string {
	parent := "FilterIterator"
	return &parent
}

func (c *CallbackFilterIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator"}
}

func (c *CallbackFilterIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *CallbackFilterIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *CallbackFilterIteratorClass) GetConstruct() data.Method {
	return &CallbackFilterIteratorConstructMethod{}
}

func (c *CallbackFilterIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(filterValidKey, data.NewBoolValue(false))
	cv.SetProperty(filterCurValKey, data.NewNullValue())
	cv.SetProperty(filterCurKeyKey, data.NewNullValue())
	cv.SetProperty(cfiCallbackKey, data.NewNullValue())
	cv.SetProperty(cfiFlagKey, data.NewIntValue(0))
	return cv, nil
}

func (c *CallbackFilterIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &CallbackFilterIteratorConstructMethod{}, true
	case "accept":
		return &CallbackFilterIteratorAcceptMethod{}, true
	}
	return nil, false
}

func (c *CallbackFilterIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&CallbackFilterIteratorConstructMethod{},
		&CallbackFilterIteratorAcceptMethod{},
	}
}

func splInvokeCallable(ctx data.Context, callback data.GetValue, args []data.Value) (data.GetValue, data.Control) {
	fn := core.NewCallUserFuncFunction()
	callCtx := ctx.CreateContext(fn.GetVariables())
	callCtx.SetIndexZVal(0, data.NewZVal(splAsValue(callback)))
	for i, arg := range args {
		callCtx.SetIndexZVal(i+1, data.NewZVal(arg))
	}
	return fn.Call(callCtx)
}

func splValueToBool(v data.GetValue) bool {
	if bv, ok := v.(*data.BoolValue); ok {
		return bv.Value
	}
	if ab, ok := v.(data.AsBool); ok {
		b, _ := ab.AsBool()
		return b
	}
	return v != nil && v != data.NewNullValue()
}

func cfiCopyState(from, to *data.ClassValue) {
	if cb, _ := from.ObjectValue.GetProperty(cfiCallbackKey); cb != nil {
		to.ObjectValue.SetProperty(cfiCallbackKey, cb)
	}
	if flag, _ := from.ObjectValue.GetProperty(cfiFlagKey); flag != nil {
		to.ObjectValue.SetProperty(cfiFlagKey, flag)
	}
}

type CallbackFilterIteratorConstructMethod struct{}

func (m *CallbackFilterIteratorConstructMethod) GetName() string { return "__construct" }
func (m *CallbackFilterIteratorConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *CallbackFilterIteratorConstructMethod) GetIsStatic() bool         { return false }
func (m *CallbackFilterIteratorConstructMethod) GetReturnType() data.Types { return nil }
func (m *CallbackFilterIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
		node.NewParameter(nil, "callback", 1, nil, nil),
		node.NewParameter(nil, "mode", 2, data.NewIntValue(0), data.Int{}),
	}
}
func (m *CallbackFilterIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
		node.NewVariable(nil, "callback", 1, data.Mixed{}),
		node.NewVariable(nil, "mode", 2, data.Int{}),
	}
}
func (m *CallbackFilterIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	it, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	mode, _ := ctx.GetIndexValue(2)
	filterSetInner(cv, it)
	if cb != nil {
		cv.ObjectValue.SetProperty(cfiCallbackKey, cb)
	}
	if mode != nil {
		cv.ObjectValue.SetProperty(cfiFlagKey, mode)
	} else {
		cv.ObjectValue.SetProperty(cfiFlagKey, data.NewIntValue(0))
	}
	filterSetValid(cv, false)
	filterSetCurVal(cv, data.NewNullValue())
	filterSetCurKey(cv, data.NewNullValue())
	return nil, nil
}

type CallbackFilterIteratorAcceptMethod struct{}

func (m *CallbackFilterIteratorAcceptMethod) GetName() string            { return "accept" }
func (m *CallbackFilterIteratorAcceptMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CallbackFilterIteratorAcceptMethod) GetIsStatic() bool          { return false }
func (m *CallbackFilterIteratorAcceptMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *CallbackFilterIteratorAcceptMethod) GetParams() []data.GetValue { return nil }
func (m *CallbackFilterIteratorAcceptMethod) GetVariables() []data.Variable {
	return nil
}
func (m *CallbackFilterIteratorAcceptMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(true), nil
	}
	cb, _ := cv.ObjectValue.GetProperty(cfiCallbackKey)
	if cb == nil || cb == data.NewNullValue() {
		return data.NewBoolValue(true), nil
	}
	flagVal, _ := cv.ObjectValue.GetProperty(cfiFlagKey)
	flag := 0
	if iv, ok := flagVal.(*data.IntValue); ok {
		flag = iv.Value
	}
	cur := filterGetCurVal(cv)
	key := filterGetCurKey(cv)
	var args []data.Value
	switch flag {
	case 1: // ARRAY_FILTER_USE_KEY
		args = []data.Value{key}
	case 2: // ARRAY_FILTER_USE_BOTH
		args = []data.Value{cur, key}
	default:
		args = []data.Value{cur}
	}
	result, ctl := splInvokeCallable(ctx, cb, args)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(splValueToBool(result)), nil
}

// RecursiveCallbackFilterIteratorClass ?? PHP �?RecursiveCallbackFilterIterator
type RecursiveCallbackFilterIteratorClass struct {
	node.Node
}

func NewRecursiveCallbackFilterIteratorClass() *RecursiveCallbackFilterIteratorClass {
	return &RecursiveCallbackFilterIteratorClass{}
}

func (c *RecursiveCallbackFilterIteratorClass) GetName() string {
	return "RecursiveCallbackFilterIterator"
}

func (c *RecursiveCallbackFilterIteratorClass) GetExtend() *string {
	parent := "RecursiveFilterIterator"
	return &parent
}

func (c *RecursiveCallbackFilterIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator", "RecursiveIterator"}
}

func (c *RecursiveCallbackFilterIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *RecursiveCallbackFilterIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *RecursiveCallbackFilterIteratorClass) GetConstruct() data.Method {
	return &CallbackFilterIteratorConstructMethod{}
}

func (c *RecursiveCallbackFilterIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(filterValidKey, data.NewBoolValue(false))
	cv.SetProperty(filterCurValKey, data.NewNullValue())
	cv.SetProperty(filterCurKeyKey, data.NewNullValue())
	cv.SetProperty(cfiCallbackKey, data.NewNullValue())
	cv.SetProperty(cfiFlagKey, data.NewIntValue(0))
	return cv, nil
}

func (c *RecursiveCallbackFilterIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &CallbackFilterIteratorConstructMethod{}, true
	case "accept":
		return &CallbackFilterIteratorAcceptMethod{}, true
	case "getChildren":
		return &RecursiveCallbackFilterIteratorGetChildrenMethod{}, true
	}
	return nil, false
}

func (c *RecursiveCallbackFilterIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&CallbackFilterIteratorConstructMethod{},
		&CallbackFilterIteratorAcceptMethod{},
		&RecursiveCallbackFilterIteratorGetChildrenMethod{},
	}
}

type RecursiveCallbackFilterIteratorGetChildrenMethod struct{}

func (m *RecursiveCallbackFilterIteratorGetChildrenMethod) GetName() string { return "getChildren" }
func (m *RecursiveCallbackFilterIteratorGetChildrenMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *RecursiveCallbackFilterIteratorGetChildrenMethod) GetIsStatic() bool          { return false }
func (m *RecursiveCallbackFilterIteratorGetChildrenMethod) GetReturnType() data.Types  { return nil }
func (m *RecursiveCallbackFilterIteratorGetChildrenMethod) GetParams() []data.GetValue { return nil }
func (m *RecursiveCallbackFilterIteratorGetChildrenMethod) GetVariables() []data.Variable {
	return nil
}
func (m *RecursiveCallbackFilterIteratorGetChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	inner := filterGetInner(cv)
	childInner, ctl := filterCallInnerMethod(inner, "getChildren")
	if ctl != nil {
		return nil, ctl
	}
	if childInner == nil {
		return data.NewNullValue(), nil
	}
	childCV, ctl := splInstantiateWithArgs(ctx, cv.Class, []data.Value{splAsValue(childInner)})
	if ctl != nil || childCV == nil {
		return childCV, ctl
	}
	cfiCopyState(cv, childCV)
	return childCV, nil
}
