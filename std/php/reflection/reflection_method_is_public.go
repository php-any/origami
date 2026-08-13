package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodIsPublicMethod 实现 ReflectionMethod::isPublic
// 检查被反射的方法是否为公开方法
type ReflectionMethodIsPublicMethod struct{}

// GetName 返回方法名 "isPublic"
func (m *ReflectionMethodIsPublicMethod) GetName() string { return "isPublic" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodIsPublicMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodIsPublicMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodIsPublicMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodIsPublicMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionMethodIsPublicMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 isPublic 方法
// 检查被反射的方法是否为公开方法
func (m *ReflectionMethodIsPublicMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(method.GetModifier() == data.ModifierPublic), nil
}
