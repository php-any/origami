package protowire

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EncodeFixed32Method implements Protowire::encodeFixed32()
type EncodeFixed32Method struct{}

func NewEncodeFixed32Method() data.Method {
	return &EncodeFixed32Method{}
}

func (m *EncodeFixed32Method) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	n, err := toUint64(val)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendFixed32(nil, uint32(n))
	return data.NewStringValue(string(buf)), nil
}

func (m *EncodeFixed32Method) GetName() string {
	return "encodeFixed32"
}

func (m *EncodeFixed32Method) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *EncodeFixed32Method) GetIsStatic() bool {
	return true
}

func (m *EncodeFixed32Method) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (m *EncodeFixed32Method) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

func (m *EncodeFixed32Method) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
