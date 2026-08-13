package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterIsOptionalMethod 实现 ReflectionParameter::isOptional
// 检查被反射参数是否为可选参数（有默认值）
type ReflectionParameterIsOptionalMethod struct{}

// GetName 返回方法名 "isOptional"
func (m *ReflectionParameterIsOptionalMethod) GetName() string { return "isOptional" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterIsOptionalMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterIsOptionalMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionParameterIsOptionalMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionParameterIsOptionalMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionParameterIsOptionalMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

// Call 执行 isOptional 方法
// 检查被反射参数是否为可选参数（有默认值）
func (m *ReflectionParameterIsOptionalMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, param := getReflectionParameterInfo(ctx)
	if param == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查参数是否有默认值
	if paramInterface, ok := param.(data.Parameter); ok {
		defaultValue := paramInterface.GetDefaultValue()
		return data.NewBoolValue(defaultValue != nil), nil
	} else if paramNode, ok := param.(*node.Parameter); ok {
		return data.NewBoolValue(paramNode.DefaultValue != nil), nil
	}

	return data.NewBoolValue(false), nil
}
