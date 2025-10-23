package http

import (
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
)

// RequestAllMethod 获取所有输入数据
type RequestAllMethod struct {
	source *httpsrc.Request
}

func (h *RequestAllMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewObjectValue(), nil
	}

	// 合并所有输入数据
	result := data.NewObjectValue()

	// 从查询参数获取
	for key, values := range h.source.URL.Query() {
		if len(values) == 1 {
			result.SetProperty(key, data.NewStringValue(values[0]))
		} else {
			result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
		}
	}

	// 从表单数据获取
	if h.source.Form != nil {
		for key, values := range h.source.Form {
			if len(values) == 1 {
				result.SetProperty(key, data.NewStringValue(values[0]))
			} else {
				result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
			}
		}
	}

	return result, nil
}

func (h *RequestAllMethod) GetName() string               { return "all" }
func (h *RequestAllMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestAllMethod) GetIsStatic() bool             { return false }
func (h *RequestAllMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestAllMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestAllMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
