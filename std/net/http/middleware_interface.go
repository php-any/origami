package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Middleware 定义洋葱模型中间件接口
// PHP 类实现此接口后，可通过 #[Middleware(XxxMiddleware::class)] 注解注册
// 洋葱模型：handler($request, $response, $next)
// - $request: 请求对象
// - $response: 响应对象
// - $next: 调用下一个中间件或控制器的回调
type Middleware interface {
	Handle(request data.Value, response data.Value, next data.Value)
}

// NewMiddlewareInterface 创建洋葱模型中间件接口定义
// PHP 接口签名: interface Middleware { public function handle($request, $response, $next); }
func NewMiddlewareInterface() data.InterfaceStmt {
	methods := []data.Method{
		// handle($request, $response, $next): void
		node.NewInterfaceMethod(
			nil,
			"handle",
			"public",
			[]data.GetValue{
				node.NewParameter(nil, "request", 0, nil, nil),
				node.NewParameter(nil, "response", 1, nil, nil),
				node.NewParameter(nil, "next", 2, nil, nil),
			},
			data.NewBaseType("void"),
		),
	}
	return node.NewInterfaceStatement(nil, "Middleware", nil, methods)
}
