package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestContextMethod struct {
	source *httpsrc.Request
}

func (h *RequestContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.Context()
	return data.NewAnyValue(ret0), nil
}

func (h *RequestContextMethod) GetName() string               { return "context" }
func (h *RequestContextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestContextMethod) GetIsStatic() bool             { return true }
func (h *RequestContextMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestContextMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestContextMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
