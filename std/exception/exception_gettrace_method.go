package exception

import (
	"github.com/php-any/origami/data"
)

type ExceptionGetTraceMethod struct {
	source *Exception
}

func (h *ExceptionGetTraceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewArrayValue(h.source.GetTraceValues()), nil
}

func (h *ExceptionGetTraceMethod) GetName() string            { return "getTrace" }
func (h *ExceptionGetTraceMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ExceptionGetTraceMethod) GetIsStatic() bool          { return false }
func (h *ExceptionGetTraceMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (h *ExceptionGetTraceMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *ExceptionGetTraceMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}
