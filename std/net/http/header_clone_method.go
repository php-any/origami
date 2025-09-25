package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type HeaderCloneMethod struct {
	source *httpsrc.Header
}

func (h *HeaderCloneMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0 := h.source.Clone()
	retPtr := &ret0
	return data.NewClassValue(NewHeaderClassFrom(retPtr), ctx), nil
}

func (h *HeaderCloneMethod) GetName() string               { return "clone" }
func (h *HeaderCloneMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *HeaderCloneMethod) GetIsStatic() bool             { return false }
func (h *HeaderCloneMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *HeaderCloneMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *HeaderCloneMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
