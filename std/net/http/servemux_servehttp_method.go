package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	httpsrc "net/http"
)

type ServeMuxServeHTTPMethod struct {
	source *httpsrc.ServeMux
}

func (h *ServeMuxServeHTTPMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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

func (h *ServeMuxServeHTTPMethod) GetName() string            { return "serveHTTP" }
func (h *ServeMuxServeHTTPMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServeMuxServeHTTPMethod) GetIsStatic() bool          { return true }
func (h *ServeMuxServeHTTPMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *ServeMuxServeHTTPMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *ServeMuxServeHTTPMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
