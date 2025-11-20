package context

import (
	"context"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type WithoutCancelFunction struct{}

func NewWithoutCancelFunction() data.FuncStmt {
	return &WithoutCancelFunction{}
}

func (h *WithoutCancelFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

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
	ret0 := context.WithoutCancel(arg0)
	return data.NewClassValue(NewContextClassFrom(ret0), ctx), nil
}

func (h *WithoutCancelFunction) GetName() string            { return "context\\withoutCancel" }
func (h *WithoutCancelFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *WithoutCancelFunction) GetIsStatic() bool          { return true }
func (h *WithoutCancelFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "parent", 0, nil, nil),
	}
}
func (h *WithoutCancelFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "parent", 0, nil),
	}
}
func (h *WithoutCancelFunction) GetReturnType() data.Types { return data.NewBaseType("void") }
