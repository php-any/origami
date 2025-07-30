package http

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ServerConstructMethod struct {
	source *Server
}

func (h *ServerConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, port: 0"))
	}
	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, addr: 1"))
	}

	port, err := a0.(*data.IntValue).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	addr := ""
	if v, ok := a1.(*data.StringValue); ok {
		addr = v.AsString()
	}

	h.source.Construct(ctx, port, addr)
	return nil, nil
}

func (h *ServerConstructMethod) GetName() string {
	return "construct"
}

func (h *ServerConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ServerConstructMethod) GetIsStatic() bool {
	return false
}

func (h *ServerConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "port", 0, data.NewIntValue(8080), data.NewBaseType("int")),
		node.NewParameter(nil, "addr", 1, data.NewStringValue("0.0.0.0"), data.NewBaseType("string")),
	}
}

func (h *ServerConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "port", 0, nil),
		node.NewVariable(nil, "addr", 1, nil),
	}
}
