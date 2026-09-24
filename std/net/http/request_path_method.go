package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// RequestPathMethod 获取请求路径
type RequestPathMethod struct {
	source *httpsrc.Request
}

func (h *RequestPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil || h.source.URL == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(h.source.URL.Path), nil
}

func (h *RequestPathMethod) GetName() string               { return "path" }
func (h *RequestPathMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestPathMethod) GetIsStatic() bool             { return false }
var requestPathMethodGetParams = []data.GetValue{}

func (h *RequestPathMethod) GetParams() []data.GetValue    { return requestPathMethodGetParams }
var requestPathMethodGetVariables = []data.Variable{}

func (h *RequestPathMethod) GetVariables() []data.Variable { return requestPathMethodGetVariables }
func (h *RequestPathMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
