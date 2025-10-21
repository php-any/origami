package driver

import (
	driversrc "database/sql/driver"
	"github.com/php-any/origami/data"
)

type RowsCloseMethod struct {
	source driversrc.Rows
}

func (h *RowsCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Close(); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *RowsCloseMethod) GetName() string            { return "close" }
func (h *RowsCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsCloseMethod) GetIsStatic() bool          { return true }
func (h *RowsCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *RowsCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *RowsCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
