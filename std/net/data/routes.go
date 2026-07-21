package netdata

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

var httpRoutes []Route

// AppendHTTPRoute 注册一条注解路由到全局路由表。
func AppendHTTPRoute(r Route) {
	httpRoutes = append(httpRoutes, r)
}

// HTTPRoutes 返回已注册的注解路由列表。
func HTTPRoutes() []Route {
	return httpRoutes
}

// ClearHTTPRoutes 清空注解路由表，供开发模式热重载使用。
func ClearHTTPRoutes() {
	httpRoutes = nil
}

// SupportsHTTPRoutes 当前 VM 是否处于可注册注解路由的运行时上下文（全局 VM 或 TempVM）。
func SupportsHTTPRoutes(vm data.VM) bool {
	switch vm.(type) {
	case *runtime.VM, *runtime.TempVM:
		return true
	default:
		return false
	}
}
