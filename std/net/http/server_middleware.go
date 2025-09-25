package http

import (
	"github.com/php-any/origami/data"
)

type ServerMiddlewareMethod struct {
	server *ServerClass
}

func (h *ServerMiddlewareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

func (h *ServerMiddlewareMethod) GetName() string            { return "middleware" }
func (h *ServerMiddlewareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerMiddlewareMethod) GetIsStatic() bool          { return false }
func (h *ServerMiddlewareMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (h *ServerMiddlewareMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *ServerMiddlewareMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
