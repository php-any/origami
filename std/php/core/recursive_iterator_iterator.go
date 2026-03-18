package core

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RecursiveIteratorIteratorClass 实现 PHP 的 RecursiveIteratorIterator 类
type RecursiveIteratorIteratorClass struct {
	node.Node
	innerIterator  data.GetValue         // 内部迭代器
	mode           int                   // 迭代模式
	currentKey     data.Value            // 当前键
	currentValue   data.Value            // 当前值
	valid          bool                  // 是否有效
	StaticProperty map[string]data.Value // 静态属性（类常量）
}

func NewRecursiveIteratorIteratorClass() *RecursiveIteratorIteratorClass {
	return &RecursiveIteratorIteratorClass{
		innerIterator: nil,
		mode:          0,
		currentKey:    data.NewStringValue(""),
		currentValue:  nil,
		valid:         false,
		StaticProperty: map[string]data.Value{
			"SELF_FIRST":       data.NewIntValue(1),
			"CHILD_FIRST":      data.NewIntValue(2),
			"LEAVES_ONLY":      data.NewIntValue(0),
			"SELF_FIRST_SELF":  data.NewIntValue(4),
			"CHILD_FIRST_SELF": data.NewIntValue(8),
		},
	}
}

func (r *RecursiveIteratorIteratorClass) GetName() string {
	return "RecursiveIteratorIterator"
}

func (r *RecursiveIteratorIteratorClass) GetExtend() *string {
	return nil
}

func (r *RecursiveIteratorIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator"}
}

func (r *RecursiveIteratorIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (r *RecursiveIteratorIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := r.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}

func (r *RecursiveIteratorIteratorClass) GetPropertyList() []data.Property {
	return nil
}

func (r *RecursiveIteratorIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建新的实例，每个实例有自己的状态（遵循技能文档规范）
	clone := &RecursiveIteratorIteratorClass{
		innerIterator:  r.innerIterator,
		mode:           r.mode,
		currentKey:     r.currentKey,
		currentValue:   r.currentValue,
		valid:          r.valid,
		StaticProperty: r.StaticProperty,
	}
	return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}

func (r *RecursiveIteratorIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RecursiveIteratorIteratorConstruct{instance: r}, true
	case "rewind":
		return &RecursiveIteratorIteratorRewind{instance: r}, true
	case "current":
		return &RecursiveIteratorIteratorCurrent{instance: r}, true
	case "key":
		return &RecursiveIteratorIteratorKey{instance: r}, true
	case "next":
		return &RecursiveIteratorIteratorNext{instance: r}, true
	case "valid":
		return &RecursiveIteratorIteratorValid{instance: r}, true
	case "getInnerIterator":
		return &RecursiveIteratorIteratorGetInnerIterator{instance: r}, true
	}
	return nil, false
}

func (r *RecursiveIteratorIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RecursiveIteratorIteratorConstruct{instance: r},
		&RecursiveIteratorIteratorRewind{instance: r},
		&RecursiveIteratorIteratorCurrent{instance: r},
		&RecursiveIteratorIteratorKey{instance: r},
		&RecursiveIteratorIteratorNext{instance: r},
		&RecursiveIteratorIteratorValid{instance: r},
		&RecursiveIteratorIteratorGetInnerIterator{instance: r},
	}
}

func (r *RecursiveIteratorIteratorClass) GetConstruct() data.Method {
	return &RecursiveIteratorIteratorConstruct{instance: r}
}

// RecursiveIteratorIteratorConstruct 构造函数
type RecursiveIteratorIteratorConstruct struct {
	instance *RecursiveIteratorIteratorClass
}

func (m *RecursiveIteratorIteratorConstruct) GetName() string {
	return "__construct"
}

func (m *RecursiveIteratorIteratorConstruct) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorIteratorConstruct) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorIteratorConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("iterator", 0, data.NewBaseType("Traversable")),
		data.NewVariable("mode", 1, data.NewBaseType("int")),
		data.NewVariable("flags", 2, data.NewBaseType("int")),
	}
}

func (m *RecursiveIteratorIteratorConstruct) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorIteratorConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Traversable")),
		node.NewParameter(nil, "mode", 1, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "flags", 2, nil, data.NewBaseType("int")),
	}
}

func (m *RecursiveIteratorIteratorConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取内部迭代器
	iterator, exists := ctx.GetIndexValue(0)
	if !exists {
		return nil, data.NewErrorThrow(nil, errors.New("缺少必需的迭代器参数"))
	}

	// 直接设置实例状态
	m.instance.innerIterator = iterator
	m.instance.mode = 0
	m.instance.currentKey = data.NewStringValue("")
	m.instance.currentValue = nil
	m.instance.valid = false

	return nil, nil
}

// RecursiveIteratorIteratorRewind 重置迭代器
type RecursiveIteratorIteratorRewind struct {
	instance *RecursiveIteratorIteratorClass
}

func (m *RecursiveIteratorIteratorRewind) GetName() string {
	return "rewind"
}

func (m *RecursiveIteratorIteratorRewind) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorIteratorRewind) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorIteratorRewind) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorIteratorRewind) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorIteratorRewind) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorIteratorRewind) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 调用内部迭代器的 rewind
	if innerIt, ok := m.instance.innerIterator.(interface {
		GetMethod(string) (data.Method, bool)
	}); ok {
		if method, found := innerIt.GetMethod("rewind"); found {
			method.Call(ctx)
		}
	}

	// 开始递归遍历
	m.beginRecursiveIteration(ctx)

	return nil, nil
}

func (m *RecursiveIteratorIteratorRewind) beginRecursiveIteration(ctx data.Context) {
	// 检查内部迭代器是否有效
	if innerIt, ok := m.instance.innerIterator.(interface {
		GetMethod(string) (data.Method, bool)
	}); ok {
		if validMethod, found := innerIt.GetMethod("valid"); found {
			validResult, _ := validMethod.Call(ctx)
			if validBool, ok := validResult.(*data.BoolValue); ok {
				fmt.Printf("Internal iterator valid: %v\n", validBool.Value)
				if validBool.Value {
					// 获取当前值
					if currentMethod, found := innerIt.GetMethod("current"); found {
						currentVal, ctrl := currentMethod.Call(ctx)
						if ctrl != nil {
							fmt.Printf("Current method error: %v\n", ctrl)
							return
						}
						fmt.Printf("Current value: %v\n", currentVal)
						// 将 GetValue 转换为 Value
						var valueData data.Value
						if v, ok := currentVal.(data.Value); ok {
							valueData = v
						} else {
							valueData = data.NewNullValue()
						}
						m.instance.currentValue = valueData

						if keyMethod, found := innerIt.GetMethod("key"); found {
							keyVal, _ := keyMethod.Call(ctx)
							fmt.Printf("Key value: %v\n", keyVal)
							var keyData data.Value
							if k, ok := keyVal.(data.Value); ok {
								keyData = k
							} else {
								keyData = data.NewStringValue("")
							}
							m.instance.currentKey = keyData
						}

						m.instance.valid = true
						fmt.Printf("Set valid to true\n")
					} else {
						m.instance.valid = false
						fmt.Printf("Set valid to false - no current method\n")
					}
				} else {
					m.instance.valid = false
					fmt.Printf("Set valid to false - internal iterator not valid\n")
				}
			} else {
				fmt.Printf("Cannot convert valid result to BoolValue\n")
				m.instance.valid = false
			}
		} else {
			fmt.Printf("Internal iterator has no valid method\n")
			m.instance.valid = false
		}
	} else {
		fmt.Printf("Cannot cast inner iterator to method interface\n")
		m.instance.valid = false
	}
}

// RecursiveIteratorIteratorCurrent 返回当前元素
type RecursiveIteratorIteratorCurrent struct {
	instance *RecursiveIteratorIteratorClass
}

func (m *RecursiveIteratorIteratorCurrent) GetName() string {
	return "current"
}

func (m *RecursiveIteratorIteratorCurrent) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorIteratorCurrent) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorIteratorCurrent) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorIteratorCurrent) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorIteratorCurrent) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorIteratorCurrent) Call(ctx data.Context) (data.GetValue, data.Control) {
	// PHP 的 RecursiveIteratorIterator::current() 返回迭代器对象本身
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return objCtx.ClassValue, nil
	}
	// 回退到返回存储的当前值
	return m.instance.currentValue, nil
}

// RecursiveIteratorIteratorKey 返回当前键
type RecursiveIteratorIteratorKey struct {
	instance *RecursiveIteratorIteratorClass
}

func (m *RecursiveIteratorIteratorKey) GetName() string {
	return "key"
}

func (m *RecursiveIteratorIteratorKey) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorIteratorKey) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorIteratorKey) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorIteratorKey) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorIteratorKey) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorIteratorKey) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.instance.currentKey, nil
}

// RecursiveIteratorIteratorNext 移动到下一个元素
type RecursiveIteratorIteratorNext struct {
	instance *RecursiveIteratorIteratorClass
}

func (m *RecursiveIteratorIteratorNext) GetName() string {
	return "next"
}

func (m *RecursiveIteratorIteratorNext) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorIteratorNext) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorIteratorNext) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorIteratorNext) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorIteratorNext) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorIteratorNext) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 调用内部迭代器的 next
	if innerIt, ok := m.instance.innerIterator.(interface {
		GetMethod(string) (data.Method, bool)
	}); ok {
		if method, found := innerIt.GetMethod("next"); found {
			method.Call(ctx)
		}

		// 更新当前值
		rewind := &RecursiveIteratorIteratorRewind{instance: m.instance}
		rewind.beginRecursiveIteration(ctx)
	}

	return nil, nil
}

// RecursiveIteratorIteratorValid 检查当前位置是否有效
type RecursiveIteratorIteratorValid struct {
	instance *RecursiveIteratorIteratorClass
}

func (m *RecursiveIteratorIteratorValid) GetName() string {
	return "valid"
}

func (m *RecursiveIteratorIteratorValid) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorIteratorValid) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorIteratorValid) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorIteratorValid) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorIteratorValid) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorIteratorValid) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(m.instance.valid), nil
}

// RecursiveIteratorIteratorGetInnerIterator 返回内部迭代器
type RecursiveIteratorIteratorGetInnerIterator struct {
	instance *RecursiveIteratorIteratorClass
}

func (m *RecursiveIteratorIteratorGetInnerIterator) GetName() string {
	return "getInnerIterator"
}

func (m *RecursiveIteratorIteratorGetInnerIterator) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorIteratorGetInnerIterator) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorIteratorGetInnerIterator) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorIteratorGetInnerIterator) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorIteratorGetInnerIterator) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorIteratorGetInnerIterator) Call(ctx data.Context) (data.GetValue, data.Control) {
	fmt.Printf("getInnerIterator returning: %v\n", m.instance.innerIterator)
	return m.instance.innerIterator, nil
}
