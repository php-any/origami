package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionTypeConstructMethod 实现 ReflectionType::__construct
// 构造函数用于初始化 ReflectionType 实例
type ReflectionTypeConstructMethod struct{}

// GetName 返回方法名 "__construct"
func (m *ReflectionTypeConstructMethod) GetName() string { return "__construct" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionTypeConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionTypeConstructMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// ReflectionType 通常由 ReflectionParameter::getType() 创建，构造函数主要用于兼容性
func (m *ReflectionTypeConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表
func (m *ReflectionTypeConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，构造函数无返回值
func (m *ReflectionTypeConstructMethod) GetReturnType() data.Types {
	return nil
}

// Call 执行构造函数
// ReflectionType 通常由 ReflectionParameter::getType() 创建，构造函数主要用于兼容性
func (m *ReflectionTypeConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// ReflectionType 通常由 newReflectionType 辅助函数创建
	// 构造函数主要用于兼容性，实际创建通过 newReflectionType 辅助函数
	return nil, nil
}
