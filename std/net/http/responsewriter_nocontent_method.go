package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterNoContentMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterNoContentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	code := httpsrc.StatusNoContent
	if _, ok := ctx.GetIndexValue(0); ok {
		var err error
		code, err = utils.ConvertFromIndex[int](ctx, 0)
		if err != nil {
			return nil, utils.NewThrowf("参数转换失败: %v", err)
		}
	}

	h.w.NoContent(code)
	return nil, nil
}

func (h *ResponseWriterNoContentMethod) GetName() string            { return "noContent" }
func (h *ResponseWriterNoContentMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterNoContentMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterNoContentMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "statusCode", 0, data.NewIntValue(httpsrc.StatusNoContent), data.NewBaseType("int")),
	}
}
func (h *ResponseWriterNoContentMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "statusCode", 0, nil),
	}
}
func (h *ResponseWriterNoContentMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
