package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DBBeginMethod struct {
	source *sqlsrc.DB
}

func (h *DBBeginMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, err := h.source.Begin()
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewTxClassFrom(ret0), ctx), nil
}

func (h *DBBeginMethod) GetName() string            { return "begin" }
func (h *DBBeginMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBBeginMethod) GetIsStatic() bool          { return true }
func (h *DBBeginMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *DBBeginMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *DBBeginMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
