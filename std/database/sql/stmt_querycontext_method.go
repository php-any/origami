package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type StmtQueryContextMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtQueryContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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

	ret0, err := h.source.QueryContext(arg0, arg1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewRowsClassFrom(ret0), ctx), nil
}

func (h *StmtQueryContextMethod) GetName() string            { return "queryContext" }
func (h *StmtQueryContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtQueryContextMethod) GetIsStatic() bool          { return true }
func (h *StmtQueryContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "args", 1, nil, nil),
	}
}

func (h *StmtQueryContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *StmtQueryContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
