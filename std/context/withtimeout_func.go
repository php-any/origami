package context

import (
	"context"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"time"
)

type WithTimeoutFunction struct{}

func NewWithTimeoutFunction() data.FuncStmt {
	return &WithTimeoutFunction{}
}

func (h *WithTimeoutFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	arg0 := a0.(*data.AnyValue).Value.(context.Context)
	arg1Int, err := a1.(*data.IntValue).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	arg1 := int64(arg1Int)
	ret0, ret1 := context.WithTimeout(arg0, time.Duration(arg1))
	return data.NewAnyValue([]any{ret0, ret1}), nil
}

func (h *WithTimeoutFunction) GetName() string            { return "context\\withTimeout" }
func (h *WithTimeoutFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithTimeoutFunction) GetIsStatic() bool          { return true }
func (h *WithTimeoutFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
		node.NewParameter(nil, "timeout", 1, nil, nil),
	}
}
func (h *WithTimeoutFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
		node.NewVariable(nil, "timeout", 1, nil),
	}
}
func (h *WithTimeoutFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
