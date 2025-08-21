package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DBExecMethod struct {
	source *sqlsrc.DB
}

func (h *DBExecMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	arg0 := a0.(*data.StringValue).AsString()
	arg1 := *a1.(*data.ArrayValue)

	ret0, err := h.source.Exec(arg0, arg1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewAnyValue(ret0), nil
}

func (h *DBExecMethod) GetName() string            { return "exec" }
func (h *DBExecMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBExecMethod) GetIsStatic() bool          { return true }
func (h *DBExecMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "query", 0, nil, nil),
		node.NewParameter(nil, "args", 1, nil, nil),
	}
}

func (h *DBExecMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "query", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *DBExecMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
