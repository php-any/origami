package context

import (
	contextsrc "context"
	"github.com/php-any/origami/data"
)

type ContextDeadlineMethod struct {
	source contextsrc.Context
}

func (h *ContextDeadlineMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0, ret1 := h.source.Deadline()
	return data.NewAnyValue([]any{ret0, ret1}), nil
}

func (h *ContextDeadlineMethod) GetName() string            { return "deadline" }
func (h *ContextDeadlineMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ContextDeadlineMethod) GetIsStatic() bool          { return true }
func (h *ContextDeadlineMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ContextDeadlineMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *ContextDeadlineMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
