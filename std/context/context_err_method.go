package context

import (
	contextsrc "context"
	"github.com/php-any/origami/data"
)

type ContextErrMethod struct {
	source contextsrc.Context
}

func (h *ContextErrMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Err(); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *ContextErrMethod) GetName() string            { return "err" }
func (h *ContextErrMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ContextErrMethod) GetIsStatic() bool          { return true }
func (h *ContextErrMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ContextErrMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *ContextErrMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
