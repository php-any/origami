package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type StmtQueryRowContextMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtQueryRowContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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

	ret0 := h.source.QueryRowContext(arg0, arg1)
	return data.NewClassValue(NewRowClassFrom(ret0), ctx), nil
}

func (h *StmtQueryRowContextMethod) GetName() string            { return "queryRowContext" }
func (h *StmtQueryRowContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtQueryRowContextMethod) GetIsStatic() bool          { return true }
func (h *StmtQueryRowContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "args", 1, nil, nil),
	}
}

func (h *StmtQueryRowContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *StmtQueryRowContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
