package http

import (
	"github.com/php-any/origami/data"
)

type RequestPathMethod struct {
	source *Request
}

func (h *RequestPathMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	h.source.Path()
	return data.NewStringValue(""), nil
}

func (h *RequestPathMethod) GetName() string {
	return "path"
}

func (h *RequestPathMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *RequestPathMethod) GetIsStatic() bool {
	return false
}

func (h *RequestPathMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *RequestPathMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
