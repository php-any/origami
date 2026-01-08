package reflection

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterGetDefaultValueMethod 实现 ReflectionParameter::getDefaultValue
// 返回被反射参数的默认值
type ReflectionParameterGetDefaultValueMethod struct{}

// GetName 返回方法名 "getDefaultValue"
func (m *ReflectionParameterGetDefaultValueMethod) GetName() string { return "getDefaultValue" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterGetDefaultValueMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterGetDefaultValueMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionParameterGetDefaultValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionParameterGetDefaultValueMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回混合类型
func (m *ReflectionParameterGetDefaultValueMethod) GetReturnType() data.Types {
	return nil
}

// Call 执行 getDefaultValue 方法
// 返回被反射参数的默认值
// 如果参数没有默认值，抛出异常
func (m *ReflectionParameterGetDefaultValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, param := getReflectionParameterInfo(ctx)
	if param == nil {
		return nil, data.NewErrorThrow(nil, errors.New("Parameter does not have a default value"))
	}

	var defaultValue data.GetValue

	// 获取参数的默认值
	if paramInterface, ok := param.(data.Parameter); ok {
		defaultValue = paramInterface.GetDefaultValue()
	} else if paramNode, ok := param.(*node.Parameter); ok {
		defaultValue = paramNode.DefaultValue
	}

	if defaultValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("Parameter does not have a default value"))
	}

	// 获取默认值的实际值
	value, acl := defaultValue.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	return value, nil
}
