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
func (h *CookieValidMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *CookieValidMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *CookieValidMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
