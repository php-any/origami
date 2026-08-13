package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// RequestIsSecureMethod 检查是否是安全连接
type RequestIsSecureMethod struct {
	source *httpsrc.Request
}

func (h *RequestIsSecureMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(h.source.TLS != nil), nil
}

func (h *RequestIsSecureMethod) GetName() string               { return "isSecure" }
func (h *RequestIsSecureMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestIsSecureMethod) GetIsStatic() bool             { return false }
func (h *RequestIsSecureMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestIsSecureMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestIsSecureMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
