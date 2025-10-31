package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterHeaderMethod struct {
	source httpsrc.ResponseWriter
}

func (h *ResponseWriterHeaderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否有参数
	_, hasKey := ctx.GetIndexValue(0)

	// 如果没有参数，返回 Header 对象
	if !hasKey {
		ret0 := h.source.Header()
		return data.NewProxyValue(NewHeaderClassFrom(&ret0), ctx), nil
	}

	// 如果有参数，设置响应头
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	param1, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	h.source.Header().Set(param0, param1)
	return nil, nil
}

func (h *ResponseWriterHeaderMethod) GetName() string            { return "header" }
func (h *ResponseWriterHeaderMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterHeaderMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterHeaderMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}
func (h *ResponseWriterHeaderMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
		node.NewVariable(nil, "value", 1, nil),
	}
}
func (h *ResponseWriterHeaderMethod) GetReturnType() data.Types {
	return data.NewBaseType("Net\\Http\\Header")
}
