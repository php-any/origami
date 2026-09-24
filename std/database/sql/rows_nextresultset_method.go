package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
)

type RowsNextResultSetMethod struct {
	source *sqlsrc.Rows
}

func (h *RowsNextResultSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.NextResultSet()
	return data.NewBoolValue(ret0), nil
}

func (h *RowsNextResultSetMethod) GetName() string            { return "nextResultSet" }
func (h *RowsNextResultSetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsNextResultSetMethod) GetIsStatic() bool          { return true }
var rowsNextResultSetMethodGetParams = []data.GetValue{}

func (h *RowsNextResultSetMethod) GetParams() []data.GetValue {
	return rowsNextResultSetMethodGetParams
}

var rowsNextResultSetMethodGetVariables = []data.Variable{}

func (h *RowsNextResultSetMethod) GetVariables() []data.Variable {
	return rowsNextResultSetMethodGetVariables
}

func (h *RowsNextResultSetMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
