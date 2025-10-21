package driver

import (
	driversrc "database/sql/driver"
	"github.com/php-any/origami/data"
)

type StmtCloseMethod struct {
	source driversrc.Stmt
}

func (h *StmtCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Close(); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *StmtCloseMethod) GetName() string            { return "close" }
func (h *StmtCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtCloseMethod) GetIsStatic() bool          { return true }
func (h *StmtCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *StmtCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *StmtCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
