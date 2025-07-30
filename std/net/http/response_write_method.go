package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ResponseWriteMethod struct {
	source *Response
}

func (h *ResponseWriteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	h.source.Write(a0.(*data.StringValue).AsString())
	return nil, nil
}

func (h *ResponseWriteMethod) GetName() string {
	return "write"
}

func (h *ResponseWriteMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ResponseWriteMethod) GetIsStatic() bool {
	return false
}

func (h *ResponseWriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
	}
}

func (h *ResponseWriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *ResponseWriteMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
