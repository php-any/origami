package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type StmtQueryRowMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtQueryRowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0 := *a0.(*data.ArrayValue)

	ret0 := h.source.QueryRow(arg0)
	return data.NewClassValue(NewRowClassFrom(ret0), ctx), nil
}

func (h *StmtQueryRowMethod) GetName() string            { return "queryRow" }
func (h *StmtQueryRowMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtQueryRowMethod) GetIsStatic() bool          { return true }
func (h *StmtQueryRowMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "args", 0, nil, nil),
	}
}

func (h *StmtQueryRowMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}

func (h *StmtQueryRowMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
