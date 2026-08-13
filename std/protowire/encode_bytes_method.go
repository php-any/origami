package protowire

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EncodeBytesMethod implements Protowire::encodeBytes()
type EncodeBytesMethod struct{}

func NewEncodeBytesMethod() data.Method {
	return &EncodeBytesMethod{}
}

func (m *EncodeBytesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	s, ok := val.(data.AsString)
	if !ok {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendBytes(nil, []byte(s.AsString()))
	return data.NewStringValue(string(buf)), nil
}

func (m *EncodeBytesMethod) GetName() string {
	return "encodeBytes"
}

func (m *EncodeBytesMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *EncodeBytesMethod) GetIsStatic() bool {
	return true
}

func (m *EncodeBytesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (m *EncodeBytesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

func (m *EncodeBytesMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
