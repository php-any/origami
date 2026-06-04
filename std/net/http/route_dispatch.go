package http

import (
	"fmt"
	"net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

// DispatchHTTPRoutes 根据 VM 中已注册的注解路由分发当前请求。
func DispatchHTTPRoutes(vm data.VM, ctx data.Context) (data.GetValue, data.Control) {
	routes := runtime.HTTPRoutes(vm)
	if len(routes) == 0 {
		return nil, utils.NewThrowf("未注册 HTTP 路由，请确认应用入口已加载且控制器带有 @Controller/@*Mapping 注解")
	}

	mux := http.NewServeMux()
	var lastACL data.Control
	for _, rt := range routes {
		rt := rt
		mux.HandleFunc(rt.Method+" "+rt.Path, func(w http.ResponseWriter, r *http.Request) {
			request := NewRequestClassFrom(r)
			response := NewResponseWriterClassFrom(w)

			reqProxy := data.NewProxyValue(request, ctx)
			resProxy := data.NewProxyValue(response, ctx)

			// 执行中间件 preHandle
			if !executePreHandle(vm, ctx, rt.Middlewares, reqProxy, resProxy) {
				return // preHandle 返回 false，中断请求
			}

			// 执行控制器方法
			executeControllerMethod(rt, reqProxy, resProxy, ctx, &lastACL)

			// 执行中间件 postHandle
			executePostHandle(vm, ctx, rt.Middlewares, reqProxy, resProxy)

			// 执行中间件 afterCompletion
			executeAfterCompletion(vm, ctx, rt.Middlewares, reqProxy, resProxy)
		})
	}

	req, err := utils.ConvertFromIndex[*http.Request](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	res, err := utils.ConvertFromIndex[http.ResponseWriter](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	mux.ServeHTTP(res, req)
	return nil, lastACL
}

// executePreHandle 执行所有中间件的 preHandle 方法
// 返回 false 表示请求被中断
func executePreHandle(vm data.VM, ctx data.Context, middlewares []runtime.MiddlewareInfo, request data.Value, response data.Value) bool {
	if len(middlewares) == 0 {
		return true
	}

	fmt.Printf("[DEBUG] executePreHandle: middlewares=%v\n", middlewares)

	// 按书写顺序依次执行 preHandle
	for _, mw := range middlewares {
		cls, ok := vm.GetClass(mw.ClassName)
		if !ok {
			fmt.Printf("[DEBUG] executePreHandle: class %s not found\n", mw.ClassName)
			continue
		}
		fmt.Printf("[DEBUG] executePreHandle: calling preHandle on %s\n", mw.ClassName)

		// 实例化中间件类
		inst, acl := cls.GetValue(ctx)
		if acl != nil {
			continue
		}

		cv, ok := inst.(*data.ClassValue)
		if !ok {
			continue
		}

		// 调用 preHandle 方法
		method, has := cls.GetMethod("preHandle")
		if !has {
			continue
		}

		fnCtx := cv.CreateContext(method.GetVariables())
		if len(method.GetVariables()) >= 1 {
			fnCtx.SetVariableValue(method.GetVariables()[0], request)
		}
		if len(method.GetVariables()) >= 2 {
			fnCtx.SetVariableValue(method.GetVariables()[1], response)
		}

		result, acl := method.Call(fnCtx)
		if acl != nil {
			return false
		}

		// 检查返回值，false 表示中断
		if result != nil {
			if boolVal, ok := result.(*data.BoolValue); ok {
				if !boolVal.Value {
					return false
				}
			}
		}
	}

	return true
}

// executePostHandle 执行所有中间件的 postHandle 方法（逆序）
func executePostHandle(vm data.VM, ctx data.Context, middlewares []runtime.MiddlewareInfo, request data.Value, response data.Value) {
	if len(middlewares) == 0 {
		return
	}

	// 逆序执行 postHandle
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		cls, ok := vm.GetClass(mw.ClassName)
		if !ok {
			continue
		}

		// 实例化中间件类
		inst, acl := cls.GetValue(ctx)
		if acl != nil {
			continue
		}

		cv, ok := inst.(*data.ClassValue)
		if !ok {
			continue
		}

		// 调用 postHandle 方法
		method, has := cls.GetMethod("postHandle")
		if !has {
			continue
		}

		fnCtx := cv.CreateContext(method.GetVariables())
		if len(method.GetVariables()) >= 1 {
			fnCtx.SetVariableValue(method.GetVariables()[0], request)
		}
		if len(method.GetVariables()) >= 2 {
			fnCtx.SetVariableValue(method.GetVariables()[1], response)
		}

		method.Call(fnCtx)
	}
}

// executeAfterCompletion 执行所有中间件的 afterCompletion 方法（逆序）
func executeAfterCompletion(vm data.VM, ctx data.Context, middlewares []runtime.MiddlewareInfo, request data.Value, response data.Value) {
	if len(middlewares) == 0 {
		return
	}

	// 逆序执行 afterCompletion
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		cls, ok := vm.GetClass(mw.ClassName)
		if !ok {
			continue
		}

		// 实例化中间件类
		inst, acl := cls.GetValue(ctx)
		if acl != nil {
			continue
		}

		cv, ok := inst.(*data.ClassValue)
		if !ok {
			continue
		}

		// 调用 afterCompletion 方法
		method, has := cls.GetMethod("afterCompletion")
		if !has {
			continue
		}

		fnCtx := cv.CreateContext(method.GetVariables())
		if len(method.GetVariables()) >= 1 {
			fnCtx.SetVariableValue(method.GetVariables()[0], request)
		}
		if len(method.GetVariables()) >= 2 {
			fnCtx.SetVariableValue(method.GetVariables()[1], response)
		}

		method.Call(fnCtx)
	}
}

// executeControllerMethod 执行控制器方法
func executeControllerMethod(rt runtime.Route, reqProxy data.Value, resProxy data.Value, ctx data.Context, lastACL *data.Control) {
	args := []data.Value{reqProxy, resProxy}
	if len(rt.Target.GetVariables()) < len(args) {
		args = args[:len(rt.Target.GetVariables())]
	}

	if rt.Receiver != nil {
		_, *lastACL = node.CallHTTPControllerMethod(rt.Receiver, rt.Target, args)
		return
	}

	mute := ctx.CreateContext(rt.Target.GetVariables())
	for i, arg := range args {
		if i < len(rt.Target.GetVariables()) {
			mute.SetVariableValue(rt.Target.GetVariables()[i], arg)
		}
	}
	_, *lastACL = rt.Target.Call(mute)
}
