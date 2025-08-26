package context

import (
	contextsrc "context"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ContextValueMethod struct {
	source contextsrc.Context
}

func (h *ContextValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0 := a0.(*data.AnyValue).Value

	ret0 := h.source.Value(arg0)
	return data.NewAnyValue(ret0), nil
}

func (h *ContextValueMethod) GetName() string            { return "value" }
func (h *ContextValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ContextValueMethod) GetIsStatic() bool          { return true }
func (h *ContextValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}

func (h *ContextValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}

func (h *ContextValueMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
