package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type TxExecMethod struct {
	source *sqlsrc.Tx
}

func (h *TxExecMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	arg0 := a0.(*data.StringValue).AsString()
	arg1 := make([]any, 0)
	for _, v := range a1.(*data.ArrayValue).Value {
		arg1 = append(arg1, v)
	}

	ret0, err := h.source.Exec(arg0, arg1...)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewResultClassFrom(ret0), ctx), nil
}

func (h *TxExecMethod) GetName() string            { return "exec" }
func (h *TxExecMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxExecMethod) GetIsStatic() bool          { return true }
func (h *TxExecMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "query", 0, nil, nil),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (h *TxExecMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "query", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *TxExecMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
