package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterErrorMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterErrorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	message := "error"
	if v, ok := ctx.GetIndexValue(0); ok {
		message = v.AsString()
	}

	code := httpsrc.StatusInternalServerError
	if _, ok := ctx.GetIndexValue(1); ok {
		status, err := utils.ConvertFromIndex[int](ctx, 1)
		if err != nil {
			return nil, utils.NewThrowf("code 参数转换失败: %v", err)
		}
		code = status
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

func (h *ResponseWriterErrorMethod) GetName() string            { return "error" }
func (h *ResponseWriterErrorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterErrorMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterErrorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "message", 0, data.NewStringValue("error"), data.NewBaseType("string")),
		node.NewParameter(nil, "code", 1, data.NewIntValue(httpsrc.StatusInternalServerError), data.NewBaseType("int")),
		node.NewParameter(nil, "data", 2, data.NewNullValue(), nil),
	}
}
func (h *ResponseWriterErrorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "message", 0, nil),
		node.NewVariable(nil, "code", 1, nil),
		node.NewVariable(nil, "data", 2, nil),
	}
}
func (h *ResponseWriterErrorMethod) GetReturnType() data.Types {
	return data.Class{Name: "Net\\Http\\Response"}
}
