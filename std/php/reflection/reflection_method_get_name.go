package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodGetNameMethod 实现 ReflectionMethod::getName
// 返回被反射的方法的名称
type ReflectionMethodGetNameMethod struct{}

// GetName 返回方法名 "getName"
func (m *ReflectionMethodGetNameMethod) GetName() string { return "getName" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodGetNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodGetNameMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodGetNameMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodGetNameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回字符串类型
func (m *ReflectionMethodGetNameMethod) GetReturnType() data.Types {
	return data.String{}
}

// Call 执行 getName 方法
// 返回被反射的方法的名称
func (m *ReflectionMethodGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, methodName, _ := getReflectionMethodInfo(ctx)
	return data.NewStringValue(methodName), nil
}
