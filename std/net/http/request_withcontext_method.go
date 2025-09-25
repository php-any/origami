package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestWithContextMethod struct {
	source *httpsrc.Request
}

func (h *RequestWithContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.WithContext(ctx.GoContext())
	return data.NewClassValue(NewRequestClassFrom(ret0), ctx), nil
}

func (h *RequestWithContextMethod) GetName() string            { return "withContext" }
func (h *RequestWithContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestWithContextMethod) GetIsStatic() bool          { return true }
func (h *RequestWithContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (h *RequestWithContextMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *RequestWithContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
