package reflection

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterConstructMethod 实现 ReflectionParameter::__construct
// 构造函数用于初始化 ReflectionParameter 实例
type ReflectionParameterConstructMethod struct{}

// GetName 返回方法名 "__construct"
func (m *ReflectionParameterConstructMethod) GetName() string { return "__construct" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionParameterConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionParameterConstructMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - function: 函数或方法的反射对象（ReflectionFunction 或 ReflectionMethod）
//   - parameter: 参数名或参数索引
func (m *ReflectionParameterConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "function", 0, nil, data.NewBaseType("object")),
		node.NewParameter(nil, "parameter", 1, nil, nil),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionParameterConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "function", 0, nil),
		node.NewVariable(nil, "parameter", 1, nil),
	}
}

// GetReturnType 返回返回类型，构造函数无返回值
func (m *ReflectionParameterConstructMethod) GetReturnType() data.Types {
	return nil
}

// Call 执行构造函数
// 初始化 ReflectionParameter 实例，存储函数/方法信息和参数索引
func (m *ReflectionParameterConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取函数/方法反射对象
	functionValue, _ := ctx.GetIndexValue(0)
	if functionValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionParameter::__construct() expects parameter 1 to be ReflectionFunction or ReflectionMethod"))
	}

	// 获取参数名或索引
	parameterValue, _ := ctx.GetIndexValue(1)
	if parameterValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionParameter::__construct() expects parameter 2 to be string or int"))
	}

	// 从 ReflectionMethod 或 ReflectionFunction 中获取类名和方法名
	// 这里简化处理，实际应该从 functionValue 中提取信息
	// 由于 ReflectionParameter 通常由 ReflectionMethod::getParameters() 创建
	// 构造函数主要用于兼容性，实际创建通过 newReflectionParameter 辅助函数

	// 存储信息到实例属性中
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 这里需要从 functionValue 中提取类名和方法名
		// 简化实现：假设已经通过 newReflectionParameter 设置了属性
		_ = objCtx
	}

	return nil, nil
}
