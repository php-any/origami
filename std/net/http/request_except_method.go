package http

import (
	"fmt"
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestExceptMethod 排除指定的输入数据
type RequestExceptMethod struct {
	source *httpsrc.Request
}

func (h *RequestExceptMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	// 获取排除的键列表
	excludeKeys := make([]string, 0)
	for i := 0; ; i++ {
		value, exists := ctx.GetIndexValue(i)
		if !exists {
			break
		}
		key, err := utils.Convert[string](value)
		if err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
		}
		excludeKeys = append(excludeKeys, key)
	}

	// 创建排除键的映射
	excludeMap := make(map[string]bool)
	for _, key := range excludeKeys {
		excludeMap[key] = true
	}

	// 合并所有输入数据
	result := data.NewObjectValue()

	// 从查询参数获取
	for key, values := range h.source.URL.Query() {
		if !excludeMap[key] {
			if len(values) == 1 {
				result.SetProperty(key, data.NewStringValue(values[0]))
			} else {
				result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
			}
		}
	}

	// 从表单数据获取
	if h.source.Form != nil {
		for key, values := range h.source.Form {
			if !excludeMap[key] {
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

func (h *RequestExceptMethod) GetName() string            { return "except" }
func (h *RequestExceptMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestExceptMethod) GetIsStatic() bool          { return false }
func (h *RequestExceptMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "keys", 0, nil, nil),
	}
}
func (h *RequestExceptMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "keys", 0, nil),
	}
}
func (h *RequestExceptMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
