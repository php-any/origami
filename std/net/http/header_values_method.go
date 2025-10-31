package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type HeaderValuesMethod struct {
	source *httpsrc.Header
}

func (h *HeaderValuesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.Values(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *HeaderValuesMethod) GetName() string            { return "values" }
func (h *HeaderValuesMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *HeaderValuesMethod) GetIsStatic() bool          { return false }
func (h *HeaderValuesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *HeaderValuesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *HeaderValuesMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
