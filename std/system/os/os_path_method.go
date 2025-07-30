package os

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type OSPathMethod struct {
	source *OS
}

func (h *OSPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数: paths"))
	}

	return data.NewStringValue(h.source.Path(*a0.(*data.ArrayValue))), nil
}

func (h *OSPathMethod) GetName() string {
	return "path"
}

func (h *OSPathMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *OSPathMethod) GetIsStatic() bool {
	return true
}

func (h *OSPathMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "paths", 0, nil, nil),
	}
}

func (h *OSPathMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "paths", 0, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *OSPathMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
