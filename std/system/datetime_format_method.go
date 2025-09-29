package system

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DateTimeFormatMethod struct {
	source *DateTime
}

func (h *DateTimeFormatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	format, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0, name: format"))
	}

	return data.NewStringValue(h.source.Format(format.AsString())), nil
}

func (h *DateTimeFormatMethod) GetName() string {
	return "format"
}

func (h *DateTimeFormatMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *DateTimeFormatMethod) GetIsStatic() bool {
	return false
}

func (h *DateTimeFormatMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
	}
}

func (h *DateTimeFormatMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, nil),
	}
}

func (h *DateTimeFormatMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
