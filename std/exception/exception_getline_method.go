package exception

import (
	"github.com/php-any/origami/data"
)

type ExceptionGetLineMethod struct {
	source *Exception
}

func (h *ExceptionGetLineMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if line, ok := instancePropertyInt(ctx, "line"); ok && line > 0 {
		return data.NewIntValue(line), nil
	}
	return data.NewIntValue(h.source.GetLine()), nil
}

func (h *ExceptionGetLineMethod) GetName() string            { return "getLine" }
func (h *ExceptionGetLineMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ExceptionGetLineMethod) GetIsStatic() bool          { return false }
func (h *ExceptionGetLineMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (h *ExceptionGetLineMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *ExceptionGetLineMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}
