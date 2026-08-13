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

	flags := 0
	if flagsValue != nil {
		if asInt, ok := flagsValue.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				flags = v
			}
		}
	}
	instanceof := flags&2 != 0 // ReflectionAttribute::IS_INSTANCEOF

	// 获取类的注解/属性
	attributes := []data.Value{}

	if classStatement, ok := classStmt.(*node.ClassStatement); ok {
		if classStatement.Annotations != nil && len(classStatement.Annotations) > 0 {
			var filterName string
			if nameValue != nil {
				if _, isNull := nameValue.(*data.NullValue); !isNull {
					if strVal, ok := nameValue.(*data.StringValue); ok {
						filterName = strVal.AsString()
					} else if nameValue.AsString() != "" {
						filterName = nameValue.AsString()
					}
				}
			}

			for _, annotation := range classStatement.Annotations {
				if filterName != "" {
					annName := annotation.Class.GetName()
					if instanceof {
						if !attributeMatchesName(ctx, annotation.Class, filterName) {
							continue
						}
					} else if annName != filterName {
						continue
					}
				}
				attrValue := newReflectionAttribute(ctx, annotation)
				attributes = append(attributes, attrValue)
			}
		}
	}

	return data.NewArrayValue(attributes), nil
}

// attributeMatchesName 判断注解类名是否等于 name，或（IS_INSTANCEOF）为 name 的子类。
func attributeMatchesName(ctx data.Context, annClass data.ClassStmt, name string) bool {
	if annClass == nil {
		return false
	}
	current := annClass
	for current != nil {
		if current.GetName() == name {
			return true
		}
		extend := current.GetExtend()
		if extend == nil || *extend == "" {
			break
		}
		parent, acl := ctx.GetVM().GetOrLoadClass(*extend)
		if acl != nil || parent == nil {
			break
		}
		current = parent
	}
	return false
}
