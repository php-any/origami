package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestHasMethod 检查是否有指定的输入数据
type RequestHasMethod struct {
	source *httpsrc.Request
}

func (h *RequestHasMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewBoolValue(false), nil
	}

	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	// 检查表单数据
	if h.source.Form != nil {
		if values, exists := h.source.Form[param0]; exists && len(values) > 0 && values[0] != "" {
			return data.NewBoolValue(true), nil
		}
	}

	// 检查查询参数
	if values, exists := h.source.URL.Query()[param0]; exists && len(values) > 0 && values[0] != "" {
		return data.NewBoolValue(true), nil
	}

	return data.NewBoolValue(false), nil
}

func (h *RequestHasMethod) GetName() string            { return "has" }
func (h *RequestHasMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestHasMethod) GetIsStatic() bool          { return false }
func (h *RequestHasMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}
func (h *RequestHasMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}
func (h *RequestHasMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
