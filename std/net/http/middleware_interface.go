package http

import (
	"github.com/php-any/origami/data"
)

// MiddlewareInterceptor 定义中间件生命周期接口
// 类似 Java Spring 的 HandlerInterceptor
type MiddlewareInterceptor interface {
	// PreHandle 前置处理 - 在控制器方法执行前调用
	// 返回 false 可以中断请求
	PreHandle(request data.Value, response data.Value) bool

	// PostHandle 后置处理 - 在控制器方法执行后调用
	PostHandle(request data.Value, response data.Value)

	// AfterCompletion 完成处理 - 整个请求完成后调用（包括错误处理）
	AfterCompletion(request data.Value, response data.Value)
}
