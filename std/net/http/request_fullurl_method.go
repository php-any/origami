package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// RequestFullUrlMethod 获取完整 URL
type RequestFullUrlMethod struct {
	source *httpsrc.Request
}

func (h *RequestFullUrlMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewStringValue(""), nil
	}

	// 构建完整 URL
	scheme := "http"
	if h.source.TLS != nil {
		scheme = "https"
	}

	host := h.source.Host
	if host == "" {
		host = "localhost"
	}

	fullUrl := scheme + "://" + host + h.source.RequestURI
	return data.NewStringValue(fullUrl), nil
}

func (h *RequestFullUrlMethod) GetName() string               { return "fullUrl" }
func (h *RequestFullUrlMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestFullUrlMethod) GetIsStatic() bool             { return false }
func (h *RequestFullUrlMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestFullUrlMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestFullUrlMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
