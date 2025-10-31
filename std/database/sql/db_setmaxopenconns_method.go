package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DBSetMaxOpenConnsMethod struct {
	source *sqlsrc.DB
}

func (h *DBSetMaxOpenConnsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0, err := a0.(*data.IntValue).AsInt()
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	h.source.SetMaxOpenConns(arg0)
	return nil, nil
}

func (h *DBSetMaxOpenConnsMethod) GetName() string            { return "setMaxOpenConns" }
func (h *DBSetMaxOpenConnsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBSetMaxOpenConnsMethod) GetIsStatic() bool          { return true }
func (h *DBSetMaxOpenConnsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "n", 0, nil, nil),
	}
}

func (h *DBSetMaxOpenConnsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "n", 0, nil),
	}
}

func (h *DBSetMaxOpenConnsMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
