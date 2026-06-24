package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterRedirectMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterRedirectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	url, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	code := httpsrc.StatusFound
	if _, ok := ctx.GetIndexValue(1); ok {
		code, err = utils.ConvertFromIndex[int](ctx, 1)
		if err != nil {
			return nil, utils.NewThrowf("参数转换失败: %v", err)
		}
	}

	h.w.Redirect(url, code)
	return nil, nil
}

func (h *ResponseWriterRedirectMethod) GetName() string            { return "redirect" }
func (h *ResponseWriterRedirectMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterRedirectMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterRedirectMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "url", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "statusCode", 1, data.NewIntValue(httpsrc.StatusFound), data.NewBaseType("int")),
	}
}
func (h *ResponseWriterRedirectMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "url", 0, nil),
		node.NewVariable(nil, "statusCode", 1, nil),
	}
}
func (h *ResponseWriterRedirectMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
