package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassIsInstanceMethod 实现 ReflectionClass::isInstance
// 检查指定的对象是否是被反射类的实例或其子类的实例
type ReflectionClassIsInstanceMethod struct{}

// GetName 返回方法名 "isInstance"
func (m *ReflectionClassIsInstanceMethod) GetName() string { return "isInstance" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassIsInstanceMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassIsInstanceMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - object: 要检查的对象实例
func (m *ReflectionClassIsInstanceMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassIsInstanceMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
	}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionClassIsInstanceMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 isInstance 方法
// 检查指定的对象是否是被反射类的实例或其子类的实例
// 通过检查对象的类名和继承关系来判断
// 返回 true 表示是实例，false 表示不是
func (m *ReflectionClassIsInstanceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	objectValue, _ := ctx.GetIndexValue(0)
	if objectValue == nil {
		return data.NewBoolValue(false), nil
	}

	classValue, ok := objectValue.(*data.ClassValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	className, _ := getReflectionClassInfo(ctx)
	objectClassName := classValue.Class.GetName()

	// 检查是否是同一个类或者是子类
	if objectClassName == className {
		return data.NewBoolValue(true), nil
	}

	// 检查继承关系
	vm := ctx.GetVM()
	last := classValue.Class
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		if *ext == className {
			return data.NewBoolValue(true), nil
		}
		next, ok := vm.GetClass(*ext)
		if !ok {
			break
		}
		last = next
	}

	return data.NewBoolValue(false), nil
}
