package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type RowErrMethod struct {
	source *sqlsrc.Row
}

func (h *RowErrMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Err(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *RowErrMethod) GetName() string            { return "err" }
func (h *RowErrMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowErrMethod) GetIsStatic() bool          { return true }
func (h *RowErrMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *RowErrMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *RowErrMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
