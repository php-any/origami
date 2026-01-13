package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodGetDeclaringClassMethod 实现 ReflectionMethod::getDeclaringClass
// 返回声明该方法的类（ReflectionClass 对象）
type ReflectionMethodGetDeclaringClassMethod struct{}

// GetName 返回方法名 "getDeclaringClass"
func (m *ReflectionMethodGetDeclaringClassMethod) GetName() string { return "getDeclaringClass" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodGetDeclaringClassMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodGetDeclaringClassMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodGetDeclaringClassMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodGetDeclaringClassMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回 ReflectionClass 类型
func (m *ReflectionMethodGetDeclaringClassMethod) GetReturnType() data.Types {
	return data.Mixed{}
}

// Call 执行 getDeclaringClass 方法
// 返回声明该方法的类，返回 ReflectionClass 对象
// 如果无法获取类信息，返回 null
func (m *ReflectionMethodGetDeclaringClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	declaringClassName, _, _ := getReflectionMethodInfo(ctx)
	if declaringClassName == "" {
		return data.NewNullValue(), nil
	}

	// 创建 ReflectionClass 实例
	classClass := &ReflectionClassClass{}
	classValue := data.NewClassValue(classClass, ctx.CreateBaseContext())

	// 存储类名到实例属性中
	classValue.ObjectValue.SetProperty("_className", data.NewStringValue(declaringClassName))

	return classValue, nil
}
