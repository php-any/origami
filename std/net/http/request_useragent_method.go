package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestUserAgentMethod struct {
	source *httpsrc.Request
}

func (h *RequestUserAgentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.UserAgent()
	return data.NewAnyValue(ret0), nil
}

func (h *RequestUserAgentMethod) GetName() string               { return "userAgent" }
func (h *RequestUserAgentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestUserAgentMethod) GetIsStatic() bool             { return false }
func (h *RequestUserAgentMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestUserAgentMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestUserAgentMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
