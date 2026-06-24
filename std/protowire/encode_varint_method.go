package protowire

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EncodeVarintMethod implements Protowire::encodeVarint()
type EncodeVarintMethod struct{}

func NewEncodeVarintMethod() data.Method {
	return &EncodeVarintMethod{}
}

func (m *EncodeVarintMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewStringValue(""), nil
	}
	n, err := toUint64(val)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendVarint(nil, n)
	return data.NewStringValue(string(buf)), nil
}

func (m *EncodeVarintMethod) GetName() string {
	return "encodeVarint"
}

func (m *EncodeVarintMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *EncodeVarintMethod) GetIsStatic() bool {
	return true
}

func (m *EncodeVarintMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (m *EncodeVarintMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

func (m *EncodeVarintMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
