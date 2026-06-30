package http

import (
	"net/http"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

func mountAnnotationRoutes(server *ServerClass, vm data.VM, ctx data.Context, label string) (data.GetValue, data.Control) {
	routes := runtime.HTTPRoutes(vm)
	if len(routes) == 0 {
		return nil, utils.NewThrowf("%s: 未发现注解路由，请确认控制器带有 @Controller/@*Mapping 注解", label)
	}

	for _, rt := range routes {
		rt := rt
		// 预提取路由路径中的参数名列表
		pathParamKeys := extractPathParamKeys(rt.Path)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw, response := beginResponse(w, r)
			defer rw.commitPending()
			r, request := beginRequest(r)
			defer detachRequestAttrs(r)

			// 将路径参数键名关联到请求，以便 all/input/only/except/has/post/route 等方法获取
			setPathValueKeys(r, pathParamKeys)

			reqProxy := data.NewProxyValue(request, ctx)
			resProxy := data.NewProxyValue(response, ctx)

			if _, acl := executeMiddlewareChain(vm, ctx, rt, reqProxy, resProxy); acl != nil {
				panic(acl)
			}
		})

		final := server.finalizeHandler(handler)

		methodPath := strings.ToUpper(rt.Method) + " " + rt.Path
		server.source.Handle(methodPath, final)
	}

	routeList := make([]data.Value, 0, len(routes))
	for _, rt := range routes {
		obj := data.NewObjectValue()
		obj.SetProperty("method", data.NewStringValue(rt.Method))
		obj.SetProperty("path", data.NewStringValue(rt.Path))
		routeList = append(routeList, obj)
	}

	return data.NewArrayValue(routeList), nil
}
