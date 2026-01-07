package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionClassGetNameMethod 实现 ReflectionClass::getName
// 返回被反射的类的名称
type ReflectionClassGetNameMethod struct{}

// GetName 返回方法名 "getName"
func (m *ReflectionClassGetNameMethod) GetName() string { return "getName" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetNameMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionClassGetNameMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionClassGetNameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回字符串类型
func (m *ReflectionClassGetNameMethod) GetReturnType() data.Types {
	return data.String{}
}

// Call 执行 getName 方法
// 返回被反射的类的名称
func (m *ReflectionClassGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 从当前对象获取类名
	className, _ := getReflectionClassInfo(ctx)
	return data.NewStringValue(className), nil
}
