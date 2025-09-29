package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	httpsrc "net/http"
)

type RequestSetPathValueMethod struct {
	source *httpsrc.Request
}

func (h *RequestSetPathValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}
	param1, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	h.source.SetPathValue(param0, param1)
	return nil, nil
}

func (h *RequestSetPathValueMethod) GetName() string            { return "setPathValue" }
func (h *RequestSetPathValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestSetPathValueMethod) GetIsStatic() bool          { return false }
func (h *RequestSetPathValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *RequestSetPathValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *RequestSetPathValueMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
