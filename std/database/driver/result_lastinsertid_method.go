package driver

import (
	driversrc "database/sql/driver"
	"github.com/php-any/origami/data"
)

type ResultLastInsertIdMethod struct {
	source driversrc.Result
}

func (h *ResultLastInsertIdMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, err := h.source.LastInsertId()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
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
