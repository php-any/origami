package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodGetReturnTypeMethod 实现 ReflectionMethod::getReturnType
// 返回被反射方法的返回类型（ReflectionNamedType 对象），无返回类型时返回 null
type ReflectionMethodGetReturnTypeMethod struct{}

// GetName 返回方法名 "getReturnType"
func (m *ReflectionMethodGetReturnTypeMethod) GetName() string { return "getReturnType" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodGetReturnTypeMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodGetReturnTypeMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodGetReturnTypeMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodGetReturnTypeMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回 ReflectionNamedType 对象（或 null）
func (m *ReflectionMethodGetReturnTypeMethod) GetReturnType() data.Types {
	return data.NewBaseType("?ReflectionNamedType")
}

// Call 执行 getReturnType 方法
// 返回被反射方法的返回类型；如果方法没有声明返回类型，返回 null
func (m *ReflectionMethodGetReturnTypeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewNullValue(), nil
	}

	retType := method.GetReturnType()
	if retType == nil {
		return data.NewNullValue(), nil
	}

	return newPhpReflectionType(ctx, retType), nil
}
