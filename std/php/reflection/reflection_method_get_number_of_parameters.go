package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodGetNumberOfParametersMethod 实现 ReflectionMethod::getNumberOfParameters
// 返回被反射方法的参数数量
type ReflectionMethodGetNumberOfParametersMethod struct{}

// GetName 返回方法名 "getNumberOfParameters"
func (m *ReflectionMethodGetNumberOfParametersMethod) GetName() string {
	return "getNumberOfParameters"
}

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodGetNumberOfParametersMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodGetNumberOfParametersMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodGetNumberOfParametersMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodGetNumberOfParametersMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回整数类型
func (m *ReflectionMethodGetNumberOfParametersMethod) GetReturnType() data.Types {
	return data.Int{}
}

// Call 执行 getNumberOfParameters 方法
// 返回被反射方法的参数数量
func (m *ReflectionMethodGetNumberOfParametersMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewIntValue(0), nil
	}

	params := method.GetParams()
	return data.NewIntValue(len(params)), nil
}
