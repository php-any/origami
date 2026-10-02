package http

import (
	"context"
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
			requestVM, scoped, owner := languageRequest(ctx, w, r)
			r = scoped
			if owner {
				defer requestVM.RunShutdownCallbacks()
			}
			rw, response := beginResponse(w, r)
			defer rw.commitPending()
			r, request := beginRequest(r)
			defer detachRequestAttrs(r)

			mctx := requestVM.CreateContext(v.GetVariables())
			nextHandler := data.NewFuncValue(NextHandler{next: next})

			mctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, mctx))
			mctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, mctx))
			mctx.SetVariableValue(data.NewVariable("next", 2, nil), nextHandler)

			_, acl := v.Call(mctx)
			if acl != nil {
				finishLanguageControl(requestVM, acl, owner)
			}
		})
	}, nil
}

type Handler struct {
	Value data.FuncStmt
	Ctx   data.Context
}

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestVM, scoped, owner := languageRequest(f.Ctx, w, r)
	r = scoped
	if owner {
		defer requestVM.RunShutdownCallbacks()
	}
	node.ResetSuperglobals()
	rw, response := beginResponse(w, r)
	defer rw.commitPending()
	r, request := beginRequest(r)
	defer detachRequestAttrs(r)

	ctx := requestVM.CreateContext(f.Value.GetVariables())

	ctx.SetVariableValue(data.NewVariable("r", 0, nil), data.NewProxyValue(request, ctx))
	ctx.SetVariableValue(data.NewVariable("w", 1, nil), data.NewProxyValue(response, ctx))

	_, acl := f.Value.Call(ctx)
	if acl != nil {
		finishLanguageControl(requestVM, acl, owner)
	}
}

// HotHandler uses the same language request VM as every HTTP handler.
type HotHandler struct {
	Value data.FuncStmt
	Ctx   data.Context
}

func (f HotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Handler{Value: f.Value, Ctx: f.Ctx}.ServeHTTP(w, r)
}

type languageRequestKey struct{}

func languageRequest(ctx data.Context, response http.ResponseWriter, request *http.Request) (*runtimesrc.RequestVM, *http.Request, bool) {
	if vm, ok := request.Context().Value(languageRequestKey{}).(*runtimesrc.RequestVM); ok {
		return vm, request, false
	}
	vm := runtimesrc.NewRequestVM(ctx.GetVM()).(*runtimesrc.RequestVM)
	vm.BindHTTP(request, response)
	scoped := request.WithContext(context.WithValue(request.Context(), languageRequestKey{}, vm))
	return vm, scoped, true
}

func finishLanguageControl(vm *runtimesrc.RequestVM, control data.Control, owner bool) {
	if owner {
		_, control = vm.HandleUnhandledException(control)
		if control == nil {
			return
		}
		if exit, ok := control.(data.ExitControl); ok && exit.IsExit() {
			return
		}
	}
	panic(control)
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

var nextHandlerGetParams = []data.GetValue{
	node.NewParameter(nil, "request", 0, nil, nil),
	node.NewParameter(nil, "response", 1, nil, nil),
}

func (f NextHandler) GetParams() []data.GetValue {
	return nextHandlerGetParams
}

var nextHandlerGetVariables = []data.Variable{
	node.NewVariable(nil, "request", 0, nil),
	node.NewVariable(nil, "response", 1, nil),
}

func (f NextHandler) GetVariables() []data.Variable {
	return nextHandlerGetVariables
}
