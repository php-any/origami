package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodIsProtectedMethod 实现 ReflectionMethod::isProtected
// 检查被反射的方法是否为受保护方法
type ReflectionMethodIsProtectedMethod struct{}

// GetName 返回方法名 "isProtected"
func (m *ReflectionMethodIsProtectedMethod) GetName() string { return "isProtected" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodIsProtectedMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodIsProtectedMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodIsProtectedMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodIsProtectedMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionMethodIsProtectedMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 isProtected 方法
// 检查被反射的方法是否为受保护方法
func (m *ReflectionMethodIsProtectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(method.GetModifier() == data.ModifierProtected), nil
}
