package http

import (
	"fmt"
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestHeaderMethod 获取请求头
type RequestHeaderMethod struct {
	source *httpsrc.Request
}

func (h *RequestHeaderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	// 检查是否有参数
	_, hasKey := ctx.GetIndexValue(0)

	// 如果没有参数，返回所有请求头
	if !hasKey {
		result := data.NewObjectValue()
		for key, values := range h.source.Header {
			if len(values) > 0 {
				result.SetProperty(key, data.NewStringValue(values[0]))
			}
		}
		return result, nil
	}

	// 如果有参数，返回指定请求头的值
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	value := h.source.Header.Get(param0)
	return data.NewStringValue(value), nil
}

func (h *RequestHeaderMethod) GetName() string            { return "header" }
func (h *RequestHeaderMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestHeaderMethod) GetIsStatic() bool          { return false }
func (h *RequestHeaderMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}
func (h *RequestHeaderMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}
func (h *RequestHeaderMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
