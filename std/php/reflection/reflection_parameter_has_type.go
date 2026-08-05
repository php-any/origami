package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionParameterHasTypeMethod 实现 ReflectionParameter::hasType
type ReflectionParameterHasTypeMethod struct{}

func (m *ReflectionParameterHasTypeMethod) GetName() string { return "hasType" }

func (m *ReflectionParameterHasTypeMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ReflectionParameterHasTypeMethod) GetIsStatic() bool { return false }

func (m *ReflectionParameterHasTypeMethod) GetParams() []data.GetValue { return []data.GetValue{} }

func (m *ReflectionParameterHasTypeMethod) GetVariables() []data.Variable { return []data.Variable{} }

func (m *ReflectionParameterHasTypeMethod) GetReturnType() data.Types { return data.Bool{} }

func (m *ReflectionParameterHasTypeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, _, param := getReflectionParameterInfo(ctx)
	if param == nil {
		return data.NewBoolValue(false), nil
	}

	if vp, ok := param.(*virtualParam); ok {
		return data.NewBoolValue(vp.GetType() != nil), nil
	}

	typeStr := extractParamType(param)
	return data.NewBoolValue(typeStr != ""), nil
}
