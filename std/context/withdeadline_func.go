package context

import (
	"context"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"time"
)

type WithDeadlineFunction struct{}

func NewWithDeadlineFunction() data.FuncStmt {
	return &WithDeadlineFunction{}
}

func (h *WithDeadlineFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
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
	var arg1 time.Time
	switch v := a1.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			arg1 = p.GetSource().(time.Time)
		} else {
			return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 1"))
		}
	case *data.AnyValue:
		arg1 = v.Value.(time.Time)
	default:
		return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 1"))
	}
	ret0, ret1 := context.WithDeadline(arg0, arg1)
	return data.NewAnyValue([]any{ret0, ret1}), nil
}

func (h *WithDeadlineFunction) GetName() string            { return "context\\withDeadline" }
func (h *WithDeadlineFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithDeadlineFunction) GetIsStatic() bool          { return true }
func (h *WithDeadlineFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
		node.NewParameter(nil, "d", 1, nil, nil),
	}
}
func (h *WithDeadlineFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
		node.NewVariable(nil, "d", 1, nil),
	}
}
func (h *WithDeadlineFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
