package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterWriteHeaderMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterWriteHeaderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	code, err := utils.ConvertFromIndex[int](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}
	h.w.WriteHeader(code)
	return nil, nil
}

func (h *ResponseWriterWriteHeaderMethod) GetName() string            { return "writeHeader" }
func (h *ResponseWriterWriteHeaderMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterWriteHeaderMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterWriteHeaderMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *ResponseWriterWriteHeaderMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *ResponseWriterWriteHeaderMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
