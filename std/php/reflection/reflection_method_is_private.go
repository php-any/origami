package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodIsPrivateMethod 实现 ReflectionMethod::isPrivate
// 检查被反射的方法是否为私有方法
type ReflectionMethodIsPrivateMethod struct{}

// GetName 返回方法名 "isPrivate"
func (m *ReflectionMethodIsPrivateMethod) GetName() string { return "isPrivate" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodIsPrivateMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodIsPrivateMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodIsPrivateMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodIsPrivateMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionMethodIsPrivateMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 isPrivate 方法
// 检查被反射的方法是否为私有方法
func (m *ReflectionMethodIsPrivateMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(method.GetModifier() == data.ModifierPrivate), nil
}
