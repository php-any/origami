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
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request := NewRequestClassFrom(r)
			response := NewResponseWriterClassFrom(w)

			reqProxy := data.NewProxyValue(request, ctx)
			resProxy := data.NewProxyValue(response, ctx)

			if _, acl := executeMiddlewareChain(vm, ctx, rt, reqProxy, resProxy); acl != nil {
				panic(acl)
			}
		})

		var final http.Handler = handler
		if len(server.Middlewares) > 0 {
			final = applyMiddlewares(final, server.Middlewares)
		}

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
