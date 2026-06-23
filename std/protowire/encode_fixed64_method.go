package protowire

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EncodeFixed64Method implements Protowire::encodeFixed64()
type EncodeFixed64Method struct{}

func NewEncodeFixed64Method() data.Method {
	return &EncodeFixed64Method{}
}

func (m *EncodeFixed64Method) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	n, err := toUint64(val)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendFixed64(nil, n)
	return data.NewStringValue(string(buf)), nil
}

func (m *EncodeFixed64Method) GetName() string {
	return "encodeFixed64"
}

func (m *EncodeFixed64Method) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *EncodeFixed64Method) GetIsStatic() bool {
	return true
}

func (m *EncodeFixed64Method) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (m *EncodeFixed64Method) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

func (m *EncodeFixed64Method) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
