package http

import (
	"github.com/php-any/origami/data"
	httpsrc "net/http"
)

type RequestMultipartReaderMethod struct {
	source *httpsrc.Request
}

func (h *RequestMultipartReaderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ret0, ret1 := h.source.MultipartReader()
	return data.NewArrayValue([]data.Value{data.NewAnyValue(ret0), data.NewAnyValue(ret1)}), nil
}

func (h *RequestMultipartReaderMethod) GetName() string               { return "multipartReader" }
func (h *RequestMultipartReaderMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestMultipartReaderMethod) GetIsStatic() bool             { return false }
func (h *RequestMultipartReaderMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestMultipartReaderMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestMultipartReaderMethod) GetReturnType() data.Types     { return data.NewBaseType("void") }
