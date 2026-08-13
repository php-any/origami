package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// RequestUrlMethod 获取请求 URL
type RequestUrlMethod struct {
	source *httpsrc.Request
}

func (h *RequestUrlMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil || h.source.URL == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(h.source.URL.String()), nil
}

func (h *RequestUrlMethod) GetName() string               { return "url" }
func (h *RequestUrlMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestUrlMethod) GetIsStatic() bool             { return false }
func (h *RequestUrlMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestUrlMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestUrlMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
