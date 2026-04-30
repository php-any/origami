package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// WeakMapClass 表示 PHP 8.0+ 的 WeakMap 类
// WeakMap 允许创建对象的弱引用映射，不会阻止对象被垃圾回收
type WeakMapClass struct {
	node.Node
}

func NewWeakMapClass() *WeakMapClass {
	return &WeakMapClass{}
}

func (w *WeakMapClass) GetName() string {
	return "WeakMap"
}

func (w *WeakMapClass) GetExtend() *string {
	return nil
}

func (w *WeakMapClass) GetImplements() []string {
	return []string{"ArrayAccess", "Countable"}
}

func (w *WeakMapClass) GetMethods() []data.Method {
	return []data.Method{
		&WeakMapOffsetExistsMethod{},
		&WeakMapOffsetGetMethod{},
		&WeakMapOffsetSetMethod{},
		&WeakMapOffsetUnsetMethod{},
		&WeakMapCountMethod{},
	}
}

func (w *WeakMapClass) GetMethod(name string) (data.Method, bool) {
	methods := w.GetMethods()
	for _, method := range methods {
		if method.GetName() == name {
			return method, true
		}
	}
	return nil, false
}

func (w *WeakMapClass) GetConstruct() data.Method {
	return nil
}

func (w *WeakMapClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (w *WeakMapClass) GetPropertyList() []data.Property {
	return nil
}

func (w *WeakMapClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(w, ctx.CreateBaseContext()), nil
}

// WeakMapOffsetExistsMethod 实现 offsetExists 方法
type WeakMapOffsetExistsMethod struct{}

func (m *WeakMapOffsetExistsMethod) GetName() string {
	return "offsetExists"
}

func (m *WeakMapOffsetExistsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *WeakMapOffsetExistsMethod) GetIsStatic() bool {
	return false
}

func (m *WeakMapOffsetExistsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, nil),
	}
}

func (m *WeakMapOffsetExistsMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *WeakMapOffsetExistsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
	}
}

func (m *WeakMapOffsetExistsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简单实现，返回 false
	// 实际需要使用真正的弱引用实现
	return data.NewBoolValue(false), nil
}

// WeakMapOffsetGetMethod 实现 offsetGet 方法
type WeakMapOffsetGetMethod struct{}

func (m *WeakMapOffsetGetMethod) GetName() string {
	return "offsetGet"
}

func (m *WeakMapOffsetGetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *WeakMapOffsetGetMethod) GetIsStatic() bool {
	return false
}

func (m *WeakMapOffsetGetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, nil),
	}
}

func (m *WeakMapOffsetGetMethod) GetReturnType() data.Types {
	return data.NewBaseType("mixed")
}

func (m *WeakMapOffsetGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
	}
}

func (m *WeakMapOffsetGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 抛出异常，因为对象不存在
	panic("InvalidArgumentException: Object not found in WeakMap")
	return nil, nil
}

// WeakMapOffsetSetMethod 实现 offsetSet 方法
type WeakMapOffsetSetMethod struct{}

func (m *WeakMapOffsetSetMethod) GetName() string {
	return "offsetSet"
}

func (m *WeakMapOffsetSetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *WeakMapOffsetSetMethod) GetIsStatic() bool {
	return false
}

func (m *WeakMapOffsetSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, nil),
		node.NewVariable(nil, "value", 1, nil),
	}
}

func (m *WeakMapOffsetSetMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *WeakMapOffsetSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
		node.NewParameter(nil, "value", 1, nil, data.NewBaseType("mixed")),
	}
}

func (m *WeakMapOffsetSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简单实现，不存储任何内容
	return nil, nil
}

// WeakMapOffsetUnsetMethod 实现 offsetUnset 方法
type WeakMapOffsetUnsetMethod struct{}

func (m *WeakMapOffsetUnsetMethod) GetName() string {
	return "offsetUnset"
}

func (m *WeakMapOffsetUnsetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *WeakMapOffsetUnsetMethod) GetIsStatic() bool {
	return false
}

func (m *WeakMapOffsetUnsetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, nil),
	}
}

func (m *WeakMapOffsetUnsetMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *WeakMapOffsetUnsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
	}
}

func (m *WeakMapOffsetUnsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简单实现，不做任何操作
	return nil, nil
}

// WeakMapCountMethod 实现 count 方法
type WeakMapCountMethod struct{}

func (m *WeakMapCountMethod) GetName() string {
	return "count"
}

func (m *WeakMapCountMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *WeakMapCountMethod) GetIsStatic() bool {
	return false
}

func (m *WeakMapCountMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *WeakMapCountMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

func (m *WeakMapCountMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *WeakMapCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 返回 0，因为没有存储任何内容
	return data.NewIntValue(0), nil
}
