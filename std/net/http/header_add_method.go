package http

import (
	"fmt"
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type HeaderAddMethod struct {
	source *httpsrc.Header
}

func (h *HeaderAddMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}
	param1, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	h.source.Add(param0, param1)
	return nil, nil
}

func (h *HeaderAddMethod) GetName() string            { return "add" }
func (h *HeaderAddMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *HeaderAddMethod) GetIsStatic() bool          { return false }
func (h *HeaderAddMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *HeaderAddMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *HeaderAddMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
