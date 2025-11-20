package context

import (
	"context"
	"errors"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type WithDeadlineCauseFunction struct{}

func NewWithDeadlineCauseFunction() data.FuncStmt {
	return &WithDeadlineCauseFunction{}
}

func (h *WithDeadlineCauseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 1"))
	}

	a2, ok := ctx.GetIndexValue(2)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 2"))
	}

	var arg0 context.Context
	switch v := a0.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			if src := p.GetSource(); src != nil {
				arg0 = src.(context.Context)
			}
		} else {
			return nil, utils.NewThrow(errors.New("参数类型不支持, index: 0"))
		}
	case *data.AnyValue:
		if v.Value != nil {
			arg0 = v.Value.(context.Context)
		}
	default:
		return nil, utils.NewThrow(errors.New("参数类型不支持, index: 0"))
	}
	var arg1 time.Time
	switch v := a1.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			arg1 = p.GetSource().(time.Time)
		} else {
			return nil, utils.NewThrow(errors.New("参数类型不支持, index: 1"))
		}
	case *data.AnyValue:
		arg1 = v.Value.(time.Time)
	default:
		return nil, utils.NewThrow(errors.New("参数类型不支持, index: 1"))
	}
	var arg2 error
	switch v := a2.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			if src := p.GetSource(); src != nil {
				arg2 = src.(error)
			}
		} else {
			return nil, utils.NewThrow(errors.New("参数类型不支持, index: 2"))
		}
	case *data.AnyValue:
		if v.Value != nil {
			arg2 = v.Value.(error)
		}
	default:
		return nil, utils.NewThrow(errors.New("参数类型不支持, index: 2"))
	}
	ret0, ret1 := context.WithDeadlineCause(arg0, arg1, arg2)
	return data.NewAnyValue([]any{ret0, ret1}), nil
}

func (h *WithDeadlineCauseFunction) GetName() string            { return "context\\withDeadlineCause" }
func (h *WithDeadlineCauseFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithDeadlineCauseFunction) GetIsStatic() bool          { return true }
func (h *WithDeadlineCauseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
		node.NewParameter(nil, "d", 1, nil, nil),
		node.NewParameter(nil, "cause", 2, nil, nil),
	}
}
func (h *WithDeadlineCauseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
		node.NewVariable(nil, "d", 1, nil),
		node.NewVariable(nil, "cause", 2, nil),
	}
}
func (h *WithDeadlineCauseFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
