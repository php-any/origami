package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	"io"
	httpsrc "net/http"
)

type RequestWriteProxyMethod struct {
	source *httpsrc.Request
}

func (h *RequestWriteProxyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[io.Writer](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	ret0 := h.source.WriteProxy(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *RequestWriteProxyMethod) GetName() string            { return "writeProxy" }
func (h *RequestWriteProxyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestWriteProxyMethod) GetIsStatic() bool          { return false }
func (h *RequestWriteProxyMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestWriteProxyMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestWriteProxyMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
