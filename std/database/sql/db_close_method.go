package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
)

type DBCloseMethod struct {
	source *sqlsrc.DB
}

func (h *DBCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Close(); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *DBCloseMethod) GetName() string            { return "close" }
func (h *DBCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBCloseMethod) GetIsStatic() bool          { return true }
func (h *DBCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *DBCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *DBCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
