package http

import (
	"github.com/php-any/origami/data"
)

type ServerRunMethod struct {
	source *Server
}

func (h *ServerRunMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	h.source.Run(ctx)
	return nil, nil
}

func (h *ServerRunMethod) GetName() string {
	return "get"
}

func (h *ServerRunMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ServerRunMethod) GetIsStatic() bool {
	return false
}

func (h *ServerRunMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ServerRunMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (h *ServerRunMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
