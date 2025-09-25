package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type CookieStringMethod struct {
	source *httpsrc.Cookie
}

func (h *CookieStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.String()
	return data.NewAnyValue(ret0), nil
}

func (h *CookieStringMethod) GetName() string               { return "string" }
func (h *CookieStringMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *CookieStringMethod) GetIsStatic() bool             { return true }
func (h *CookieStringMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *CookieStringMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *CookieStringMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
