package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ReflectionMethodGetAttributesMethod struct{}

func (m *ReflectionMethodGetAttributesMethod) GetName() string { return "getAttributes" }
func (m *ReflectionMethodGetAttributesMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionMethodGetAttributesMethod) GetIsStatic() bool { return false }
func (m *ReflectionMethodGetAttributesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Mixed{}),
	}
}
func (m *ReflectionMethodGetAttributesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Mixed{}),
	}
}
func (m *ReflectionMethodGetAttributesMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

func (m *ReflectionMethodGetAttributesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, reflected := getReflectionMethodInfo(ctx)
	method, ok := reflected.(*node.ClassMethod)
	if !ok || len(method.Annotations) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	filterName := ""
	if name, hasName := ctx.GetIndexValue(0); hasName && name != nil {
		if _, isNull := name.(*data.NullValue); !isNull {
			filterName = name.AsString()
		}
	}

	attributes := make([]data.Value, 0, len(method.Annotations))
	for _, annotation := range method.Annotations {
		if annotation == nil {
			continue
		}
		if filterName != "" && annotation.GetName() != filterName {
			continue
		}
		attributes = append(attributes, newReflectionAttribute(ctx, annotation))
	}
	return data.NewArrayValue(attributes), nil
}
