package context

import (
	"context"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type WithCancelCauseFunction struct{}

func NewWithCancelCauseFunction() data.FuncStmt {
	return &WithCancelCauseFunction{}
}

func (h *WithCancelCauseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
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
	ret0, ret1 := context.WithCancelCause(arg0)
	return data.NewAnyValue([]any{ret0, ret1}), nil
}

func (h *WithCancelCauseFunction) GetName() string            { return "context\\withCancelCause" }
func (h *WithCancelCauseFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithCancelCauseFunction) GetIsStatic() bool          { return true }
func (h *WithCancelCauseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
	}
}
func (h *WithCancelCauseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
	}
}
func (h *WithCancelCauseFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
