package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewMiddlewareInterface 注册 Net\Http\Middleware 洋葱模型中间件接口。
// PHP: interface Middleware { public function handle($request, $response, $next); }
func NewMiddlewareInterface() data.InterfaceStmt {
	methods := []data.Method{
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
	return node.NewInterfaceStatement(nil, "Net\\Http\\Middleware", nil, methods)
}
