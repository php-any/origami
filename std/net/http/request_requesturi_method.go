package http

import (
	"github.com/php-any/origami/data"
)

type RequestRequestURIMethod struct {
	source *Request
}

func (h *RequestRequestURIMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	h.source.RequestURI()
	return data.NewStringValue(""), nil
}

func (h *RequestRequestURIMethod) GetName() string {
	return "requesturi"
}

func (h *RequestRequestURIMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *RequestRequestURIMethod) GetIsStatic() bool {
	return false
}

func (h *RequestRequestURIMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *RequestRequestURIMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
