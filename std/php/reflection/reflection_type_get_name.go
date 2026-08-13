package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionTypeGetNameMethod 实现 ReflectionType::getName
// 返回被反射类型的名称
type ReflectionTypeGetNameMethod struct{}

// GetName 返回方法名 "getName"
func (m *ReflectionTypeGetNameMethod) GetName() string { return "getName" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionTypeGetNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionTypeGetNameMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionTypeGetNameMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionTypeGetNameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回字符串类型
func (m *ReflectionTypeGetNameMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

// Call 执行 getName 方法
// 返回被反射类型的名称
func (m *ReflectionTypeGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	typeName, _ := getReflectionTypeInfo(ctx)
	return data.NewStringValue(typeName), nil
}
