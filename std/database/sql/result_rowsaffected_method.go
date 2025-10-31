package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type ResultRowsAffectedMethod struct {
	source sqlsrc.Result
}

func (h *ResultRowsAffectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, err := h.source.RowsAffected()
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewIntValue(int(ret0)), nil
}

func (h *ResultRowsAffectedMethod) GetName() string            { return "rowsAffected" }
func (h *ResultRowsAffectedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResultRowsAffectedMethod) GetIsStatic() bool          { return true }
func (h *ResultRowsAffectedMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ResultRowsAffectedMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *ResultRowsAffectedMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
