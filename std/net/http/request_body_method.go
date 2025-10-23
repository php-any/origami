package http

import (
	"io"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// RequestBodyMethod 获取请求体数据
type RequestBodyMethod struct {
	source *httpsrc.Request
}

func (h *RequestBodyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewStringValue(""), nil
	}

	// 获取原始请求体
	if h.source.Body != nil {
		body, err := io.ReadAll(h.source.Body)
		if err == nil {
			return data.NewStringValue(string(body)), nil
		}
	}

	// 如果没有请求体，返回空字符串
	return data.NewStringValue(""), nil
}

func (h *RequestBodyMethod) GetName() string               { return "body" }
func (h *RequestBodyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestBodyMethod) GetIsStatic() bool             { return false }
func (h *RequestBodyMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestBodyMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestBodyMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
