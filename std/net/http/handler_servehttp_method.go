package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	httpsrc "net/http"
)

type HandlerServeHTTPMethod struct {
	source httpsrc.Handler
}

func (h *HandlerServeHTTPMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[httpsrc.ResponseWriter](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}
	param1, err := utils.ConvertFromIndex[*httpsrc.Request](ctx, 1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	h.source.ServeHTTP(param0, param1)
	return nil, nil
}

func (h *HandlerServeHTTPMethod) GetName() string            { return "serveHTTP" }
func (h *HandlerServeHTTPMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *HandlerServeHTTPMethod) GetIsStatic() bool          { return false }
func (h *HandlerServeHTTPMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *HandlerServeHTTPMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *HandlerServeHTTPMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
