package http

import (
	"errors"
	"net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ServerMiddlewareMethod struct {
	server *ServerClass
}

func (h *ServerMiddlewareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("需要传入中间件闭包或类实例"))
	}

	priority := 0
	if _, has := ctx.GetIndexValue(1); has {
		var err error
		priority, err = utils.ConvertFromIndex[int](ctx, 1)
		if err != nil {
			return nil, utils.NewThrowf("priority 参数转换失败: %v", err)
		}
	}

	if fv, ok := v.(*data.FuncValue); ok {
		mw, err := newMiddleware(fv.Value, ctx)
		if err != nil {
			return nil, utils.NewThrow(err)
		}
		h.server.middlewares = append(h.server.middlewares, middlewareEntry{priority: priority, fn: mw})
		return nil, nil
	}

	if cv, ok := v.(*data.ClassValue); ok {
		method, has := cv.GetMethod("handle")
		if !has {
			return nil, utils.NewThrow(errors.New("中间件类必须实现 handle($request, $response, $next) 方法"))
		}

		mw := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rw, response := beginResponse(w)
				defer rw.commitPending()
				r, request := beginRequest(r)
				defer detachRequestAttrs(r)

				mctx := cv.CreateContext(method.GetVariables())
				nextHandler := data.NewFuncValue(ServerMiddlewareNext{
					next: next,
					vars: method.GetVariables(),
				})

				mctx.SetVariableValue(method.GetVariables()[0], data.NewProxyValue(request, mctx))
				mctx.SetVariableValue(method.GetVariables()[1], data.NewProxyValue(response, mctx))
				mctx.SetVariableValue(method.GetVariables()[2], nextHandler)

				_, acl := method.Call(mctx)
				if acl != nil {
					panic(acl)
				}
			})
		}
		h.server.middlewares = append(h.server.middlewares, middlewareEntry{priority: priority, fn: mw})
		return nil, nil
	}

	return nil, utils.NewThrow(errors.New("参数必须是中间件闭包或实现 handle 方法的类实例"))
}

// ServerMiddlewareNext 用于 $server->middleware() 注册的类中间件中 $next 回调
type ServerMiddlewareNext struct {
	next http.Handler
	vars []data.Variable
}

func (f ServerMiddlewareNext) Call(ctx data.Context) (_ data.GetValue, acl data.Control) {
	request, err := utils.ConvertFromIndex[*http.Request](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	response, err := utils.ConvertFromIndex[http.ResponseWriter](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	defer func() {
		if r := recover(); r != nil {
			if acl2, ok2 := r.(data.Control); ok2 {
				acl = acl2
				return
			}
			panic(r)
		}
	}()

	f.next.ServeHTTP(response, request)
	return nil, acl
}

func (f ServerMiddlewareNext) GetName() string { return "next" }
func (f ServerMiddlewareNext) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
	}
}
func (f ServerMiddlewareNext) GetVariables() []data.Variable {
	if len(f.vars) >= 3 {
		return []data.Variable{
			f.vars[0],
			f.vars[1],
		}
	}
	return []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
	}
}

func (h *ServerMiddlewareMethod) GetName() string            { return "middleware" }
func (h *ServerMiddlewareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerMiddlewareMethod) GetIsStatic() bool          { return false }
func (h *ServerMiddlewareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "mid", 0, nil, nil),
		node.NewParameter(nil, "priority", 1, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (h *ServerMiddlewareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "mid", 0, nil),
		node.NewVariable(nil, "priority", 1, nil),
	}
}
func (h *ServerMiddlewareMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
