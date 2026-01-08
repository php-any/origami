package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterGetNameMethod 实现 ReflectionParameter::getName
// 返回被反射参数的名称
type ReflectionParameterGetNameMethod struct{}

// GetName 返回方法名 "getName"
func (m *ReflectionParameterGetNameMethod) GetName() string { return "getName" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterGetNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterGetNameMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionParameterGetNameMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionParameterGetNameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回字符串类型
func (m *ReflectionParameterGetNameMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

// Call 执行 getName 方法
// 返回被反射参数的名称
func (m *ReflectionParameterGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, param := getReflectionParameterInfo(ctx)
	if param == nil {
		return data.NewStringValue(""), nil
	}

	var paramName string
	// 尝试多种类型断言来获取参数名
	if paramVar, ok := param.(data.Variable); ok {
		paramName = paramVar.GetName()
	} else if paramInterface, ok := param.(data.Parameter); ok {
		paramName = paramInterface.GetName()
	} else if paramNode, ok := param.(*node.Parameter); ok {
		paramName = paramNode.GetName()
	} else if paramNodes, ok := param.(*node.Parameters); ok {
		paramName = paramNodes.GetName()
	}

	return data.NewStringValue(paramName), nil
}
