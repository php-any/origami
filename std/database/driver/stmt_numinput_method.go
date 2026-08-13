package driver

import (
	driversrc "database/sql/driver"
	"github.com/php-any/origami/data"
)

type StmtNumInputMethod struct {
	source driversrc.Stmt
}

func (h *StmtNumInputMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.NumInput()
	return data.NewIntValue(ret0), nil
}

func (h *StmtNumInputMethod) GetName() string            { return "numInput" }
func (h *StmtNumInputMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtNumInputMethod) GetIsStatic() bool          { return true }
func (h *StmtNumInputMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *StmtNumInputMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *StmtNumInputMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
