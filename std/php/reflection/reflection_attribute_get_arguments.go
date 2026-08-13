package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionAttributeGetArgumentsMethod 实现 ReflectionAttribute::getArguments
// 返回传递给属性（注解）的参数数组
type ReflectionAttributeGetArgumentsMethod struct{}

// GetName 返回方法名 "getArguments"
func (m *ReflectionAttributeGetArgumentsMethod) GetName() string { return "getArguments" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionAttributeGetArgumentsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionAttributeGetArgumentsMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionAttributeGetArgumentsMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionAttributeGetArgumentsMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回数组类型
func (m *ReflectionAttributeGetArgumentsMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

// Call 执行 getArguments 方法
// 返回传递给属性（注解）的参数数组
func (m *ReflectionAttributeGetArgumentsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	annotationValue := getReflectionAttributeInfo(ctx)
	if annotationValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 注解对象是 ClassValue，但我们已经实例化了注解
	// 原始参数已经被传递给构造函数并使用了
	// 由于 ClassStatement.Annotations 只存储了 ClassValue，我们无法直接获取原始参数
	// 这里返回空数组，因为参数已经在实例化时被使用了
	// TODO: 如果需要支持 getArguments，需要在 ClassStatement 中同时存储 Annotation 节点
	return data.NewArrayValue([]data.Value{}), nil
}
