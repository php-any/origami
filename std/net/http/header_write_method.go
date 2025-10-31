package http

import (
	"io"
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type HeaderWriteMethod struct {
	source *httpsrc.Header
}

func (h *HeaderWriteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[io.Writer](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.Write(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *HeaderWriteMethod) GetName() string            { return "write" }
func (h *HeaderWriteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *HeaderWriteMethod) GetIsStatic() bool          { return false }
func (h *HeaderWriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *HeaderWriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *HeaderWriteMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
