package context

import (
	"context"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type WithCancelFunction struct{}

func NewWithCancelFunction() data.FuncStmt {
	return &WithCancelFunction{}
}

func (h *WithCancelFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0 := a0.(*data.AnyValue).Value.(context.Context)
	ret0, ret1 := context.WithCancel(arg0)
	return data.NewAnyValue([]any{ret0, ret1}), nil
}

func (h *WithCancelFunction) GetName() string            { return "context\\withCancel" }
func (h *WithCancelFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithCancelFunction) GetIsStatic() bool          { return true }
func (h *WithCancelFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
	}
}
func (h *WithCancelFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
	}
}
func (h *WithCancelFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
