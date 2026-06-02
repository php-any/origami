package http

import (
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
			args := []data.Value{reqProxy, resProxy}
			if len(rt.Target.GetVariables()) < len(args) {
				args = args[:len(rt.Target.GetVariables())]
			}

			if rt.Receiver != nil {
				_, lastACL = node.CallHTTPControllerMethod(rt.Receiver, rt.Target, args)
				return
			}

			mute := ctx.CreateContext(rt.Target.GetVariables())
			for i, arg := range args {
				if i < len(rt.Target.GetVariables()) {
					mute.SetVariableValue(rt.Target.GetVariables()[i], arg)
				}
			}
			_, lastACL = rt.Target.Call(mute)
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
