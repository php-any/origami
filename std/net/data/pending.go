package netdata

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var (
	pendingRoutes         []PendingRoute
	pendingControllers    = make(map[string]pendingController)
	controllerMiddlewares = make(map[string][]string)
)

// ControllerInstantiator 在路由注册阶段实例化控制器；默认直接 new，不经过容器。
var ControllerInstantiator = node.InstantiateController

func AddPendingRoute(r PendingRoute) {
	pendingRoutes = append(pendingRoutes, r)
}

func RegisterDeferredController(name string, cls data.ClassStmt, ctx data.Context) {
	pendingControllers[name] = pendingController{ClassStmt: cls, Ctx: ctx}
}

func AddControllerMiddleware(controllerName, middlewareClassName string) {
	controllerMiddlewares[controllerName] = append(controllerMiddlewares[controllerName], middlewareClassName)
}

func instantiatePendingControllers() (map[string]data.GetValue, data.Control) {
	instances := make(map[string]data.GetValue, len(pendingControllers))
	for name, pc := range pendingControllers {
		inst, acl := ControllerInstantiator(pc.ClassStmt, pc.Ctx)
		if acl != nil {
			return nil, acl
		}
		instances[name] = inst
	}
	return instances, nil
}

// RegisterPendingRoutes 实例化待注册控制器并将待注册路由写入全局路由表。
func RegisterPendingRoutes() data.Control {
	instances, acl := instantiatePendingControllers()
	if acl != nil {
		return acl
	}

	for _, pr := range pendingRoutes {
		middlewares := []MiddlewareInfo{}
		if mws, ok := controllerMiddlewares[pr.ControllerName]; ok {
			for _, className := range mws {
				middlewares = append(middlewares, MiddlewareInfo{ClassName: className})
			}
		}
		for _, className := range pr.Middlewares {
			middlewares = append(middlewares, MiddlewareInfo{ClassName: className})
		}

		receiver := pr.Receiver
		if receiver == nil && pr.StaticReceiver != nil {
			receiver = pr.StaticReceiver
		}
		if receiver == nil {
			receiver = instances[pr.ControllerName]
		}

		AppendHTTPRoute(Route{
			Method:      pr.Method,
			Path:        pr.Path,
			Target:      pr.Target,
			Receiver:    receiver,
			Middlewares: middlewares,
			Operation:   pr.Operation,
			HandlerSpec: pr.HandlerSpec,
		})
	}

	pendingRoutes = nil
	pendingControllers = make(map[string]pendingController)
	controllerMiddlewares = make(map[string][]string)
	return nil
}
