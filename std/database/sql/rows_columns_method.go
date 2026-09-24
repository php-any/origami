package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type RowsColumnsMethod struct {
	source *sqlsrc.Rows
}

func (h *RowsColumnsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, err := h.source.Columns()
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewAnyValue(ret0), nil
}

func (h *RowsColumnsMethod) GetName() string            { return "columns" }
func (h *RowsColumnsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsColumnsMethod) GetIsStatic() bool          { return true }
var rowsColumnsMethodGetParams = []data.GetValue{}

func (h *RowsColumnsMethod) GetParams() []data.GetValue {
	return rowsColumnsMethodGetParams
}

var rowsColumnsMethodGetVariables = []data.Variable{}

func (h *RowsColumnsMethod) GetVariables() []data.Variable {
	return rowsColumnsMethodGetVariables
}

func (h *RowsColumnsMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
