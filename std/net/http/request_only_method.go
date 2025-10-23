package http

import (
	"fmt"
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
			return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return data.NewObjectValue(), nil
	}

	// 只返回指定的键
	result := data.NewObjectValue()

	// 从表单数据获取
	if h.source.Form != nil {
		for _, key := range keys {
			if values, exists := h.source.Form[key]; exists && len(values) > 0 {
				if len(values) == 1 {
					result.SetProperty(key, data.NewStringValue(values[0]))
				} else {
					result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
				}
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
func (h *RequestOnlyMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
