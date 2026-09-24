package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// RequestMethodMethod 获取请求方法
type RequestMethodMethod struct {
	source *httpsrc.Request
}

func (h *RequestMethodMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(h.source.Method), nil
}

func (h *RequestMethodMethod) GetName() string               { return "method" }
func (h *RequestMethodMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestMethodMethod) GetIsStatic() bool             { return false }
var requestMethodMethodGetParams = []data.GetValue{}

func (h *RequestMethodMethod) GetParams() []data.GetValue    { return requestMethodMethodGetParams }
var requestMethodMethodGetVariables = []data.Variable{}

func (h *RequestMethodMethod) GetVariables() []data.Variable { return requestMethodMethodGetVariables }
func (h *RequestMethodMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
