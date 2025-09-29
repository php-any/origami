package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestCookiesMethod struct {
	source *httpsrc.Request
}

func (h *RequestCookiesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.Cookies()
	return data.NewAnyValue(ret0), nil
}

func (h *RequestCookiesMethod) GetName() string               { return "cookies" }
func (h *RequestCookiesMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestCookiesMethod) GetIsStatic() bool             { return false }
func (h *RequestCookiesMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestCookiesMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestCookiesMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
