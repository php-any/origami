package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestParseFormMethod struct {
	source *httpsrc.Request
}

func (h *RequestParseFormMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.ParseForm()
	return data.NewAnyValue(ret0), nil
}

func (h *RequestParseFormMethod) GetName() string               { return "parseForm" }
func (h *RequestParseFormMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestParseFormMethod) GetIsStatic() bool             { return true }
func (h *RequestParseFormMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestParseFormMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestParseFormMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
