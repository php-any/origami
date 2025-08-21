package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ConnBeginTxMethod struct {
	source *sqlsrc.Conn
}

func (h *ConnBeginTxMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	arg0 := a0.(*data.AnyValue).Value.(context.Context)
	arg1Class := a1.(*data.ClassValue).Class.(*TxOptionsClass)
	arg1 := arg1Class.source

	ret0, err := h.source.BeginTx(arg0, arg1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewTxClassFrom(ret0), ctx), nil
}

func (h *ConnBeginTxMethod) GetName() string            { return "beginTx" }
func (h *ConnBeginTxMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnBeginTxMethod) GetIsStatic() bool          { return true }
func (h *ConnBeginTxMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "opts", 1, nil, nil),
	}
}

func (h *ConnBeginTxMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "opts", 1, nil),
	}
}

func (h *ConnBeginTxMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
