package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"

	"time"
)

type DBSetConnMaxIdleTimeMethod struct {
	source *sqlsrc.DB
}

func (h *DBSetConnMaxIdleTimeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0Int, err := a0.(*data.IntValue).AsInt()
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	arg0 := time.Duration(arg0Int)

	h.source.SetConnMaxIdleTime(arg0)
	return nil, nil
}

func (h *DBSetConnMaxIdleTimeMethod) GetName() string            { return "setConnMaxIdleTime" }
func (h *DBSetConnMaxIdleTimeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBSetConnMaxIdleTimeMethod) GetIsStatic() bool          { return true }
func (h *DBSetConnMaxIdleTimeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "d", 0, nil, nil),
	}
}

func (h *DBSetConnMaxIdleTimeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "d", 0, nil),
	}
}

func (h *DBSetConnMaxIdleTimeMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
