package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetAttributesMethod 实现 ReflectionClass::getAttributes
// 返回类的所有属性（attributes/annotations），返回 ReflectionAttribute 对象数组
type ReflectionClassGetAttributesMethod struct{}

func (m *ReflectionClassGetAttributesMethod) GetName() string { return "getAttributes" }

func (m *ReflectionClassGetAttributesMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ReflectionClassGetAttributesMethod) GetIsStatic() bool { return false }

func (m *ReflectionClassGetAttributesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Mixed{}),
	}
}

func (m *ReflectionClassGetAttributesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Mixed{}),
	}
}

func (m *ReflectionClassGetAttributesMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

// newReflectionAttribute 创建一个新的 ReflectionAttribute 实例
func newReflectionAttribute(ctx data.Context, annotation *data.ClassValue) *data.ClassValue {
	attrClass := &ReflectionAttributeClass{}
	attrValue := data.NewClassValue(attrClass, ctx.CreateBaseContext())

	// 存储注解对象到实例属性中
	attrValue.ObjectValue.SetProperty("_annotation", annotation)

	return attrValue
}

func (m *ReflectionClassGetAttributesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取被反射的类信息
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 获取参数
	nameValue, _ := ctx.GetIndexValue(0)  // name 参数，可选，默认为 null
	flagsValue, _ := ctx.GetIndexValue(1) // flags 参数，可选，默认为 0

	// 获取类的注解/属性
	// 检查是否是 ClassStatement 类型，并获取其注解
	attributes := []data.Value{}

	// 尝试将 classStmt 转换为 ClassStatement 以访问注解
	if classStatement, ok := classStmt.(*node.ClassStatement); ok {
		// 检查注解列表是否存在
		if classStatement.Annotations != nil && len(classStatement.Annotations) > 0 {
			// 如果指定了 name 参数，过滤特定名称的注解
			var filterName string
			if nameValue != nil {
				if strVal, ok := nameValue.(*data.StringValue); ok {
					filterName = strVal.AsString()
				} else if nameValue.AsString() != "" {
					filterName = nameValue.AsString()
				}
			}

			// 将注解转换为 ReflectionAttribute 对象
			for _, annotation := range classStatement.Annotations {
				// 如果指定了 name，只返回匹配的注解
				if filterName != "" {
					if annotation.Class.GetName() != filterName {
						continue
					}
				}
				// 创建 ReflectionAttribute 实例
				attrValue := newReflectionAttribute(ctx, annotation)
				attributes = append(attributes, attrValue)
			}
		}
	}

	// TODO: 实现 flags 参数的过滤逻辑
	// flags 可以用于过滤继承的注解等
	_ = flagsValue

	return data.NewArrayValue(attributes), nil
}
