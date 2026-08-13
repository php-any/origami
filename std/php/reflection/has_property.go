package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassHasPropertyMethod 实现 ReflectionClass::hasProperty
// 检查被反射的类是否包含指定的属性
type ReflectionClassHasPropertyMethod struct{}

// GetName 返回方法名 "hasProperty"
func (m *ReflectionClassHasPropertyMethod) GetName() string { return "hasProperty" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassHasPropertyMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassHasPropertyMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - name: 属性名（字符串）
func (m *ReflectionClassHasPropertyMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassHasPropertyMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionClassHasPropertyMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 hasProperty 方法
// 检查被反射的类是否包含指定的属性
// 返回 true 表示属性存在，false 表示不存在
func (m *ReflectionClassHasPropertyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	propertyNameValue, _ := ctx.GetIndexValue(0)
	if propertyNameValue == nil {
		return data.NewBoolValue(false), nil
	}

	propertyName := propertyNameValue.AsString()
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	_, exists := classStmt.GetProperty(propertyName)
	return data.NewBoolValue(exists), nil
}
