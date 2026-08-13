package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"

	"time"
)

type DBSetConnMaxLifetimeMethod struct {
	source *sqlsrc.DB
}

func (h *DBSetConnMaxLifetimeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0Int, err := a0.(*data.IntValue).AsInt()
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	arg0 := time.Duration(arg0Int)

	h.source.SetConnMaxLifetime(arg0)
	return nil, nil
}

func (h *DBSetConnMaxLifetimeMethod) GetName() string            { return "setConnMaxLifetime" }
func (h *DBSetConnMaxLifetimeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBSetConnMaxLifetimeMethod) GetIsStatic() bool          { return true }
func (h *DBSetConnMaxLifetimeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "d", 0, nil, nil),
	}
}

func (h *DBSetConnMaxLifetimeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "d", 0, nil),
	}
}

func (h *DBSetConnMaxLifetimeMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
