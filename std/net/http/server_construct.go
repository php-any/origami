package http

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ServerConstructMethod struct {
	source *ServerClass
}

func (h *ServerConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}
	param1, err := utils.ConvertFromIndex[int](ctx, 1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	h.source.Host = param0
	h.source.Port = param1
	return nil, nil
}

func (h *ServerConstructMethod) GetName() string            { return "__construct" }
func (h *ServerConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerConstructMethod) GetIsStatic() bool          { return false }
func (h *ServerConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "host", 0, data.NewStringValue("0.0.0.0"), data.String{}),
		node.NewParameter(nil, "port", 1, data.NewIntValue(80), data.Int{}),
	}
}
func (h *ServerConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "host", 0, data.String{}),
		node.NewVariable(nil, "port", 1, data.Int{}),
	}
}
func (h *ServerConstructMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
