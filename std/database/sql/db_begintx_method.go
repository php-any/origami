package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DBBeginTxMethod struct {
	source *sqlsrc.DB
}

func (h *DBBeginTxMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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

func (h *DBBeginTxMethod) GetName() string            { return "beginTx" }
func (h *DBBeginTxMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBBeginTxMethod) GetIsStatic() bool          { return true }
func (h *DBBeginTxMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "opts", 1, nil, nil),
	}
}

func (h *DBBeginTxMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "opts", 1, nil),
	}
}

func (h *DBBeginTxMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
