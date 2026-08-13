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
	Method          string
	Path            string
	Target          data.Method
	Receiver        data.GetValue                                    // 静态方法或显式注册对象的接收者
	ReceiverFactory func(data.Context) (data.GetValue, data.Control) // 普通控制器按请求创建，避免跨请求共享对象状态
	Middlewares     []MiddlewareInfo
	Operation       *OperationInfo
	HandlerSpec     HandlerSpec
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
