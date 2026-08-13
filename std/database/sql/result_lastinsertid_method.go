package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type ResultLastInsertIdMethod struct {
	source sqlsrc.Result
}

func (h *ResultLastInsertIdMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, err := h.source.LastInsertId()
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewIntValue(int(ret0)), nil
}

func (h *ResultLastInsertIdMethod) GetName() string            { return "lastInsertId" }
func (h *ResultLastInsertIdMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResultLastInsertIdMethod) GetIsStatic() bool          { return true }
func (h *ResultLastInsertIdMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ResultLastInsertIdMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *ResultLastInsertIdMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
