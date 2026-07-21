package netdata

import "github.com/php-any/origami/data"

// OperationInfo 接口操作元数据（从 @Operation 注解收集）。
type OperationInfo struct {
	Summary     string
	Description string
	Tags        []string
	OperationID string
	Deprecated  bool
	Hidden      bool
}

// Route 已注册的 HTTP 注解路由。
type Route struct {
	Method      string
	Path        string
	Target      data.Method
	Receiver    data.GetValue // 注册路由时已实例化的控制器（或静态方法的 ClassValue）；非空时在其上调用 Target
	Middlewares []MiddlewareInfo
	Operation   *OperationInfo
	HandlerSpec HandlerSpec
}

// MiddlewareInfo 中间件信息（从 @Middleware 注解收集）。
type MiddlewareInfo struct {
	ClassName string
}

// PendingRoute 扫描阶段待注册的路由。
type PendingRoute struct {
	Method         string
	Path           string
	Target         data.Method
	Receiver       data.GetValue
	ControllerName string
	StaticReceiver data.GetValue
	Operation      *OperationInfo
	HandlerSpec    HandlerSpec
	Middlewares    []string // 路由级中间件（Route::middleware / group）
}

type pendingController struct {
	ClassStmt data.ClassStmt
	Ctx       data.Context
}
