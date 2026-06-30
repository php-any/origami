package http

import (
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestOnlyMethod 只获取指定的输入数据
type RequestOnlyMethod struct {
	source *httpsrc.Request
}

func (h *RequestOnlyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	// 获取参数列表
	keys := make([]string, 0)
	for i := 0; ; i++ {
		value, exists := ctx.GetIndexValue(i)
		if !exists {
			break
		}
		key, err := utils.Convert[string](value)
		if err != nil {
			return nil, utils.NewThrowf("参数转换失败: %v", err)
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return data.NewObjectValue(), nil
	}

	// 只返回指定的键，优先级：路由参数 → 查询参数 → 表单数据
	result := data.NewObjectValue()

	// 从路由参数获取（最低优先级）
	for _, key := range keys {
		if pathVals := collectPathValues(h.source); pathVals != nil {
			if val, exists := pathVals[key]; exists {
				result.SetProperty(key, data.NewStringValue(val))
			}
		}
	}

	// 从查询参数获取
	for _, key := range keys {
		if values, exists := h.source.URL.Query()[key]; exists && len(values) > 0 {
			if len(values) == 1 {
				result.SetProperty(key, data.NewStringValue(values[0]))
			} else {
				result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
			}
		}
	}

	// 从 POST 表单数据获取（最高优先级）
	if h.source.PostForm != nil {
		for _, key := range keys {
			if values, exists := h.source.PostForm[key]; exists && len(values) > 0 {
				if len(values) == 1 {
					result.SetProperty(key, data.NewStringValue(values[0]))
				} else {
					result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
				}
			}
		}
	}

	return result, nil
}

func (h *RequestOnlyMethod) GetName() string            { return "only" }
func (h *RequestOnlyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestOnlyMethod) GetIsStatic() bool          { return false }
func (h *RequestOnlyMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "keys", 0, nil, nil),
	}
}
func (h *RequestOnlyMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "keys", 0, nil),
	}
}
func (h *RequestOnlyMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
