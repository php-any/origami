package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterSuccessMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterSuccessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	var payload data.Value
	if v, ok := ctx.GetIndexValue(0); ok {
		payload = v
	}

	message := "success"
	if v, ok := ctx.GetIndexValue(1); ok {
		message = v.AsString()
	}

	status := httpsrc.StatusOK
	if _, ok := ctx.GetIndexValue(2); ok {
		code, err := utils.ConvertFromIndex[int](ctx, 2)
		if err != nil {
			return nil, utils.NewThrowf("status 参数转换失败: %v", err)
		}
		status = code
	}

	if err := writeFormattedResponse(h.w, ctx, status, message, payload); err != nil {
		return nil, utils.NewThrow(err)
	}
	return responseSelf(h.w, ctx)
}

func (h *ResponseWriterSuccessMethod) GetName() string            { return "success" }
func (h *ResponseWriterSuccessMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterSuccessMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterSuccessMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, data.NewNullValue(), nil),
		node.NewParameter(nil, "message", 1, data.NewStringValue("success"), data.NewBaseType("string")),
		node.NewParameter(nil, "status", 2, data.NewIntValue(httpsrc.StatusOK), data.NewBaseType("int")),
	}
}
func (h *ResponseWriterSuccessMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, nil),
		node.NewVariable(nil, "message", 1, nil),
		node.NewVariable(nil, "status", 2, nil),
	}
}
func (h *ResponseWriterSuccessMethod) GetReturnType() data.Types {
	return data.Class{Name: "Net\\Http\\Response"}
}
