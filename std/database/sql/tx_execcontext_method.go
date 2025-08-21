package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type TxExecContextMethod struct {
	source *sqlsrc.Tx
}

func (h *TxExecContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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

	arg0 := a0.(*data.AnyValue).Value.(context.Context)
	arg1 := a1.(*data.StringValue).AsString()
	arg2 := *a2.(*data.ArrayValue)

	ret0, err := h.source.ExecContext(arg0, arg1, arg2)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewAnyValue(ret0), nil
}

func (h *TxExecContextMethod) GetName() string            { return "execContext" }
func (h *TxExecContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxExecContextMethod) GetIsStatic() bool          { return true }
func (h *TxExecContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "query", 1, nil, nil),
		node.NewParameter(nil, "args", 2, nil, nil),
	}
}

func (h *TxExecContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "query", 1, nil),
		node.NewVariable(nil, "args", 2, nil),
	}
}

func (h *TxExecContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
