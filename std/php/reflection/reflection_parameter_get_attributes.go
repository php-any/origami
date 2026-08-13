package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterGetAttributesMethod 实现 ReflectionParameter::getAttributes。
type ReflectionParameterGetAttributesMethod struct{}

func (m *ReflectionParameterGetAttributesMethod) GetName() string { return "getAttributes" }

func (m *ReflectionParameterGetAttributesMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ReflectionParameterGetAttributesMethod) GetIsStatic() bool { return false }

func (m *ReflectionParameterGetAttributesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Mixed{}),
	}
}

func (m *ReflectionParameterGetAttributesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Mixed{}),
	}
}

func (m *ReflectionParameterGetAttributesMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

func (m *ReflectionParameterGetAttributesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, reflected := getReflectionParameterInfo(ctx)
	param, ok := reflected.(*node.Parameter)
	if !ok || len(param.Annotations) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	filterName := ""
	if name, hasName := ctx.GetIndexValue(0); hasName && name != nil {
		if _, isNull := name.(*data.NullValue); !isNull {
			filterName = name.AsString()
		}
	}

	attributes := make([]data.Value, 0, len(param.Annotations))
	for _, annotation := range param.Annotations {
		if annotation == nil || annotation.Class == nil {
			continue
		}
		if filterName != "" && annotation.Class.GetName() != filterName {
			continue
		}
		attributes = append(attributes, newReflectionAttribute(ctx, annotation))
	}
	return data.NewArrayValue(attributes), nil
}
