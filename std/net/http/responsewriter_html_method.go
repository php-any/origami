package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterHtmlMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterHtmlMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	content, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrowf("html 方法缺少内容参数")
	}

	if _, hasStatus := ctx.GetIndexValue(1); hasStatus {
		code, err := utils.ConvertFromIndex[int](ctx, 1)
		if err != nil {
			return nil, utils.NewThrowf("status 参数转换失败: %v", err)
		}
		h.w.SetStatus(code)
	}

	if err := h.w.WriteHTML([]byte(content.AsString())); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ResponseWriterHtmlMethod) GetName() string            { return "html" }
func (h *ResponseWriterHtmlMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterHtmlMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterHtmlMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "content", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "statusCode", 1, data.NewIntValue(200), data.NewBaseType("int")),
	}
}
func (h *ResponseWriterHtmlMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "content", 0, nil),
		node.NewVariable(nil, "statusCode", 1, nil),
	}
}
func (h *ResponseWriterHtmlMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
