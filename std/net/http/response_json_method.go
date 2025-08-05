package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ResponseJsonMethod struct {
	source *Response
}

func (h *ResponseJsonMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	return h.source.Json(a0)
}

func (h *ResponseJsonMethod) GetName() string {
	return "json"
}

func (h *ResponseJsonMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ResponseJsonMethod) GetIsStatic() bool {
	return false
}

func (h *ResponseJsonMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
	}
}

func (h *ResponseJsonMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *ResponseJsonMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}
