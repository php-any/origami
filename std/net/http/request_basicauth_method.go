package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestBasicAuthMethod struct {
	source *httpsrc.Request
}

func (h *RequestBasicAuthMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0, ret1, ret2 := h.source.BasicAuth()
	return data.NewArrayValue([]data.Value{data.NewAnyValue(ret0), data.NewAnyValue(ret1), data.NewAnyValue(ret2)}), nil
}

func (h *RequestBasicAuthMethod) GetName() string               { return "basicAuth" }
func (h *RequestBasicAuthMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestBasicAuthMethod) GetIsStatic() bool             { return false }
func (h *RequestBasicAuthMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestBasicAuthMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestBasicAuthMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
