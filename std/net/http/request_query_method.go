package http

import (
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
)

// RequestQueryMethod 获取查询参数
type RequestQueryMethod struct {
	source *httpsrc.Request
}

func (h *RequestQueryMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil || h.source.URL == nil {
		return data.NewObjectValue(), nil
	}

	// 将查询参数转换为对象
	result := data.NewObjectValue()
	for key, values := range h.source.URL.Query() {
		if len(values) == 1 {
			result.SetProperty(key, data.NewStringValue(values[0]))
		} else {
			// 如果有多个值，用逗号分隔
			result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
		}
	}

	return result, nil
}

func (h *RequestQueryMethod) GetName() string               { return "query" }
func (h *RequestQueryMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestQueryMethod) GetIsStatic() bool             { return false }
func (h *RequestQueryMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestQueryMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestQueryMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
