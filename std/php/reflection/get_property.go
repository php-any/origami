package reflection

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetPropertyMethod 实现 ReflectionClass::getProperty
// 根据属性名获取被反射类的指定属性
type ReflectionClassGetPropertyMethod struct{}

// GetName 返回方法名 "getProperty"
func (m *ReflectionClassGetPropertyMethod) GetName() string { return "getProperty" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetPropertyMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetPropertyMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - name: 属性名（字符串）
func (m *ReflectionClassGetPropertyMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassGetPropertyMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}

// GetReturnType 返回返回类型，返回混合类型（当前实现返回字符串）
func (m *ReflectionClassGetPropertyMethod) GetReturnType() data.Types {
	return data.Mixed{}
}

// Call 执行 getProperty 方法
// 根据属性名查找并返回对应的属性
// 如果属性不存在，抛出异常
// TODO: 当前实现返回属性名字符串，实际应该返回 ReflectionProperty 对象
func (m *ReflectionClassGetPropertyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	propertyNameValue, _ := ctx.GetIndexValue(0)
	if propertyNameValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionClass::getProperty() expects parameter 1 to be string"))
	}

	propertyName := propertyNameValue.AsString()
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Property %s does not exist", propertyName))
	}

	// 查找属性
	prop, exists := classStmt.GetProperty(propertyName)
	if !exists {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Property %s does not exist", propertyName))
	}

	// 返回属性名（简化实现，实际应该返回 ReflectionProperty 对象）
	return data.NewStringValue(prop.GetName()), nil
}
