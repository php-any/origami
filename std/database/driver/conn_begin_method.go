package driver

import (
	driversrc "database/sql/driver"

	"github.com/php-any/origami/data"
)

type ConnBeginMethod struct {
	source driversrc.Conn
}

func (h *ConnBeginMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, err := h.source.Begin()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewTxClassFrom(ret0), ctx), nil
}

func (h *ConnBeginMethod) GetName() string            { return "begin" }
func (h *ConnBeginMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnBeginMethod) GetIsStatic() bool          { return true }
func (h *ConnBeginMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ConnBeginMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *ConnBeginMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
