package annotation

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

// PendingRoute 待注册路由
type PendingRoute struct {
	Method         string
	Path           string
	Target         data.Method
	Receiver       data.GetValue
	ControllerName string // 控制器名称，用于关联中间件
}

// pendingRoutes 全局存储待注册路由
var pendingRoutes []PendingRoute

// controllerMiddlewares 全局存储控制器中间件
// key: 控制器名称, value: 中间件类名列表
var controllerMiddlewares = make(map[string][]string)

// AddPendingRoute 添加待注册路由
func AddPendingRoute(r PendingRoute) {
	pendingRoutes = append(pendingRoutes, r)
}

// AddControllerMiddleware 添加控制器中间件
func AddControllerMiddleware(controllerName, middlewareClassName string) {
	controllerMiddlewares[controllerName] = append(controllerMiddlewares[controllerName], middlewareClassName)
}

// RegisterPendingRoutes 注册所有待处理路由
// 合并路由信息和中间件信息后注册
func RegisterPendingRoutes(vm data.VM) {
	for _, pr := range pendingRoutes {
		// 获取该控制器的中间件
		middlewares := []runtime.MiddlewareInfo{}
		if mws, ok :=
			controllerMiddlewares[pr.ControllerName]; ok {
			for _, className := range mws {
				middlewares = append(middlewares, runtime.MiddlewareInfo{ClassName: className})
			}
		}

		fmt.Printf("[DEBUG] RegisterPendingRoutes: %s %s, controller=%s, middlewares=%v\n", pr.Method, pr.Path, pr.ControllerName, middlewares)

		// 注册路由
		runtime.AppendHTTPRoute(vm, runtime.Route{
			Method:      pr.Method,
			Path:        pr.Path,
			Target:      pr.Target,
			Receiver:    pr.Receiver,
			Middlewares: middlewares,
		})
	}
	// 清理
	pendingRoutes = nil
	controllerMiddlewares = make(map[string][]string)
}
