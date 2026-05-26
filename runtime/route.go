package runtime

import "github.com/php-any/origami/data"

// AppendHTTPRoute 向 VM 或 TempVM 注册注解路由（生产模式写入全局 VM，开发模式写入请求级 TempVM）。
func AppendHTTPRoute(vm data.VM, r Route) {
	switch v := vm.(type) {
	case *TempVM:
		v.Cache = append(v.Cache, r)
	case *VM:
		v.mu.Lock()
		v.httpRoutes = append(v.httpRoutes, r)
		v.mu.Unlock()
	}
}

// HTTPRoutes 返回已注册的注解路由列表。
func HTTPRoutes(vm data.VM) []Route {
	switch v := vm.(type) {
	case *TempVM:
		return v.Cache
	case *VM:
		v.mu.RLock()
		defer v.mu.RUnlock()
		return v.httpRoutes
	default:
		return nil
	}
}

// SupportsHTTPRoutes 当前 VM 是否支持注解 HTTP 路由注册。
func SupportsHTTPRoutes(vm data.VM) bool {
	switch vm.(type) {
	case *TempVM, *VM:
		return true
	default:
		return false
	}
}
