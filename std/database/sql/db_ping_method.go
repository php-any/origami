package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
)

type DBPingMethod struct {
	source *sqlsrc.DB
}

func (h *DBPingMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Ping(); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *DBPingMethod) GetName() string            { return "ping" }
func (h *DBPingMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBPingMethod) GetIsStatic() bool          { return true }
func (h *DBPingMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *DBPingMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *DBPingMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
