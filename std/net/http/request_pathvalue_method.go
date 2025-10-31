package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestPathValueMethod struct {
	source *httpsrc.Request
}

func (h *RequestPathValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.PathValue(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *RequestPathValueMethod) GetName() string            { return "pathValue" }
func (h *RequestPathValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestPathValueMethod) GetIsStatic() bool          { return false }
func (h *RequestPathValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestPathValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestPathValueMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
