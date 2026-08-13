package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterIsVariadicMethod 实现 ReflectionParameter::isVariadic
// 检查被反射参数是否为可变参数（variadic parameter，使用 ... 语法）
type ReflectionParameterIsVariadicMethod struct{}

// GetName 返回方法名 "isVariadic"
func (m *ReflectionParameterIsVariadicMethod) GetName() string { return "isVariadic" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterIsVariadicMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterIsVariadicMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionParameterIsVariadicMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionParameterIsVariadicMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionParameterIsVariadicMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

// Call 执行 isVariadic 方法
// 检查被反射参数是否为可变参数（variadic parameter）
// 可变参数使用 ... 语法定义，例如 function foo(...$args)
func (m *ReflectionParameterIsVariadicMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, param := getReflectionParameterInfo(ctx)
	if param == nil {
		return data.NewBoolValue(false), nil
	}

	// 优先处理 virtualParam（Closure 参数）
	if vp, ok := param.(*virtualParam); ok {
		return data.NewBoolValue(vp.IsVariadic()), nil
	}

	// 检查参数是否是可变参数类型
	// node.Parameters 和 data.Parameters 都表示可变参数
	if _, ok := param.(*node.Parameters); ok {
		return data.NewBoolValue(true), nil
	}
	if _, ok := param.(data.Parameters); ok {
		return data.NewBoolValue(true), nil
	}

	return data.NewBoolValue(false), nil
}
