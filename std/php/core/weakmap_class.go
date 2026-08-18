package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// WeakMapClass 表示 PHP 8.0+ 的 WeakMap 类
// WeakMap 允许创建对象的弱引用映射，不会阻止对象被垃圾回收
// 在 Origami 中，使用同步 map 实现对象键映射
type WeakMapClass struct {
	node.Node
	mu    sync.RWMutex
	store map[string]data.Value // 使用对象 ID 作为键
}

func NewWeakMapClass() *WeakMapClass {
	return &WeakMapClass{
		store: make(map[string]data.Value),
	}
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
	return data.NewProxyValue(NewWeakMapClass(), ctx.CreateBaseContext()), nil
}

func (w *WeakMapClass) GetSource() any { return w }

func objectKey(v data.Value) string {
	switch t := v.(type) {
	case *data.ClassValue:
		if t.ObjectValue != nil {
			return fmt.Sprintf("obj:%p", t.ObjectValue)
		}
		return "class:" + t.Class.GetName()
	case *data.ThisValue:
		if t.ClassValue != nil && t.ClassValue.ObjectValue != nil {
			return fmt.Sprintf("obj:%p", t.ClassValue.ObjectValue)
		}
		return "class:" + t.ClassValue.Class.GetName()
	case *data.StringValue:
		return "str:" + t.Value
	default:
		return "val:" + v.AsString()
	}
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
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
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
	objVal, has := ctx.GetIndexValue(0)
	if !has {
		return data.NewBoolValue(false), nil
	}

	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || classCtx.ClassValue == nil {
		return data.NewBoolValue(false), nil
	}

	wm, ok := classCtx.ClassValue.Class.(*WeakMapClass)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	wm.mu.RLock()
	defer wm.mu.RUnlock()
	_, exists := wm.store[objectKey(objVal)]
	return data.NewBoolValue(exists), nil
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
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
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
	objVal, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("InvalidArgumentException: Object not found in WeakMap"))
	}

	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || classCtx.ClassValue == nil {
		return nil, utils.NewThrow(errors.New("InvalidArgumentException: Object not found in WeakMap"))
	}

	wm, ok := classCtx.ClassValue.Class.(*WeakMapClass)
	if !ok {
		return nil, utils.NewThrow(errors.New("InvalidArgumentException: Object not found in WeakMap"))
	}

	wm.mu.RLock()
	val, exists := wm.store[objectKey(objVal)]
	wm.mu.RUnlock()

	if !exists {
		return nil, utils.NewThrow(errors.New("InvalidArgumentException: Object not found in WeakMap"))
	}
	return val, nil
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
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
		node.NewVariable(nil, "value", 1, data.NewBaseType("mixed")),
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
	objVal, has1 := ctx.GetIndexValue(0)
	val, has2 := ctx.GetIndexValue(1)
	if !has1 || !has2 {
		return nil, nil
	}

	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || classCtx.ClassValue == nil {
		return nil, nil
	}

	wm, ok := classCtx.ClassValue.Class.(*WeakMapClass)
	if !ok {
		return nil, nil
	}

	wm.mu.Lock()
	wm.store[objectKey(objVal)] = val
	wm.mu.Unlock()
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
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
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
	objVal, has := ctx.GetIndexValue(0)
	if !has {
		return nil, nil
	}

	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || classCtx.ClassValue == nil {
		return nil, nil
	}

	wm, ok := classCtx.ClassValue.Class.(*WeakMapClass)
	if !ok {
		return nil, nil
	}

	wm.mu.Lock()
	delete(wm.store, objectKey(objVal))
	wm.mu.Unlock()
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
	return nil
}

func (m *WeakMapCountMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

func (m *WeakMapCountMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *WeakMapCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || classCtx.ClassValue == nil {
		return data.NewIntValue(0), nil
	}

	wm, ok := classCtx.ClassValue.Class.(*WeakMapClass)
	if !ok {
		return data.NewIntValue(0), nil
	}

	wm.mu.RLock()
	count := len(wm.store)
	wm.mu.RUnlock()
	return data.NewIntValue(count), nil
}
