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
	var arg1 interface{}
	switch v := a1.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			if src := p.GetSource(); src != nil {
				arg1 = src.(interface{})
			}
		} else {
			return nil, utils.NewThrow(errors.New("参数类型不支持, index: 1"))
		}
	case *data.AnyValue:
		if v.Value != nil {
			arg1 = v.Value.(interface{})
		}
	default:
		return nil, utils.NewThrow(errors.New("参数类型不支持, index: 1"))
	}
	var arg2 interface{}
	switch v := a2.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			if src := p.GetSource(); src != nil {
				arg2 = src.(interface{})
			}
		} else {
			return nil, utils.NewThrow(errors.New("参数类型不支持, index: 2"))
		}
	case *data.AnyValue:
		if v.Value != nil {
			arg2 = v.Value.(interface{})
		}
	default:
		return nil, utils.NewThrow(errors.New("参数类型不支持, index: 2"))
	}
	ret0 := context.WithValue(arg0, arg1, arg2)
	return data.NewClassValue(NewContextClassFrom(ret0), ctx), nil
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
