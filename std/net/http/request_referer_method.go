package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestRefererMethod struct {
	source *httpsrc.Request
}

func (h *RequestRefererMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.Referer()
	return data.NewAnyValue(ret0), nil
}

func (h *RequestRefererMethod) GetName() string               { return "referer" }
func (h *RequestRefererMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestRefererMethod) GetIsStatic() bool             { return false }
func (h *RequestRefererMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestRefererMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestRefererMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
