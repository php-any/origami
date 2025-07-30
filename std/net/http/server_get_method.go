package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ServerGetMethod struct {
	source *Server
}

func (h *ServerGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a1, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a2, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	h.source.Get(ctx, *(a1.(*data.StringValue)), *(a2.(*data.FuncValue)))
	return nil, nil
}

func (h *ServerGetMethod) GetName() string {
	return "get"
}

func (h *ServerGetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ServerGetMethod) GetIsStatic() bool {
	return false
}

func (h *ServerGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "handler", 1, nil, data.NewBaseType("callable")),
	}
}

func (h *ServerGetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "handler", 1, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *ServerGetMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
