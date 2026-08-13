package http

import (
	"net/http"
	"strings"

	"github.com/php-any/origami/data"
	netdata "github.com/php-any/origami/std/net/data"
	"github.com/php-any/origami/utils"
)

func mountAnnotationRoutes(server *ServerClass, vm data.VM, ctx data.Context, label string) (data.GetValue, data.Control) {
	routes := netdata.HTTPRoutes()
	if len(routes) == 0 {
		return nil, utils.NewThrowf("%s: 未发现注解路由，请确认控制器带有 @Controller/@*Mapping 注解", label)
	}

	for _, rt := range routes {
		rt := rt
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw, response := beginResponse(w, r)
			defer rw.commitPending()
			r, request := beginRequest(r)
			defer detachRequestAttrs(r)

			reqProxy := data.NewProxyValue(request, ctx)
			resProxy := data.NewProxyValue(response, ctx)

			ret, acl := executeMiddlewareChain(vm, ctx, rt, reqProxy, resProxy)
			if acl != nil {
				panic(acl)
			}
			if acl := writeHandlerReturnValue(resProxy, ret); acl != nil {
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
		if rt.Operation != nil {
			op := data.NewObjectValue()
			op.SetProperty("summary", data.NewStringValue(rt.Operation.Summary))
			op.SetProperty("description", data.NewStringValue(rt.Operation.Description))
			op.SetProperty("operationId", data.NewStringValue(rt.Operation.OperationID))
			op.SetProperty("deprecated", data.NewBoolValue(rt.Operation.Deprecated))
			op.SetProperty("hidden", data.NewBoolValue(rt.Operation.Hidden))
			tagValues := make([]data.Value, 0, len(rt.Operation.Tags))
			for _, tag := range rt.Operation.Tags {
				tagValues = append(tagValues, data.NewStringValue(tag))
			}
			op.SetProperty("tags", data.NewArrayValue(tagValues))
			obj.SetProperty("operation", op)
		}
		routeList = append(routeList, obj)
	}

	return data.NewArrayValue(routeList), nil
}
