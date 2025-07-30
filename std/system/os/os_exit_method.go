package os

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type OSExitMethod struct {
	source *OS
}

func (h *OSExitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	i, _ := a0.(*data.IntValue).AsInt()
	h.source.Exit(i)
	return nil, nil
}

func (h *OSExitMethod) GetName() string {
	return "exit"
}

func (h *OSExitMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *OSExitMethod) GetIsStatic() bool {
	return false
}

func (h *OSExitMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}

func (h *OSExitMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *OSExitMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
