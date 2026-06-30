package http

import (
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestInputMethod 获取输入数据
type RequestInputMethod struct {
	source *httpsrc.Request
}

func (h *RequestInputMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	// 检查是否有参数
	_, hasKey := ctx.GetIndexValue(0)

	// 如果没有参数，返回所有输入数据（路由参数 + 查询参数 + 表单数据）
	if !hasKey {
		result := data.NewObjectValue()

		// 添加路由参数（最低优先级）
		for key, val := range collectPathValues(h.source) {
			result.SetProperty(key, data.NewStringValue(val))
		}

		// 添加查询参数
		for key, values := range h.source.URL.Query() {
			if len(values) == 1 {
				result.SetProperty(key, data.NewStringValue(values[0]))
			} else {
				result.SetProperty(key, data.NewStringValue(strings.Join(values, ",")))
			}
		}

		// 添加 POST 表单数据（最高优先级）
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

	// 优先从 POST 表单数据获取
	if h.source.PostForm != nil {
		if values, exists := h.source.PostForm[param0]; exists && len(values) > 0 {
			return data.NewStringValue(values[0]), nil
		}
	}

	// 然后从查询参数获取
	if values, exists := h.source.URL.Query()[param0]; exists && len(values) > 0 {
		return data.NewStringValue(values[0]), nil
	}

	// 最后从路由参数获取
	if pathVals := collectPathValues(h.source); pathVals != nil {
		if val, exists := pathVals[param0]; exists {
			return data.NewStringValue(val), nil
		}
	}

	return data.NewAnyValue(nil), nil
}

func (h *RequestInputMethod) GetName() string            { return "input" }
func (h *RequestInputMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestInputMethod) GetIsStatic() bool          { return false }
func (h *RequestInputMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}
func (h *RequestInputMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}
func (h *RequestInputMethod) GetReturnType() data.Types {
	return data.NewUnionType([]data.Types{
		data.NewBaseType("array"),
		data.NewBaseType("string"),
		data.NewBaseType("null"),
	})
}
