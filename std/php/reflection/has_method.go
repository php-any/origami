package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassHasMethodMethod 实现 ReflectionClass::hasMethod
// 检查被反射的类是否包含指定的方法
type ReflectionClassHasMethodMethod struct{}

// GetName 返回方法名 "hasMethod"
func (m *ReflectionClassHasMethodMethod) GetName() string { return "hasMethod" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassHasMethodMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassHasMethodMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - name: 方法名（字符串）
func (m *ReflectionClassHasMethodMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassHasMethodMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionClassHasMethodMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 hasMethod 方法
// 检查被反射的类是否包含指定的方法
// 返回 true 表示方法存在，false 表示不存在
func (m *ReflectionClassHasMethodMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	methodNameValue, _ := ctx.GetIndexValue(0)
	if methodNameValue == nil {
		return data.NewBoolValue(false), nil
	}

	methodName := methodNameValue.AsString()
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	_, exists := classStmt.GetMethod(methodName)
	return data.NewBoolValue(exists), nil
}
