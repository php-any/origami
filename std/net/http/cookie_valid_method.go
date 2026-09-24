package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

type CookieValidMethod struct {
	source *httpsrc.Cookie
}

func (h *CookieValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.Valid()
	return data.NewAnyValue(ret0), nil
}

func (h *CookieValidMethod) GetName() string               { return "valid" }
func (h *CookieValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *CookieValidMethod) GetIsStatic() bool             { return false }
var cookieValidMethodGetParams = []data.GetValue{}

func (h *CookieValidMethod) GetParams() []data.GetValue    { return cookieValidMethodGetParams }
var cookieValidMethodGetVariables = []data.Variable{}

func (h *CookieValidMethod) GetVariables() []data.Variable { return cookieValidMethodGetVariables }
func (h *CookieValidMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
