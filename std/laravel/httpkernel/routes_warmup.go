package httpkernel

import "github.com/php-any/origami/data"

// warmRouteCompiles 启动期把每条路由编译一次。
//
// Route::compileRoute（Route.php:373）把结果缓存在 **Route 对象自身** 的 $compiled 上：
//
//	if (! $this->compiled) { $this->compiled = $this->toSymfonyRoute()->compile(); }
//
// 它只由 uri / wheres / domain / methods 决定，全是启动期常量；但请求期用的是
// 通用对象图策略产生的 Route 请求副本（属性按引用拷一份，$compiled 初值为 null），
// 副本上编译完也只写进副本 —— 全局那份永远是 null，于是**每个候选路由每请求都要重跑**
// RouteCompiler（拼正则、preg_replace、preg_match_all）。
// 实测 /hello：compileRoute 14 次/请求、toSymfonyRoute 13、getOptionalParameterNames 13、
// Route::getDomain 14、Route::uri 26。
//
// 先在全局 Route 上编译一次，请求复制时 $compiled 会按引用带给每个副本
// （CompiledRoute 编译后只读，跨请求共享安全），副本上的 compileRoute 直接命中缓存短路。
func warmRouteCompiles(ctx data.Context, router data.Value) {
	if ctx == nil || router == nil {
		return
	}
	rv, ok := router.(*data.ClassValue)
	if !ok || rv == nil {
		return
	}
	raw, ctl := rv.GetProperty("routes")
	if ctl != nil || raw == nil {
		return
	}
	col, ok := raw.(*data.ClassValue)
	if !ok || col == nil {
		return
	}
	list, ctl := callObjectMethodInContext(ctx, col, "getRoutes")
	if ctl != nil || list == nil {
		return
	}
	arr, ok := list.(*data.ArrayValue)
	if !ok || arr == nil {
		return
	}
	for arraySlots77, arrayPosition77 := arr.View(), 0; arrayPosition77 < arraySlots77.Len(); arrayPosition77++ {
		z := arraySlots77.At(arrayPosition77)
		if z == nil || z.ReadValue() == nil {
			continue
		}
		// 只编译 Route 本体，别的实现类（或未加载的类）交给请求期原路径。
		if !isRouteValue(ctx, z.ReadValue()) {
			continue
		}
		_, _ = callObjectMethodInContext(ctx, z.ReadValue(), "compileRoute")
	}
}

func isRouteValue(ctx data.Context, value data.Value) bool {
	var object *data.ClassValue
	switch v := value.(type) {
	case *data.ClassValue:
		object = v
	case *data.ThisValue:
		object = v.ClassValue
	}
	return object != nil && data.NominalIsA(object.Class, "Illuminate\\Routing\\Route", ctx.GetVM())
}
