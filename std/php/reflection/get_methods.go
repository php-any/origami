package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetMethodsMethod 实现 ReflectionClass::getMethods
// 返回被反射的类的所有方法列表
type ReflectionClassGetMethodsMethod struct{}

// GetName 返回方法名 "getMethods"
func (m *ReflectionClassGetMethodsMethod) GetName() string { return "getMethods" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetMethodsMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetMethodsMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - filter: 可选的过滤器标志（?int），用于过滤方法的可见性，默认为 null
func (m *ReflectionClassGetMethodsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filter", 0, node.NewNullLiteral(nil), data.NewBaseType("?int")),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassGetMethodsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filter", 0, data.NewBaseType("?int")),
	}
}

// GetReturnType 返回返回类型，返回数组类型
func (m *ReflectionClassGetMethodsMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

// Call 执行 getMethods 方法
// 返回被反射的类的所有方法列表（当前实现返回方法名数组）
// TODO: 实现完整的过滤器逻辑，当前忽略 filter 参数
func (m *ReflectionClassGetMethodsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取类信息
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 获取所有方法
	methods := classStmt.GetMethods()
	result := make([]data.Value, 0, len(methods))

	// 获取 filter 参数（可选）
	_, _ = ctx.GetIndexValue(0)
	// TODO: 实现完整的过滤器逻辑

	for _, method := range methods {
		// 应用过滤器（简化实现，暂时忽略过滤器）
		// TODO: 实现完整的过滤器逻辑
		methodName := method.GetName()
		// 创建 ReflectionMethod 对象（简化实现，返回方法名）
		result = append(result, data.NewStringValue(methodName))
	}

	return data.NewArrayValue(result), nil
}
