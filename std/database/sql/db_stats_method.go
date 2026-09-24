package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
)

type DBStatsMethod struct {
	source *sqlsrc.DB
}

func (h *DBStatsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.Stats()
	return data.NewAnyValue(ret0), nil
}

func (h *DBStatsMethod) GetName() string            { return "stats" }
func (h *DBStatsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBStatsMethod) GetIsStatic() bool          { return true }
var dBStatsMethodGetParams = []data.GetValue{}

func (h *DBStatsMethod) GetParams() []data.GetValue {
	return dBStatsMethodGetParams
}

var dBStatsMethodGetVariables = []data.Variable{}

func (h *DBStatsMethod) GetVariables() []data.Variable {
	return dBStatsMethodGetVariables
}

func (h *DBStatsMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
