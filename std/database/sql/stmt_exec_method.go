package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type StmtExecMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtExecMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0 := make([]any, 0)
	for _, v := range a0.(*data.ArrayValue).Value {
		arg0 = append(arg0, v)
	}

	ret0, err := h.source.Exec(arg0...)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewResultClassFrom(ret0), ctx), nil
}

func (h *StmtExecMethod) GetName() string            { return "exec" }
func (h *StmtExecMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtExecMethod) GetIsStatic() bool          { return true }
func (h *StmtExecMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}

func (h *StmtExecMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}

func (h *StmtExecMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
