package http

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type RequestInputMethod struct {
	source *Request
}

func (h *RequestInputMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	h.source.Input(a0.(*data.StringValue).AsString())
	return nil, nil
}

func (h *RequestInputMethod) GetName() string {
	return "input"
}

func (h *RequestInputMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *RequestInputMethod) GetIsStatic() bool {
	return false
}

func (h *RequestInputMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}

func (h *RequestInputMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
