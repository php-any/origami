package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

// PendingRoute 待注册路由
type PendingRoute struct {
	Method         string
	Path           string
	Target         data.Method
	Receiver       data.GetValue
	ControllerName string
	StaticReceiver data.GetValue
}

type pendingController struct {
	ClassStmt data.ClassStmt
	Ctx       data.Context
}

var pendingRoutes []PendingRoute

// 扫描阶段收集控制器，全部文件加载后再实例化（保证 DI 依赖已注册）
var pendingControllers = make(map[string]pendingController)

var controllerMiddlewares = make(map[string][]string)

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

func RegisterPendingRoutes(vm data.VM) data.Control {
	instances, acl := instantiatePendingControllers()
	if acl != nil {
		return acl
	}

	for _, pr := range pendingRoutes {
		middlewares := []runtime.MiddlewareInfo{}
		if mws, ok := controllerMiddlewares[pr.ControllerName]; ok {
			for _, className := range mws {
				middlewares = append(middlewares, runtime.MiddlewareInfo{ClassName: className})
			}
		}

		receiver := pr.Receiver
		if receiver == nil && pr.StaticReceiver != nil {
			receiver = pr.StaticReceiver
		}
		if receiver == nil {
			receiver = instances[pr.ControllerName]
		}

		runtime.AppendHTTPRoute(vm, runtime.Route{
			Method:      pr.Method,
			Path:        pr.Path,
			Target:      pr.Target,
			Receiver:    receiver,
			Middlewares: middlewares,
		})
	}

	pendingRoutes = nil
	pendingControllers = make(map[string]pendingController)
	controllerMiddlewares = make(map[string][]string)
	return nil
}
