package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FilterIteratorClass 实现 PHP 的 FilterIterator 抽象类
// 子类通过重写 accept() 方法来决定是否接受当前元素
type FilterIteratorClass struct {
	node.Node
	innerIterator *data.ClassValue // 内部迭代器
}

func NewFilterIteratorClass() *FilterIteratorClass {
	return &FilterIteratorClass{}
}

func (f *FilterIteratorClass) GetName() string {
	extend := "FilterIterator"
	_ = extend
	return "FilterIterator"
}

func (f *FilterIteratorClass) GetExtend() *string { return nil }

func (f *FilterIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator"}
}

func (f *FilterIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (f *FilterIteratorClass) GetPropertyList() []data.Property { return nil }

func (f *FilterIteratorClass) GetConstruct() data.Method {
	return &FilterIteratorConstruct{instance: f}
}

func (f *FilterIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	clone := &FilterIteratorClass{
		innerIterator: f.innerIterator,
	}
	return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}

func (f *FilterIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &FilterIteratorConstruct{instance: f}, true
	case "rewind":
		return &FilterIteratorRewind{instance: f}, true
	case "next":
		return &FilterIteratorNext{instance: f}, true
	case "current":
		return &FilterIteratorCurrent{instance: f}, true
	case "key":
		return &FilterIteratorKey{instance: f}, true
	case "valid":
		return &FilterIteratorValid{instance: f}, true
	case "getInnerIterator":
		return &FilterIteratorGetInnerIterator{instance: f}, true
	case "accept":
		// 默认 accept 返回 true（抽象方法，子类应覆盖）
		return &FilterIteratorAccept{}, true
	}
	return nil, false
}

func (f *FilterIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&FilterIteratorConstruct{instance: f},
		&FilterIteratorRewind{instance: f},
		&FilterIteratorNext{instance: f},
		&FilterIteratorCurrent{instance: f},
		&FilterIteratorKey{instance: f},
		&FilterIteratorValid{instance: f},
		&FilterIteratorGetInnerIterator{instance: f},
		&FilterIteratorAccept{},
	}
}

// callInnerMethod 调用内部迭代器的指定方法
func callInnerMethod(inner *data.ClassValue, name string, ctx data.Context) (data.GetValue, data.Control) {
	if inner == nil {
		return nil, nil
	}
	method, ok := inner.GetMethod(name)
	if !ok {
		return nil, nil
	}
	innerCtx := inner.CreateContext(nil)
	return method.Call(innerCtx)
}

// callAccept 调用 ctx 中当前对象的 accept() 方法（支持子类覆盖）
func callAccept(ctx data.Context) bool {
	cmc, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return true
	}
	method, ok := cmc.ClassValue.GetMethod("accept")
	if !ok {
		return true
	}
	acceptCtx := cmc.ClassValue.CreateContext(nil)
	result, ctrl := method.Call(acceptCtx)
	if ctrl != nil {
		return false
	}
	if bv, ok := result.(*data.BoolValue); ok {
		return bv.Value
	}
	if result != nil {
		return true
	}
	return true
}

// advanceToAccepted 推进内部迭代器直到 accept() 为 true 或无效
func advanceToAccepted(instance *FilterIteratorClass, ctx data.Context) {
	for {
		// 检查内部迭代器是否有效
		validResult, _ := callInnerMethod(instance.innerIterator, "valid", ctx)
		bv, ok := validResult.(*data.BoolValue)
		if !ok || !bv.Value {
			return
		}
		// 检查 accept()
		if callAccept(ctx) {
			return
		}
		// 不接受，前进
		callInnerMethod(instance.innerIterator, "next", ctx)
	}
}

// __construct
type FilterIteratorConstruct struct {
	instance *FilterIteratorClass
}

func (m *FilterIteratorConstruct) GetName() string            { return "__construct" }
func (m *FilterIteratorConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilterIteratorConstruct) GetIsStatic() bool          { return false }
func (m *FilterIteratorConstruct) GetReturnType() data.Types  { return nil }
func (m *FilterIteratorConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
	}
}
func (m *FilterIteratorConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
	}
}
func (m *FilterIteratorConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	iterVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}
	if cv, ok := iterVal.(*data.ClassValue); ok {
		m.instance.innerIterator = cv
	}
	return nil, nil
}

// rewind
type FilterIteratorRewind struct {
	instance *FilterIteratorClass
}

func (m *FilterIteratorRewind) GetName() string               { return "rewind" }
func (m *FilterIteratorRewind) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorRewind) GetIsStatic() bool             { return false }
func (m *FilterIteratorRewind) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorRewind) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorRewind) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorRewind) Call(ctx data.Context) (data.GetValue, data.Control) {
	callInnerMethod(m.instance.innerIterator, "rewind", ctx)
	advanceToAccepted(m.instance, ctx)
	return nil, nil
}

// next
type FilterIteratorNext struct {
	instance *FilterIteratorClass
}

func (m *FilterIteratorNext) GetName() string               { return "next" }
func (m *FilterIteratorNext) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorNext) GetIsStatic() bool             { return false }
func (m *FilterIteratorNext) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorNext) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorNext) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorNext) Call(ctx data.Context) (data.GetValue, data.Control) {
	callInnerMethod(m.instance.innerIterator, "next", ctx)
	advanceToAccepted(m.instance, ctx)
	return nil, nil
}

// current
type FilterIteratorCurrent struct {
	instance *FilterIteratorClass
}

func (m *FilterIteratorCurrent) GetName() string               { return "current" }
func (m *FilterIteratorCurrent) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorCurrent) GetIsStatic() bool             { return false }
func (m *FilterIteratorCurrent) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorCurrent) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorCurrent) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorCurrent) Call(ctx data.Context) (data.GetValue, data.Control) {
	return callInnerMethod(m.instance.innerIterator, "current", ctx)
}

// key
type FilterIteratorKey struct {
	instance *FilterIteratorClass
}

func (m *FilterIteratorKey) GetName() string               { return "key" }
func (m *FilterIteratorKey) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorKey) GetIsStatic() bool             { return false }
func (m *FilterIteratorKey) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorKey) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorKey) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorKey) Call(ctx data.Context) (data.GetValue, data.Control) {
	return callInnerMethod(m.instance.innerIterator, "key", ctx)
}

// valid
type FilterIteratorValid struct {
	instance *FilterIteratorClass
}

func (m *FilterIteratorValid) GetName() string               { return "valid" }
func (m *FilterIteratorValid) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorValid) GetIsStatic() bool             { return false }
func (m *FilterIteratorValid) GetReturnType() data.Types     { return data.Bool{} }
func (m *FilterIteratorValid) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorValid) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorValid) Call(ctx data.Context) (data.GetValue, data.Control) {
	return callInnerMethod(m.instance.innerIterator, "valid", ctx)
}

// getInnerIterator
type FilterIteratorGetInnerIterator struct {
	instance *FilterIteratorClass
}

func (m *FilterIteratorGetInnerIterator) GetName() string               { return "getInnerIterator" }
func (m *FilterIteratorGetInnerIterator) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorGetInnerIterator) GetIsStatic() bool             { return false }
func (m *FilterIteratorGetInnerIterator) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorGetInnerIterator) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorGetInnerIterator) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorGetInnerIterator) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.instance.innerIterator, nil
}

// accept（抽象方法默认实现，子类应覆盖）
type FilterIteratorAccept struct{}

func (m *FilterIteratorAccept) GetName() string               { return "accept" }
func (m *FilterIteratorAccept) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorAccept) GetIsStatic() bool             { return false }
func (m *FilterIteratorAccept) GetReturnType() data.Types     { return data.Bool{} }
func (m *FilterIteratorAccept) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorAccept) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorAccept) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}
