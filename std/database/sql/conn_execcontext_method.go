package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ConnExecContextMethod struct {
	source *sqlsrc.Conn
}

func (h *ConnExecContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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
			arg0 = p.GetSource().(context.Context)
		} else {
			return nil, utils.NewThrow(errors.New("参数类型不支持, index: 0"))
		}
	case *data.AnyValue:
		arg0 = v.Value.(context.Context)
	default:
		return nil, utils.NewThrow(errors.New("参数类型不支持, index: 0"))
	}
	arg1 := a1.(*data.StringValue).AsString()
	arg2 := make([]any, 0)
	for _, v := range a2.(*data.ArrayValue).Value {
		arg2 = append(arg2, ConvertValueToGoType(v))
	}

	ret0, err := h.source.ExecContext(arg0, arg1, arg2...)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewResultClassFrom(ret0), ctx), nil
}

func (h *ConnExecContextMethod) GetName() string            { return "execContext" }
func (h *ConnExecContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnExecContextMethod) GetIsStatic() bool          { return true }
func (h *ConnExecContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "query", 1, nil, nil),
		node.NewParameters(nil, "args", 2, nil, nil),
	}
}

func (h *ConnExecContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "query", 1, nil),
		node.NewVariable(nil, "args", 2, nil),
	}
}

func (h *ConnExecContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
