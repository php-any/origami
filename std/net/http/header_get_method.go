package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	httpsrc "net/http"
)

type HeaderGetMethod struct {
	source *httpsrc.Header
}

func (h *HeaderGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	ret0 := h.source.Get(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *HeaderGetMethod) GetName() string            { return "get" }
func (h *HeaderGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *HeaderGetMethod) GetIsStatic() bool          { return false }
func (h *HeaderGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *HeaderGetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *HeaderGetMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
