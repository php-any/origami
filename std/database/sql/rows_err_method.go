package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type RowsErrMethod struct {
	source *sqlsrc.Rows
}

func (h *RowsErrMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Err(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *RowsErrMethod) GetName() string            { return "err" }
func (h *RowsErrMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsErrMethod) GetIsStatic() bool          { return true }
var rowsErrMethodGetParams = []data.GetValue{}

func (h *RowsErrMethod) GetParams() []data.GetValue {
	return rowsErrMethodGetParams
}

var rowsErrMethodGetVariables = []data.Variable{}

func (h *RowsErrMethod) GetVariables() []data.Variable {
	return rowsErrMethodGetVariables
}

func (h *RowsErrMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
