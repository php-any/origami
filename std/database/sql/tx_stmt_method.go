package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type TxStmtMethod struct {
	source *sqlsrc.Tx
}

func (h *TxStmtMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0Class := a0.(*data.ClassValue).Class.(*StmtClass)
	arg0 := arg0Class.source

	ret0 := h.source.Stmt(arg0)
	return data.NewClassValue(NewStmtClassFrom(ret0), ctx), nil
}

func (h *TxStmtMethod) GetName() string            { return "stmt" }
func (h *TxStmtMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxStmtMethod) GetIsStatic() bool          { return true }
func (h *TxStmtMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stmt", 0, nil, nil),
	}
}

func (h *TxStmtMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stmt", 0, nil),
	}
}

func (h *TxStmtMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
