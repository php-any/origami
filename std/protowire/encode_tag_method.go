package protowire

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EncodeTagMethod implements Protowire::encodeTag()
type EncodeTagMethod struct{}

func NewEncodeTagMethod() data.Method {
	return &EncodeTagMethod{}
}

func (m *EncodeTagMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	numVal, _ := ctx.GetIndexValue(0)
	wtypeVal, _ := ctx.GetIndexValue(1)
	if numVal == nil || wtypeVal == nil {
		return data.NewStringValue(""), nil
	}
	num, err := toUint64(numVal)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	wtype, err := toInt32(wtypeVal)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	buf := pwAppendTag(nil, pwNumber(int(num)), pwType(int(wtype)))
	return data.NewStringValue(string(buf)), nil
}

func (m *EncodeTagMethod) GetName() string {
	return "encodeTag"
}

func (m *EncodeTagMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *EncodeTagMethod) GetIsStatic() bool {
	return true
}

func (m *EncodeTagMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "number", 0, nil, nil),
		node.NewParameter(nil, "wireType", 1, nil, nil),
	}
}

func (m *EncodeTagMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "number", 0, nil),
		node.NewVariable(nil, "wireType", 1, nil),
	}
}

func (m *EncodeTagMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
