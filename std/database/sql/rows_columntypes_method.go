package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
)

type RowsColumnTypesMethod struct {
	source *sqlsrc.Rows
}

func (h *RowsColumnTypesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, err := h.source.ColumnTypes()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewAnyValue(ret0), nil
}

func (h *RowsColumnTypesMethod) GetName() string            { return "columnTypes" }
func (h *RowsColumnTypesMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsColumnTypesMethod) GetIsStatic() bool          { return true }
func (h *RowsColumnTypesMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *RowsColumnTypesMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *RowsColumnTypesMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
