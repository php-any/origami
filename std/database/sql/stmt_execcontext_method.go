package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type StmtExecContextMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtExecContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	arg0 := a0.(*data.AnyValue).Value.(context.Context)
	arg1 := *a1.(*data.ArrayValue)

	ret0, err := h.source.ExecContext(arg0, arg1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewAnyValue(ret0), nil
}

func (h *StmtExecContextMethod) GetName() string            { return "execContext" }
func (h *StmtExecContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtExecContextMethod) GetIsStatic() bool          { return true }
func (h *StmtExecContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "args", 1, nil, nil),
	}
}

func (h *StmtExecContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *StmtExecContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
