package http

import (
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
)

// RequestIpMethod 获取客户端 IP
type RequestIpMethod struct {
	source *httpsrc.Request
}

func (h *RequestIpMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewStringValue(""), nil
	}

	// 优先从 X-Forwarded-For 获取
	if xff := h.source.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return data.NewStringValue(strings.TrimSpace(ips[0])), nil
		}
	}

	// 从 X-Real-IP 获取
	if xri := h.source.Header.Get("X-Real-IP"); xri != "" {
		return data.NewStringValue(xri), nil
	}

	// 最后使用 RemoteAddr
	return data.NewStringValue(h.source.RemoteAddr), nil
}

func (h *RequestIpMethod) GetName() string               { return "ip" }
func (h *RequestIpMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestIpMethod) GetIsStatic() bool             { return false }
func (h *RequestIpMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestIpMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestIpMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
