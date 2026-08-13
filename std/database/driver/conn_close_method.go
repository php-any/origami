package driver

import (
	driversrc "database/sql/driver"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type ConnCloseMethod struct {
	source driversrc.Conn
}

func (h *ConnCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Close(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ConnCloseMethod) GetName() string            { return "close" }
func (h *ConnCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnCloseMethod) GetIsStatic() bool          { return true }
func (h *ConnCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ConnCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *ConnCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
