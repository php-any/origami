package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

type ResponseWriterHeaderMethod struct {
	source httpsrc.ResponseWriter
}

func (h *ResponseWriterHeaderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.Header()
	return data.NewProxyValue(NewHeaderClassFrom(&ret0), ctx), nil
}

func (h *ResponseWriterHeaderMethod) GetName() string               { return "header" }
func (h *ResponseWriterHeaderMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *ResponseWriterHeaderMethod) GetIsStatic() bool             { return false }
func (h *ResponseWriterHeaderMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *ResponseWriterHeaderMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *ResponseWriterHeaderMethod) GetReturnType() data.Types {
	return data.NewBaseType("Net\\Http\\Header")
}
