package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionTypeToStringMethod 实现 ReflectionType::__toString
// 返回被反射类型的字符串表示
type ReflectionTypeToStringMethod struct{}

// GetName 返回方法名 "__toString"
func (m *ReflectionTypeToStringMethod) GetName() string { return "__toString" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionTypeToStringMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionTypeToStringMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionTypeToStringMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionTypeToStringMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回字符串类型
func (m *ReflectionTypeToStringMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

// Call 执行 __toString 方法
// 返回被反射类型的字符串表示
func (m *ReflectionTypeToStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	typeName, _ := getReflectionTypeInfo(ctx)
	return data.NewStringValue(typeName), nil
}
