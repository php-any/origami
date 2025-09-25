package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	httpsrc "net/http"
)

type RequestFormFileMethod struct {
	source *httpsrc.Request
}

func (h *RequestFormFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	ret0, ret1, ret2 := h.source.FormFile(param0)
	return data.NewArrayValue([]data.Value{data.NewAnyValue(ret0), data.NewAnyValue(ret1), data.NewAnyValue(ret2)}), nil
}

func (h *RequestFormFileMethod) GetName() string            { return "formFile" }
func (h *RequestFormFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestFormFileMethod) GetIsStatic() bool          { return true }
func (h *RequestFormFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestFormFileMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestFormFileMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
