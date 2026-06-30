package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	httpstd "github.com/php-any/origami/std/net/http"
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

// registerMethodBindings 读取方法上的参数绑定注解，注册到 http 包的绑定表
func registerMethodBindings(target data.Method, cm *node.ClassMethod) {
	vars := target.GetVariables()
	if len(vars) == 0 {
		return
	}

	bindings := make([]httpstd.ParamBind, len(vars))
	// 默认全部为自动绑定
	for i := range bindings {
		bindings[i] = httpstd.BindAuto
	}

	// 读取方法注解中的参数绑定声明
	for _, ann := range cm.Annotations {
		switch a := ann.Class.(type) {
		case *PathVariableClass:
			for _, name := range a.ParamNames() {
				for i, v := range vars {
					if v.GetName() == name {
						bindings[i] = httpstd.BindPath
					}
				}
			}
		case *RequestParamClass:
			for _, name := range a.ParamNames() {
				for i, v := range vars {
					if v.GetName() == name {
						bindings[i] = httpstd.BindQuery
					}
				}
			}
		case *RequestBodyClass:
			for _, name := range a.ParamNames() {
				for i, v := range vars {
					if v.GetName() == name {
						bindings[i] = httpstd.BindBody
					}
				}
			}
		}
	}

	httpstd.RegisterParamBindings(httpstd.MethodID(target), bindings)
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
