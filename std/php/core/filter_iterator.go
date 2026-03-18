package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FilterIterator 内部状态属性名
// 状态存储在 ClassValue.ObjectValue 属性中，而不是 Go struct 字段
// 这样继承链方法分发时可以通过 ctx 访问到当前 PHP 实例的状态
const (
	filterInnerKey  = "__fi_inner__"  // 内部迭代器
	filterValidKey  = "__fi_valid__"  // 是否有效
	filterCurValKey = "__fi_curval__" // 当前元素値缓存
	filterCurKeyKey = "__fi_curkey__" // 当前元素键缓存
)

// FilterIteratorClass 实现 PHP 的 FilterIterator 抽象类
// FilterIterator 是 OuterIterator 的抽象实现，子类需覆盖 accept() 方法
// 状态存储在 ClassValue 属性中，通过 ctx 访问，不存在于 Go struct
type FilterIteratorClass struct {
	node.Node
}

// NewFilterIteratorClass 创建 FilterIteratorClass 实例
func NewFilterIteratorClass() *FilterIteratorClass {
	return &FilterIteratorClass{}
}

// GetValue 每次 new 时创建地址的副本（Class 类型不变，状态由 ClassValue.ObjectValue 属性承载）
func (c *FilterIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	// 初始化内部状态属性
	cv.SetProperty(filterValidKey, data.NewBoolValue(false))
	cv.SetProperty(filterCurValKey, data.NewNullValue())
	cv.SetProperty(filterCurKeyKey, data.NewNullValue())
	return cv, nil
}

func (c *FilterIteratorClass) GetName() string { return "FilterIterator" }

func (c *FilterIteratorClass) GetExtend() *string { return nil }

func (c *FilterIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator"}
}

func (c *FilterIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *FilterIteratorClass) GetPropertyList() []data.Property { return nil }

func (c *FilterIteratorClass) GetConstruct() data.Method {
	return &FilterIteratorConstructMethod{}
}

func (c *FilterIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &FilterIteratorConstructMethod{}, true
	case "rewind":
		return &FilterIteratorRewindMethod{}, true
	case "valid":
		return &FilterIteratorValidMethod{}, true
	case "current":
		return &FilterIteratorCurrentMethod{}, true
	case "key":
		return &FilterIteratorKeyMethod{}, true
	case "next":
		return &FilterIteratorNextMethod{}, true
	case "accept":
		return &FilterIteratorAcceptMethod{}, true
	case "getInnerIterator":
		return &FilterIteratorGetInnerIteratorMethod{}, true
	}
	return nil, false
}

func (c *FilterIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&FilterIteratorConstructMethod{},
		&FilterIteratorRewindMethod{},
		&FilterIteratorValidMethod{},
		&FilterIteratorCurrentMethod{},
		&FilterIteratorKeyMethod{},
		&FilterIteratorNextMethod{},
		&FilterIteratorAcceptMethod{},
		&FilterIteratorGetInnerIteratorMethod{},
	}
}

// ---- 从 ctx 中获取当前 ClassValue（支持 ClassMethodContext 和 ClassValue两种入口） ----

func filterGetClassValue(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

// ---- 状态读写辅助函数（通过 ClassValue.ObjectValue 属性） ----

func filterGetInner(cv *data.ClassValue) data.GetValue {
	v, _ := cv.ObjectValue.GetProperty(filterInnerKey)
	if v == nil {
		return nil
	}
	if _, ok := v.(*data.NullValue); ok {
		return nil
	}
	return v
}

func filterSetInner(cv *data.ClassValue, inner data.GetValue) {
	if inner == nil {
		cv.ObjectValue.SetProperty(filterInnerKey, data.NewNullValue())
		return
	}
	if v, ok := inner.(data.Value); ok {
		cv.ObjectValue.SetProperty(filterInnerKey, v)
	}
}

func filterIsValid(cv *data.ClassValue) bool {
	v, _ := cv.ObjectValue.GetProperty(filterValidKey)
	if bv, ok := v.(*data.BoolValue); ok {
		return bv.Value
	}
	return false
}

func filterSetValid(cv *data.ClassValue, valid bool) {
	cv.ObjectValue.SetProperty(filterValidKey, data.NewBoolValue(valid))
}

func filterGetCurVal(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(filterCurValKey)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func filterSetCurVal(cv *data.ClassValue, val data.Value) {
	if val == nil {
		val = data.NewNullValue()
	}
	cv.ObjectValue.SetProperty(filterCurValKey, val)
}

func filterGetCurKey(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(filterCurKeyKey)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func filterSetCurKey(cv *data.ClassValue, key data.Value) {
	if key == nil {
		key = data.NewNullValue()
	}
	cv.ObjectValue.SetProperty(filterCurKeyKey, key)
}

// ---- 内部迭代器方法调用辅助函数 ----

// callInnerMethod 调用内部迭代器的某个方法
// 必须基于内部迭代器自身的 ClassValue 创建 context，不得使用外层 ctx
func filterCallInnerMethod(inner data.GetValue, name string) (data.GetValue, data.Control) {
	if inner == nil {
		return nil, nil
	}
	// 展开 ThisValue
	if tv, ok := inner.(*data.ThisValue); ok {
		inner = tv.ClassValue
	}
	if cv, ok := inner.(*data.ClassValue); ok {
		if method, found := cv.GetMethod(name); found {
			innerCtx := cv.CreateContext(method.GetVariables())
			return method.Call(innerCtx)
		}
	}
	return nil, nil
}

func filterInnerValid(inner data.GetValue) bool {
	result, _ := filterCallInnerMethod(inner, "valid")
	if bv, ok := result.(*data.BoolValue); ok {
		return bv.Value
	}
	return false
}

func filterInnerCurrent(inner data.GetValue) data.Value {
	result, _ := filterCallInnerMethod(inner, "current")
	if v, ok := result.(data.Value); ok {
		return v
	}
	return data.NewNullValue()
}

func filterInnerKeyVal(inner data.GetValue) data.Value {
	result, _ := filterCallInnerMethod(inner, "key")
	if v, ok := result.(data.Value); ok {
		return v
	}
	return data.NewNullValue()
}

// ---- callAccept：通过 ctx 中的 ClassValue 动态调用子类的 accept ----

// callAccept 通过 ctx 内的 ClassValue 动态查找 accept 方法（支持 PHP 子类覆盖）
func filterCallAccept(ctx data.Context) bool {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return true
	}
	if method, found := cv.GetMethod("accept"); found {
		acceptCtx := cv.CreateContext(method.GetVariables())
		result, _ := method.Call(acceptCtx)
		if bv, ok := result.(*data.BoolValue); ok {
			return bv.Value
		}
		if v, ok := result.(data.AsBool); ok {
			b, _ := v.AsBool()
			return b
		}
	}
	return true
}

// advanceToAccepted 从内部迭代器当前位置向前推进，直到找到第一个通过 accept 的位置
// 必须先更新 currentValue 缓存，再调用 accept（PHP 子类 accept() 内用 $this->current() 需要正确内部迭代器当前值）
func advanceToAccepted(cv *data.ClassValue, ctx data.Context) {
	inner := filterGetInner(cv)
	for filterInnerValid(inner) {
		// 先更新当前值缓存，保证 accept() 内 $this->current() 获取到正确内部迭代器当前值
		filterSetCurVal(cv, filterInnerCurrent(inner))
		filterSetCurKey(cv, filterInnerKeyVal(inner))
		if filterCallAccept(ctx) {
			filterSetValid(cv, true)
			return
		}
		filterCallInnerMethod(inner, "next")
	}
	// 内部迭代器耗尽
	filterSetValid(cv, false)
	filterSetCurVal(cv, data.NewNullValue())
	filterSetCurKey(cv, data.NewNullValue())
}

// ---- __construct($iterator) ----

type FilterIteratorConstructMethod struct{}

func (m *FilterIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *FilterIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilterIteratorConstructMethod) GetIsStatic() bool          { return false }
func (m *FilterIteratorConstructMethod) GetReturnType() data.Types  { return nil }
func (m *FilterIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
	}
}
func (m *FilterIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
	}
}
func (m *FilterIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	filterSetInner(cv, it)
	filterSetValid(cv, false)
	filterSetCurVal(cv, data.NewNullValue())
	filterSetCurKey(cv, data.NewNullValue())
	return nil, nil
}

// ---- rewind() ----

type FilterIteratorRewindMethod struct{}

func (m *FilterIteratorRewindMethod) GetName() string               { return "rewind" }
func (m *FilterIteratorRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorRewindMethod) GetIsStatic() bool             { return false }
func (m *FilterIteratorRewindMethod) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorRewindMethod) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	inner := filterGetInner(cv)
	filterCallInnerMethod(inner, "rewind")
	advanceToAccepted(cv, ctx)
	return nil, nil
}

// ---- valid() ----

type FilterIteratorValidMethod struct{}

func (m *FilterIteratorValidMethod) GetName() string               { return "valid" }
func (m *FilterIteratorValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorValidMethod) GetIsStatic() bool             { return false }
func (m *FilterIteratorValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *FilterIteratorValidMethod) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorValidMethod) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(filterIsValid(cv)), nil
}

// ---- current() ----

type FilterIteratorCurrentMethod struct{}

func (m *FilterIteratorCurrentMethod) GetName() string               { return "current" }
func (m *FilterIteratorCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorCurrentMethod) GetIsStatic() bool             { return false }
func (m *FilterIteratorCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *FilterIteratorCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return filterGetCurVal(cv), nil
}

// ---- key() ----

type FilterIteratorKeyMethod struct{}

func (m *FilterIteratorKeyMethod) GetName() string               { return "key" }
func (m *FilterIteratorKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorKeyMethod) GetIsStatic() bool             { return false }
func (m *FilterIteratorKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *FilterIteratorKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorKeyMethod) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return filterGetCurKey(cv), nil
}

// ---- next() ----

type FilterIteratorNextMethod struct{}

func (m *FilterIteratorNextMethod) GetName() string               { return "next" }
func (m *FilterIteratorNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorNextMethod) GetIsStatic() bool             { return false }
func (m *FilterIteratorNextMethod) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorNextMethod) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorNextMethod) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	inner := filterGetInner(cv)
	filterCallInnerMethod(inner, "next")
	advanceToAccepted(cv, ctx)
	return nil, nil
}

// ---- accept(): 抽象方法，默认返回 true，PHP 子类覆盖 ----

type FilterIteratorAcceptMethod struct{}

func (m *FilterIteratorAcceptMethod) GetName() string               { return "accept" }
func (m *FilterIteratorAcceptMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FilterIteratorAcceptMethod) GetIsStatic() bool             { return false }
func (m *FilterIteratorAcceptMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *FilterIteratorAcceptMethod) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorAcceptMethod) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorAcceptMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 抽象方法默认实现：返回 true（真实逻辑由 PHP 子类覆盖）
	return data.NewBoolValue(true), nil
}

// ---- getInnerIterator() ----

type FilterIteratorGetInnerIteratorMethod struct{}

func (m *FilterIteratorGetInnerIteratorMethod) GetName() string { return "getInnerIterator" }
func (m *FilterIteratorGetInnerIteratorMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *FilterIteratorGetInnerIteratorMethod) GetIsStatic() bool             { return false }
func (m *FilterIteratorGetInnerIteratorMethod) GetReturnType() data.Types     { return nil }
func (m *FilterIteratorGetInnerIteratorMethod) GetParams() []data.GetValue    { return nil }
func (m *FilterIteratorGetInnerIteratorMethod) GetVariables() []data.Variable { return nil }
func (m *FilterIteratorGetInnerIteratorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := filterGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	inner := filterGetInner(cv)
	if inner == nil {
		return data.NewNullValue(), nil
	}
	return inner, nil
}
