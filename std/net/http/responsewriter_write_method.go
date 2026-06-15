package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterWriteMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterWriteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrowf("write方法缺少参数: %v", 0)
	}

	n, err := h.w.Write([]byte(param0.AsString()))
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	_ = n
	return responseSelf(h.w, ctx)
}

func (h *ResponseWriterWriteMethod) GetName() string            { return "write" }
func (h *ResponseWriterWriteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterWriteMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterWriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *ResponseWriterWriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *ResponseWriterWriteMethod) GetReturnType() data.Types {
	return data.Class{Name: "Net\\Http\\Response"}
}
