package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionParameterGetPositionMethod 实现 ReflectionParameter::getPosition
// 返回被反射参数的位置（索引）
type ReflectionParameterGetPositionMethod struct{}

// GetName 返回方法名 "getPosition"
func (m *ReflectionParameterGetPositionMethod) GetName() string { return "getPosition" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterGetPositionMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterGetPositionMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionParameterGetPositionMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionParameterGetPositionMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回整数类型
func (m *ReflectionParameterGetPositionMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

// Call 执行 getPosition 方法
// 返回被反射参数的位置（索引）
func (m *ReflectionParameterGetPositionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, paramIndex, _ := getReflectionParameterInfo(ctx)
	return data.NewIntValue(paramIndex), nil
}
