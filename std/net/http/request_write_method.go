package http

import (
	"io"
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestWriteMethod struct {
	source *httpsrc.Request
}

func (h *RequestWriteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[io.Writer](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.Write(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *RequestWriteMethod) GetName() string            { return "write" }
func (h *RequestWriteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestWriteMethod) GetIsStatic() bool          { return false }
func (h *RequestWriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestWriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestWriteMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
