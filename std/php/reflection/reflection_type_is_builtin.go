package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionTypeIsBuiltinMethod 实现 ReflectionType::isBuiltin
// 检查被反射类型是否为内置类型
type ReflectionTypeIsBuiltinMethod struct{}

// GetName 返回方法名 "isBuiltin"
func (m *ReflectionTypeIsBuiltinMethod) GetName() string { return "isBuiltin" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionTypeIsBuiltinMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionTypeIsBuiltinMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionTypeIsBuiltinMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionTypeIsBuiltinMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionTypeIsBuiltinMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

// Call 执行 isBuiltin 方法
// 检查被反射类型是否为内置类型
// PHP 内置类型包括：int, float, string, bool, array, object, callable, iterable, void, mixed 等
func (m *ReflectionTypeIsBuiltinMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	typeName, _ := getReflectionTypeInfo(ctx)
	if typeName == "" {
		return data.NewBoolValue(false), nil
	}

	// 检查是否为内置类型
	isBuiltin := data.ISBaseType(typeName)

	// 如果是可空类型（以 ? 开头），检查基础类型是否为内置类型
	if len(typeName) > 0 && typeName[0] == '?' {
		baseTypeName := typeName[1:]
		isBuiltin = data.ISBaseType(baseTypeName)
	}

	return data.NewBoolValue(isBuiltin), nil
}
