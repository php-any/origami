package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterFormatMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterFormatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	code, err := utils.ConvertFromIndex[int](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("code 参数转换失败: %v", err)
	}

	message, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, utils.NewThrowf("message 参数转换失败: %v", err)
	}

	var payload data.Value
	if v, ok := ctx.GetIndexValue(2); ok {
		payload = v
	}

	if err := writeFormattedResponse(h.w, ctx, code, message, payload); err != nil {
		return nil, utils.NewThrow(err)
	}
	return responseSelf(h.w, ctx)
}

func (h *ResponseWriterFormatMethod) GetName() string            { return "format" }
func (h *ResponseWriterFormatMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterFormatMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterFormatMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "code", 0, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "message", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "data", 2, data.NewNullValue(), nil),
	}
}
func (h *ResponseWriterFormatMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "code", 0, nil),
		node.NewVariable(nil, "message", 1, nil),
		node.NewVariable(nil, "data", 2, nil),
	}
}
func (h *ResponseWriterFormatMethod) GetReturnType() data.Types {
	return data.Class{Name: "Net\\Http\\Response"}
}
