package context

import (
	"context"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"time"
)

type WithTimeoutCauseFunction struct{}

func NewWithTimeoutCauseFunction() data.FuncStmt {
	return &WithTimeoutCauseFunction{}
}

func (h *WithTimeoutCauseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

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

	var arg0 context.Context
	switch v := a0.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			if src := p.GetSource(); src != nil {
				arg0 = src.(context.Context)
			}
		} else {
			return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 0"))
		}
	case *data.AnyValue:
		if v.Value != nil {
			arg0 = v.Value.(context.Context)
		}
	default:
		return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 0"))
	}
	arg1Int, err := a1.(*data.IntValue).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	arg1 := time.Duration(arg1Int)
	var arg2 error
	switch v := a2.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			if src := p.GetSource(); src != nil {
				arg2 = src.(error)
			}
		} else {
			return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 2"))
		}
	case *data.AnyValue:
		if v.Value != nil {
			arg2 = v.Value.(error)
		}
	default:
		return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 2"))
	}
	ret0, ret1 := context.WithTimeoutCause(arg0, arg1, arg2)
	return data.NewAnyValue([]any{ret0, ret1}), nil
}

func (h *WithTimeoutCauseFunction) GetName() string            { return "context\\withTimeoutCause" }
func (h *WithTimeoutCauseFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithTimeoutCauseFunction) GetIsStatic() bool          { return true }
func (h *WithTimeoutCauseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
		node.NewParameter(nil, "timeout", 1, nil, nil),
		node.NewParameter(nil, "cause", 2, nil, nil),
	}
}
func (h *WithTimeoutCauseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
		node.NewVariable(nil, "timeout", 1, nil),
		node.NewVariable(nil, "cause", 2, nil),
	}
}
func (h *WithTimeoutCauseFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
