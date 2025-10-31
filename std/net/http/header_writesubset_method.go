package http

import (
	"io"
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type HeaderWriteSubsetMethod struct {
	source *httpsrc.Header
}

func (h *HeaderWriteSubsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[io.Writer](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}
	param1, err := utils.ConvertFromIndex[map[string]bool](ctx, 1)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.WriteSubset(param0, param1)
	return data.NewAnyValue(ret0), nil
}

func (h *HeaderWriteSubsetMethod) GetName() string            { return "writeSubset" }
func (h *HeaderWriteSubsetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *HeaderWriteSubsetMethod) GetIsStatic() bool          { return false }
func (h *HeaderWriteSubsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *HeaderWriteSubsetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *HeaderWriteSubsetMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
