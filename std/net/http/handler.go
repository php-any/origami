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

// MiddlewareFunc 定义：接收下一个 http.Handler，返回包装后的 http.Handler
type MiddlewareFunc func(http.Handler) http.Handler

func newMiddleware(v data.FuncStmt, ctx data.Context) (MiddlewareFunc, error) {
	if len(v.GetVariables()) < 3 {
		return nil, errors.New("invalid variable definition")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw, response := beginResponse(w)
			defer rw.commitPending()
			r, request := beginRequest(r)
			defer detachRequestAttrs(r)

			mctx := ctx.CreateContext(v.GetVariables())
			nextHandler := data.NewFuncValue(NextHandler{next: next})

			mctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, mctx))
			mctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, mctx))
			mctx.SetVariableValue(data.NewVariable("next", 2, nil), nextHandler)

			_, acl := v.Call(mctx)
			if acl != nil {
				ctx.GetVM().ThrowControl(acl)
			}
		})
	}, nil
}

type Handler struct {
	Value data.FuncStmt
	Ctx   data.Context
}

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	node.ResetSuperglobals()
	rw, response := beginResponse(w)
	defer rw.commitPending()
	r, request := beginRequest(r)
	defer detachRequestAttrs(r)

	ctx := f.Ctx.CreateContext(f.Value.GetVariables())

	ctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, ctx))
	ctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, ctx))

	_, acl := f.Value.Call(ctx)
	if acl != nil {
		panic(acl)
	}
}

// HotHandler 专用于启用 HotReload 的场景：负责请求期清理
type HotHandler struct {
	Value data.FuncStmt
	Ctx   data.Context
}

func (f HotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	node.ResetSuperglobals()
	rw, response := beginResponse(w)
	defer rw.commitPending()
	r, request := beginRequest(r)
	defer detachRequestAttrs(r)

	ctx := f.Ctx.CreateContext(f.Value.GetVariables())
	ctx.SetVM(runtimesrc.NewTempVM(f.Ctx.GetVM()))

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
