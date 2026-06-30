package http

import (
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestPostMethod 获取 POST 表单数据
// 不带参数时返回所有 POST 数据，带参数时返回指定键的值
type RequestPostMethod struct {
	source *httpsrc.Request
}

func (h *RequestPostMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	// 检查是否有参数
	_, hasKey := ctx.GetIndexValue(0)

	// 如果没有参数，返回所有 POST 表单数据
	if !hasKey {
		result := data.NewObjectValue()

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

	// 如果有参数，返回指定键的值
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	if h.source.PostForm != nil {
		if values, exists := h.source.PostForm[param0]; exists && len(values) > 0 {
			return data.NewStringValue(values[0]), nil
		}
	}

	return data.NewAnyValue(nil), nil
}

func (h *RequestPostMethod) GetName() string            { return "post" }
func (h *RequestPostMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestPostMethod) GetIsStatic() bool          { return false }
func (h *RequestPostMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}
func (h *RequestPostMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}
func (h *RequestPostMethod) GetReturnType() data.Types {
	return data.NewUnionType([]data.Types{
		data.NewBaseType("array"),
		data.NewBaseType("string"),
		data.NewBaseType("null"),
	})
}
