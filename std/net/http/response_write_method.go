package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	"io"
	httpsrc "net/http"
)

type ResponseWriteMethod struct {
	source *httpsrc.Response
}

func (h *ResponseWriteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[io.Writer](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	ret0 := h.source.Write(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *ResponseWriteMethod) GetName() string            { return "write" }
func (h *ResponseWriteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriteMethod) GetIsStatic() bool          { return true }
func (h *ResponseWriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *ResponseWriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *ResponseWriteMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
