package http

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ServerGroupMethod struct {
	server *ServerClass
}

func (h *ServerGroupMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	prefix, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("prefix 参数转换失败: %v", err))
	}

	ret := NewServerClassFromGroup(prefix, h.server)

	return data.NewProxyValue(ret, ctx.CreateBaseContext()), nil
}

func (h *ServerGroupMethod) GetName() string            { return "group" }
func (h *ServerGroupMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerGroupMethod) GetIsStatic() bool          { return false }
func (h *ServerGroupMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "prefix", 0, nil, nil),
	}
}
func (h *ServerGroupMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "prefix", 0, nil),
	}
}
func (h *ServerGroupMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
