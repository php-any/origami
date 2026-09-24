package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type RowsCloseMethod struct {
	source *sqlsrc.Rows
}

func (h *RowsCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Close(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *RowsCloseMethod) GetName() string            { return "close" }
func (h *RowsCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsCloseMethod) GetIsStatic() bool          { return true }
var rowsCloseMethodGetParams = []data.GetValue{}

func (h *RowsCloseMethod) GetParams() []data.GetValue {
	return rowsCloseMethodGetParams
}

var rowsCloseMethodGetVariables = []data.Variable{}

func (h *RowsCloseMethod) GetVariables() []data.Variable {
	return rowsCloseMethodGetVariables
}

func (h *RowsCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
