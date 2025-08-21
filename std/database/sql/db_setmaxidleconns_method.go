package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DBSetMaxIdleConnsMethod struct {
	source *sqlsrc.DB
}

func (h *DBSetMaxIdleConnsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0, err := a0.(*data.IntValue).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	h.source.SetMaxIdleConns(arg0)
	return nil, nil
}

func (h *DBSetMaxIdleConnsMethod) GetName() string            { return "setMaxIdleConns" }
func (h *DBSetMaxIdleConnsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBSetMaxIdleConnsMethod) GetIsStatic() bool          { return true }
func (h *DBSetMaxIdleConnsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "n", 0, nil, nil),
	}
}

func (h *DBSetMaxIdleConnsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "n", 0, nil),
	}
}

func (h *DBSetMaxIdleConnsMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
