package http

import (
	"fmt"
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ResponseWriterWriteMethod struct {
	source httpsrc.ResponseWriter
}

func (h *ResponseWriterWriteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("write方法缺少参数: %v", 0))
	}

	ret0, ret1 := h.source.Write([]byte(param0.AsString()))
	if ret1 != nil {
		return nil, data.NewErrorThrow(nil, ret1)
	}
	return data.NewIntValue(ret0), nil
}

func (h *ResponseWriterWriteMethod) GetName() string            { return "write" }
func (h *ResponseWriterWriteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterWriteMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterWriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *ResponseWriterWriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *ResponseWriterWriteMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
