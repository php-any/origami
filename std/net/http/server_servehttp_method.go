package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ServerServeHTTPMethod 暴露底层 ServeMux 的 ServeHTTP，便于不监听端口直接分发
type ServerServeHTTPMethod struct{ server *ServerClass }

func (h *ServerServeHTTPMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	w, err := utils.ConvertFromIndex[httpsrc.ResponseWriter](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}
	r, err := utils.ConvertFromIndex[*httpsrc.Request](ctx, 1)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}
	h.server.source.ServeHTTP(w, r)
	return nil, nil
}

func (h *ServerServeHTTPMethod) GetName() string            { return "serveHTTP" }
func (h *ServerServeHTTPMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerServeHTTPMethod) GetIsStatic() bool          { return false }
func (h *ServerServeHTTPMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *ServerServeHTTPMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *ServerServeHTTPMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
