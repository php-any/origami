package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterHeaderMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterHeaderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, hasKey := ctx.GetIndexValue(0)
	if !hasKey {
		ret0 := h.w.Header()
		return data.NewProxyValue(NewHeaderClassFrom(&ret0), ctx), nil
	}

	key, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}
	value, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	h.w.SetHeader(key, value)
	return responseSelf(h.w, ctx)
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
