package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type TxPrepareMethod struct {
	source *sqlsrc.Tx
}

func (h *TxPrepareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0 := a0.(*data.StringValue).AsString()

	ret0, err := h.source.Prepare(arg0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewStmtClassFrom(ret0), ctx), nil
}

func (h *TxPrepareMethod) GetName() string            { return "prepare" }
func (h *TxPrepareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxPrepareMethod) GetIsStatic() bool          { return true }
func (h *TxPrepareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "query", 0, nil, nil),
	}
}

func (h *TxPrepareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "query", 0, nil),
	}
}

func (h *TxPrepareMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
