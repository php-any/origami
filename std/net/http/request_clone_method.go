package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestCloneMethod struct {
	source *httpsrc.Request
}

func (h *RequestCloneMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.Clone(ctx.GoContext())
	return data.NewClassValue(NewRequestClassFrom(ret0), ctx), nil
}

func (h *RequestCloneMethod) GetName() string            { return "clone" }
func (h *RequestCloneMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestCloneMethod) GetIsStatic() bool          { return false }
var requestCloneMethodGetParams = []data.GetValue{}

func (h *RequestCloneMethod) GetParams() []data.GetValue {
	return requestCloneMethodGetParams
}
var requestCloneMethodGetVariables = []data.Variable{}

func (h *RequestCloneMethod) GetVariables() []data.Variable {
	return requestCloneMethodGetVariables
}
func (h *RequestCloneMethod) GetReturnType() data.Types {
	return data.Class{Name: "Net\\Http\\Request"}
}
