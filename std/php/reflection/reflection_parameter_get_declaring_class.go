package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionParameterGetDeclaringClassMethod 实现 ReflectionParameter::getDeclaringClass
// 返回声明该参数的类（即包含该参数的方法所属的类）
type ReflectionParameterGetDeclaringClassMethod struct{}

// GetName 返回方法名 "getDeclaringClass"
func (m *ReflectionParameterGetDeclaringClassMethod) GetName() string { return "getDeclaringClass" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterGetDeclaringClassMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterGetDeclaringClassMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionParameterGetDeclaringClassMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionParameterGetDeclaringClassMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回 ReflectionClass 对象（或 null）
func (m *ReflectionParameterGetDeclaringClassMethod) GetReturnType() data.Types {
	return data.NewBaseType("?ReflectionClass")
}

// Call 执行 getDeclaringClass 方法
// 返回声明该参数的类，返回 ReflectionClass 对象
// 如果无法获取类信息，返回 null
func (m *ReflectionParameterGetDeclaringClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, _, _, _ := getReflectionParameterInfo(ctx)
	if className == "" {
		return data.NewNullValue(), nil
	}

	// 创建 ReflectionClass 实例
	reflectionClass := &ReflectionClassClass{}
	reflectionClassValue := data.NewClassValue(reflectionClass, ctx.CreateBaseContext())

	// 存储类名到实例属性中
	reflectionClassValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))

	return reflectionClassValue, nil
}
