package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type RequestGetMethod struct {
	source *Request
}

func (h *RequestGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	h.source.Get(a0.(*data.StringValue).AsString())
	return nil, nil
}

func (h *RequestGetMethod) GetName() string {
	return "get"
}

func (h *RequestGetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *RequestGetMethod) GetIsStatic() bool {
	return false
}

func (h *RequestGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "handle", 1, nil, data.NewBaseType("callable")),
	}
}

func (h *RequestGetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
		node.NewVariable(nil, "handle", 1, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *RequestGetMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
