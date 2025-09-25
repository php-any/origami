package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ServerHandleMethod struct {
	server *ServerClass
	name   string
}

func (h *ServerHandleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("路由必须是字符串: %v", err))
	}
	param1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("需要传入闭包"))
	}
	if fv, ok := param1.(*data.FuncValue); ok {
		handle, err := newHandler(fv.Value, ctx)
		if err != nil {
			return nil, data.NewErrorThrow(nil, errors.New("路由处理函数入参不对"))
		}
		router := h.server.Prefix + param0
		var finalHandler http.Handler = handle
		if len(h.server.Middlewares) > 0 {
			finalHandler = applyMiddlewares(finalHandler, h.server.Middlewares)
		}
		h.server.source.Handle(router, finalHandler)
		return nil, nil
	}
	return nil, data.NewErrorThrow(nil, errors.New("第二个参数必须是路由处理函数"))
}

func (h *ServerHandleMethod) GetName() string            { return h.name }
func (h *ServerHandleMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerHandleMethod) GetIsStatic() bool          { return false }
func (h *ServerHandleMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, nil),
		node.NewParameter(nil, "handle", 1, nil, nil),
	}
}
func (h *ServerHandleMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
		node.NewVariable(nil, "handle", 1, nil),
	}
}
func (h *ServerHandleMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
