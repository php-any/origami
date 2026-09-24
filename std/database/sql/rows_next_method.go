package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
)

type RowsNextMethod struct {
	source *sqlsrc.Rows
}

func (h *RowsNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.Next()
	return data.NewBoolValue(ret0), nil
}

func (h *RowsNextMethod) GetName() string            { return "next" }
func (h *RowsNextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsNextMethod) GetIsStatic() bool          { return true }
var rowsNextMethodGetParams = []data.GetValue{}

func (h *RowsNextMethod) GetParams() []data.GetValue {
	return rowsNextMethodGetParams
}

var rowsNextMethodGetVariables = []data.Variable{}

func (h *RowsNextMethod) GetVariables() []data.Variable {
	return rowsNextMethodGetVariables
}

func (h *RowsNextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
