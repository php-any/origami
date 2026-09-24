package context

import (
	contextsrc "context"
	"github.com/php-any/origami/data"
)

type ContextDoneMethod struct {
	source contextsrc.Context
}

func (h *ContextDoneMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.Done()
	return data.NewAnyValue(ret0), nil
}

func (h *ContextDoneMethod) GetName() string            { return "done" }
func (h *ContextDoneMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ContextDoneMethod) GetIsStatic() bool          { return true }
var contextDoneMethodGetParams = []data.GetValue{}

func (h *ContextDoneMethod) GetParams() []data.GetValue {
	return contextDoneMethodGetParams
}

var contextDoneMethodGetVariables = []data.Variable{}

func (h *ContextDoneMethod) GetVariables() []data.Variable {
	return contextDoneMethodGetVariables
}

func (h *ContextDoneMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
