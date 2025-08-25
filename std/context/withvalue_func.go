package context

import (
	"context"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type WithValueFunction struct{}

func NewWithValueFunction() data.FuncStmt {
	return &WithValueFunction{}
}

func (h *WithValueFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	a2, ok := ctx.GetIndexValue(2)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 2"))
	}

	arg0 := a0.(*data.AnyValue).Value.(context.Context)
	arg1 := a1.(*data.AnyValue).Value
	arg2 := a2.(*data.AnyValue).Value
	ret0 := context.WithValue(arg0, arg1, arg2)
	return data.NewAnyValue(ret0), nil
}

func (h *WithValueFunction) GetName() string            { return "context\\withValue" }
func (h *WithValueFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithValueFunction) GetIsStatic() bool          { return true }
func (h *WithValueFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
		node.NewParameter(nil, "key", 1, nil, nil),
		node.NewParameter(nil, "val", 2, nil, nil),
	}
}
func (h *WithValueFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
		node.NewVariable(nil, "key", 1, nil),
		node.NewVariable(nil, "val", 2, nil),
	}
}
func (h *WithValueFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
