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

	// 合并所有输入数据：路由参数 → 查询参数 → 表单数据
	// 优先级：表单 > 查询 > 路由参数
	result := data.NewObjectValue()

	// 从路由参数获取（最低优先级）
	for key, val := range collectPathValues(h.source) {
		result.SetProperty(key, data.NewStringValue(val))
	}

	// 从查询参数获取
	for key, values := range h.source.URL.Query() {
		if len(values) == 1 {
			result.SetProperty(key, data.NewStringValue(values[0]))
		} else {
			result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
		}
	}

	// 从 POST 表单数据获取（最高优先级）
	// 使用 PostForm 而非 Form，因为 Form 已合并了 URL 查询参数
	if h.source.PostForm != nil {
		for key, values := range h.source.PostForm {
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
func (h *RequestAllMethod) GetReturnType() data.Types     { return data.NewBaseType("array") }
