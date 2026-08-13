package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodGetParametersMethod 实现 ReflectionMethod::getParameters
// 返回被反射方法的参数列表
type ReflectionMethodGetParametersMethod struct{}

// GetName 返回方法名 "getParameters"
func (m *ReflectionMethodGetParametersMethod) GetName() string { return "getParameters" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodGetParametersMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodGetParametersMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodGetParametersMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodGetParametersMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回数组类型
func (m *ReflectionMethodGetParametersMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

// Call 执行 getParameters 方法
// 返回被反射方法的参数列表，返回 ReflectionParameter 对象数组
func (m *ReflectionMethodGetParametersMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, methodName, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	params := method.GetParams()
	result := make([]data.Value, 0, len(params))

	for index, param := range params {
		if param == nil {
			continue
		}

		// 创建 ReflectionParameter 对象
		paramObj := newReflectionParameter(ctx, className, methodName, index, param)
		result = append(result, paramObj)
	}

	return data.NewArrayValue(result), nil
}
