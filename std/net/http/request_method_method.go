package http

import (
	"github.com/php-any/origami/data"
)

type RequestMethodMethod struct {
	source *Request
}

func (h *RequestMethodMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	h.source.Method()
	return data.NewStringValue(""), nil
}

func (h *RequestMethodMethod) GetName() string {
	return "method"
}

func (h *RequestMethodMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *RequestMethodMethod) GetIsStatic() bool {
	return false
}

func (h *RequestMethodMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *RequestMethodMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
