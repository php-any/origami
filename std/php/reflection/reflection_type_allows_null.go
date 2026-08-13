package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionTypeAllowsNullMethod 实现 ReflectionType::allowsNull
// 检查被反射类型是否允许 null 值
type ReflectionTypeAllowsNullMethod struct{}

// GetName 返回方法名 "allowsNull"
func (m *ReflectionTypeAllowsNullMethod) GetName() string { return "allowsNull" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionTypeAllowsNullMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionTypeAllowsNullMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionTypeAllowsNullMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionTypeAllowsNullMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionTypeAllowsNullMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

// Call 执行 allowsNull 方法
// 检查被反射类型是否允许 null 值
// 如果类型是可空类型（以 ? 开头）或者是 mixed，则允许 null
func (m *ReflectionTypeAllowsNullMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	typeName, typeInfo := getReflectionTypeInfo(ctx)
	if typeName == "" {
		return data.NewBoolValue(true), nil // 没有类型声明时，默认允许 null
	}

	// 检查是否为可空类型（以 ? 开头）
	if len(typeName) > 0 && typeName[0] == '?' {
		return data.NewBoolValue(true), nil
	}

	// 检查是否为 mixed 类型
	if typeName == "mixed" {
		return data.NewBoolValue(true), nil
	}

	// 检查是否为 NullableType
	if _, ok := typeInfo.(data.NullableType); ok {
		return data.NewBoolValue(true), nil
	}

	return data.NewBoolValue(false), nil
}
