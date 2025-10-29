package http

import (
	"errors"
	"net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	runtimesrc "github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

func newHandler(v data.FuncStmt, ctx data.Context) (Handler, error) {
	if len(v.GetVariables()) < 2 {
		return Handler{}, errors.New("invalid variable definition")
	}
	return Handler{Value: v, Ctx: ctx.CreateContext(v.GetVariables())}, nil
}

// Middleware 定义：接收下一个 http.Handler，返回包装后的 http.Handler
type Middleware func(http.Handler) http.Handler

func newMiddleware(v data.FuncStmt, ctx data.Context) (Middleware, error) {
	if len(v.GetVariables()) < 3 {
		return nil, errors.New("invalid variable definition")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mctx := ctx.CreateContext(v.GetVariables())

			request := NewRequestClassFrom(r)
			response := NewResponseWriterClassFrom(w)
			nextHandler := data.NewFuncValue(NextHandler{next: next})

			mctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, mctx))
			mctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, mctx))
			mctx.SetVariableValue(data.NewVariable("next", 2, nil), nextHandler)

			// 检查中间件是否抛出异常
			_, acl := v.Call(mctx)
			if acl != nil {
				ctx.GetVM().ThrowControl(acl)
			}
		})
	}, nil
}

func applyMiddlewares(final http.Handler, middlewares []Middleware) http.Handler {
	h := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

type Handler struct {
	Value data.FuncStmt
	Ctx   data.Context
}

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := f.Ctx.CreateContext(f.Value.GetVariables())

	request := NewRequestClassFrom(r)
	response := NewResponseWriterClassFrom(w)

	ctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, ctx))
	ctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, ctx))

	_, acl := f.Value.Call(ctx)
	if acl != nil {
		// 使用panic抛出异常到NextHandler
		panic(acl)
	}
}

// HotHandler 专用于启用 HotReload 的场景：负责请求期清理
type HotHandler struct {
	Value data.FuncStmt
	Ctx   data.Context
}

func (f HotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 先创建本次请求的函数上下文，再将该 ctx 的 VM 替换为 TempVM
	ctx := f.Ctx.CreateContext(f.Value.GetVariables())
	ctx.SetVM(runtimesrc.NewTempVM(f.Ctx.GetVM()))

	request := NewRequestClassFrom(r)
	response := NewResponseWriterClassFrom(w)

	ctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, ctx))
	ctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, ctx))

	_, acl := f.Value.Call(ctx)
	if acl != nil {
		panic(acl)
	}
}

type NextHandler struct {
	Ctx  data.Context
	next http.Handler
}

func (f NextHandler) Call(ctx data.Context) (_ data.GetValue, acl data.Control) {
	request, err := utils.ConvertFromIndex[*http.Request](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	response, err := utils.ConvertFromIndex[http.ResponseWriter](ctx, 1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	// 使用defer和recover来捕获panic
	defer func() {
		if r := recover(); r != nil {
			if acl2, ok2 := r.(data.Control); ok2 {
				acl = acl2
				return
			}
		}
	}()

	f.next.ServeHTTP(response, request)

	return nil, acl
}

func (f NextHandler) GetName() string {
	return "next"
}

func (f NextHandler) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
	}
}

func (f NextHandler) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
	}
}
