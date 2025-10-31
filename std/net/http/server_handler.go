package http

import (
	"errors"
	"net/http"
	"strings"

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
		return nil, utils.NewThrowf("路由必须是字符串: %v", err)
	}
	param1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("需要传入闭包"))
	}
	if fv, ok := param1.(*data.FuncValue); ok {
		var final http.Handler
		final, err = newHandler(fv.Value, ctx)
		if err != nil {
			return nil, utils.NewThrow(errors.New("路由处理函数入参不对"))
		}
		router := h.server.Prefix + param0
		if len(h.server.Middlewares) > 0 {
			final = applyMiddlewares(final, h.server.Middlewares)
		}

		// 使用 Go 1.22+ 的方法路由语法
		methodPath := strings.ToUpper(h.name) + " " + router
		h.server.source.Handle(methodPath, final)
		return nil, nil
	}
	return nil, utils.NewThrow(errors.New("第二个参数必须是路由处理函数"))
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

// ServerAnyMethod 注册任意方法与任意路径的处理器：$server->any(($request, $response) => {...})
type ServerAnyMethod struct {
	server *ServerClass
}

func (h *ServerAnyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 只接受一个函数参数
	param0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("需要传入闭包"))
	}
	if fv, ok := param0.(*data.FuncValue); ok {
		var final http.Handler
		final, err := newHandler(fv.Value, ctx)
		if err != nil {
			return nil, utils.NewThrow(errors.New("路由处理函数入参不对"))
		}

		if len(h.server.Middlewares) > 0 {
			final = applyMiddlewares(final, h.server.Middlewares)
		}

		// 使用 Go 1.22+ ServeMux 的方法路由语法：为常见方法注册前缀 "/" 的通配路由
		methods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "TRACE"}
		route := h.server.Prefix + "/"
		for _, m := range methods {
			h.server.source.Handle(m+" "+route, final)
		}
		return nil, nil
	}
	return nil, utils.NewThrow(errors.New("参数必须是路由处理函数"))
}

func (h *ServerAnyMethod) GetName() string            { return "any" }
func (h *ServerAnyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerAnyMethod) GetIsStatic() bool          { return false }
func (h *ServerAnyMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "handle", 0, nil, nil),
	}
}
func (h *ServerAnyMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "handle", 0, nil),
	}
}
func (h *ServerAnyMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
