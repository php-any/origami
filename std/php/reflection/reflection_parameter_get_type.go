package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterGetTypeMethod 实现 ReflectionParameter::getType
// 返回被反射参数的类型
// 注意：PHP 中返回 ReflectionType 对象，这里简化实现返回类型名称字符串
type ReflectionParameterGetTypeMethod struct{}

// GetName 返回方法名 "getType"
func (m *ReflectionParameterGetTypeMethod) GetName() string { return "getType" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterGetTypeMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterGetTypeMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionParameterGetTypeMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionParameterGetTypeMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回 ReflectionNamedType 对象（或 null）
func (m *ReflectionParameterGetTypeMethod) GetReturnType() data.Types {
	return data.NewBaseType("?ReflectionNamedType")
}

// Call 执行 getType 方法
// 返回被反射参数的类型，返回 ReflectionNamedType 对象
// 如果参数没有类型声明，返回 null
func (m *ReflectionParameterGetTypeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, param := getReflectionParameterInfo(ctx)
	if param == nil {
		return data.NewNullValue(), nil
	}

	var paramType data.Types

	// 尝试多种类型断言来获取参数类型
	if paramVar, ok := param.(data.Variable); ok {
		paramType = paramVar.GetType()
	} else if paramInterface, ok := param.(data.Parameter); ok {
		paramType = paramInterface.GetType()
	} else if paramNode, ok := param.(*node.Parameter); ok {
		paramType = paramNode.GetType()
	} else if paramNodes, ok := param.(*node.Parameters); ok {
		paramType = paramNodes.GetType()
	}

	// 如果参数没有类型，返回 null
	if paramType == nil {
		return data.NewNullValue(), nil
	}

	// 创建并返回 ReflectionNamedType 对象
	// 在 PHP 中，getType() 返回的是 ReflectionNamedType（对于命名类型）
	return newReflectionNamedType(ctx, paramType), nil
}
