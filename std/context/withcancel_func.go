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
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
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
