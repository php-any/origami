package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type StmtQueryMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtQueryMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0 := make([]any, 0)
	for _, v := range a0.(*data.ArrayValue).Value {
		arg0 = append(arg0, ConvertValueToGoType(v))
	}

	ret0, err := h.source.Query(arg0...)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewRowsClassFrom(ret0), ctx), nil
}

func (h *StmtQueryMethod) GetName() string            { return "query" }
func (h *StmtQueryMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtQueryMethod) GetIsStatic() bool          { return true }
func (h *StmtQueryMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}

func (h *StmtQueryMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}

func (h *StmtQueryMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
