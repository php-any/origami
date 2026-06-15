package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterStatusMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterStatusMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	code, err := utils.ConvertFromIndex[int](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}
	h.w.SetStatus(code)
	return responseSelf(h.w, ctx)
}

func (h *ResponseWriterStatusMethod) GetName() string            { return "status" }
func (h *ResponseWriterStatusMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterStatusMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterStatusMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "statusCode", 0, nil, data.NewBaseType("int")),
	}
}
func (h *ResponseWriterStatusMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "statusCode", 0, nil),
	}
}
func (h *ResponseWriterStatusMethod) GetReturnType() data.Types {
	return data.Class{Name: "Net\\Http\\Response"}
}
