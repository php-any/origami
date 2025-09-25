package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	httpsrc "net/http"
)

type ServeMuxHandleFuncMethod struct {
	source *httpsrc.ServeMux
}

func (h *ServeMuxHandleFuncMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}
	param1, err := utils.ConvertFromIndex[func(httpsrc.ResponseWriter, *httpsrc.Request)](ctx, 1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	h.source.HandleFunc(param0, param1)
	return nil, nil
}

func (h *ServeMuxHandleFuncMethod) GetName() string            { return "handleFunc" }
func (h *ServeMuxHandleFuncMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServeMuxHandleFuncMethod) GetIsStatic() bool          { return true }
func (h *ServeMuxHandleFuncMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *ServeMuxHandleFuncMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *ServeMuxHandleFuncMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
