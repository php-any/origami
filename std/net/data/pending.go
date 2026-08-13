package netdata

import (
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var (
	pendingMu              sync.Mutex
	pendingRoutes          []PendingRoute
	pendingControllers     = make(map[string]pendingController)
	controllerMiddlewares  = make(map[string][]string)
	controllerInstantiator = node.InstantiateController
)

// ControllerInstantiator 在路由注册阶段实例化控制器；默认直接 new，不经过容器。
func SetControllerInstantiator(instantiator func(data.ClassStmt, data.Context) (data.GetValue, data.Control)) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	if instantiator == nil {
		controllerInstantiator = node.InstantiateController
		return
	}
	controllerInstantiator = instantiator
}

func AddPendingRoute(r PendingRoute) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	pendingRoutes = append(pendingRoutes, r)
}

func RegisterDeferredController(name string, cls data.ClassStmt, ctx data.Context) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	pendingControllers[name] = pendingController{ClassStmt: cls, Ctx: ctx}
}

func AddControllerMiddleware(controllerName, middlewareClassName string) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	controllerMiddlewares[controllerName] = append(controllerMiddlewares[controllerName], middlewareClassName)
}

// RegisterPendingRoutes 固化待注册路由；普通控制器由 ReceiverFactory 按请求实例化。
func RegisterPendingRoutes() data.Control {
	pendingMu.Lock()
	routes := append([]PendingRoute(nil), pendingRoutes...)
	controllers := pendingControllers
	middlewareByController := controllerMiddlewares
	instantiator := controllerInstantiator
	pendingRoutes = nil
	pendingControllers = make(map[string]pendingController)
	controllerMiddlewares = make(map[string][]string)
	pendingMu.Unlock()

	for _, pr := range routes {
		middlewares := []MiddlewareInfo{}
		if mws, ok := middlewareByController[pr.ControllerName]; ok {
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
		var receiverFactory func(data.Context) (data.GetValue, data.Control)
		if receiver == nil {
			if controller, ok := controllers[pr.ControllerName]; ok {
				pc := controller
				receiverFactory = func(requestCtx data.Context) (data.GetValue, data.Control) {
					if requestCtx == nil {
						requestCtx = pc.Ctx
					}
					return instantiator(pc.ClassStmt, requestCtx)
				}
			}
		}

		AppendHTTPRoute(Route{
			Method:          pr.Method,
			Path:            pr.Path,
			Target:          pr.Target,
			Receiver:        receiver,
			ReceiverFactory: receiverFactory,
			Middlewares:     middlewares,
			Operation:       pr.Operation,
			HandlerSpec:     pr.HandlerSpec,
		})
	}

	return nil
}
