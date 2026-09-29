package httpkernel

import (
	"sync"

	"github.com/php-any/origami/data"
)

// fqnRouteCollection / fqnRoute 是路由表里需要按请求复制的类。
const (
	fqnRouteCollection = "Illuminate\\Routing\\RouteCollection"
	fqnRoute           = "Illuminate\\Routing\\Route"
)

// cloneRoutesForRequest 给请求路由器换一份请求私有的路由表。
//
// 为什么必须复制：Illuminate\Routing\Router::findRoute（Router.php:779）每请求都会把
// **请求路由器**的 container 写到匹配到的 Route 对象上：
//
//	$this->current = $route = $this->routes->match($request);
//	$route->setContainer($this->container);
//
// 而 Route 对象来自共享的 RouteCollection（routerSandboxDeepKeys 只隔离
// current / currentRequest，路由集合按设计共享）。于是并发请求会依次把**同一个 Route
// 对象**的 container 改成各自的请求级 Application 克隆，然后 Route.php:254 读它：
//
//	return $this->container[CallableDispatcher::class]->dispatch($this, $callable);
//
// 读到的是别人（或自己）的克隆，谁读到谁就在那个克隆上跑 Container::resolve；
// resolve 未命中时会写 instances / resolved / with / buildStack（Container.php:927/951/962/1133），
// 于是两个 goroutine 写同一张数组 —— 实测 fatal error: concurrent map writes，
// 键长度 47 = Illuminate\Routing\Contracts\CallableDispatcher，栈是
// IndexExpression.SetValue ← Container::resolve ← Application::resolve ← Route::runCallable。
// 这也解释了塌成「写别的请求的 app 克隆」的现场。
//
// 做法：RouteCollection 的数组槽位里，Route 对象换成请求私有副本；同一个 Route 只复制一次
// （allRoutes / routes / nameList / actionList 之间互相引用同一个对象），
// 数组结构本身也复制一层，避免请求期对路由表的就地写落到共享数组上。
// 复制后 Route::setContainer 写的是本请求的副本，CallableDispatcher 在请求级 app 上解析，
// 容器表的写全部落在请求私有数组里。
func cloneRoutesForRequest(ctx data.Context, router data.Value) {
	rv, ok := router.(*data.ClassValue)
	if !ok || rv == nil {
		return
	}
	raw, ctl := rv.GetProperty("routes")
	if ctl != nil || raw == nil {
		return
	}
	col, ok := raw.(*data.ClassValue)
	if !ok || col == nil || col.GetName() != fqnRouteCollection {
		return
	}
	c := &routeCloner{ctx: ctx, memo: map[any]data.Value{}}
	clone := c.collection(col)
	if clone == nil {
		return
	}
	rv.SetProperty("routes", clone)
}

type routeCloner struct {
	ctx  data.Context
	memo map[any]data.Value // 原 Route 的 ObjectValue -> 请求私有副本
}

// collection 浅拷贝 RouteCollection：属性表复制一层，数组属性递归重写。
func (c *routeCloner) collection(col *data.ClassValue) *data.ClassValue {
	if col.ObjectValue == nil {
		return nil
	}
	clone := &data.ClassValue{
		ObjectValue: data.NewObjectValue(),
		Class:       col.Class,
		Context:     c.ctx,
	}
	// 必须带上 VM：RouteCollection 的 matchAgainstRoutes 在父类
	// AbstractRouteCollection 上，方法查找要沿 extends 链加载父类。
	if c.ctx != nil {
		clone.SetVM(c.ctx.GetVM())
	}
	col.RangeProperties(func(key string, value data.Value) bool {
		clone.SetProperty(key, c.rewrite(value, 1))
		return true
	})
	return clone
}

// rewrite 递归重写路由表里的结构：数组复制一层，Route 换请求私有副本，其它保持共享。
func (c *routeCloner) rewrite(v data.Value, depth int) data.Value {
	if v == nil || depth > 8 {
		return v
	}
	switch t := v.(type) {
	case *data.ArrayValue:
		if t == nil {
			return v
		}
		src := t.List
		list := make([]*data.ZVal, len(src))
		for i, z := range src {
			if z == nil {
				continue
			}
			list[i] = data.CopyZValKeepName(z, c.rewrite(z.Value, depth+1))
		}
		return &data.ArrayValue{
			List:                  list,
			IndirectOverloadClass: t.IndirectOverloadClass,
		}
	case *data.ObjectValue:
		if t == nil {
			return v
		}
		out := data.NewObjectValue()
		t.RangeProperties(func(key string, value data.Value) bool {
			out.SetProperty(key, c.rewrite(value, depth+1))
			return true
		})
		return out
	case *data.ThisValue:
		// 路由表里的 Route 以 ThisValue 包装存放（ThisValue 内嵌 *ClassValue，
		// 类型断言不会落进下面的 ClassValue 分支），必须显式拆包再包回去。
		if t == nil || t.ClassValue == nil {
			return v
		}
		if !c.isRoute(t.ClassValue) {
			return v
		}
		clone, ok := c.route(t.ClassValue).(*data.ClassValue)
		if !ok || clone == t.ClassValue {
			return v
		}
		return data.NewThisValue(clone)
	case *data.ClassValue:
		if t == nil {
			return v
		}
		if !c.isRoute(t) {
			return v
		}
		return c.route(t)
	}
	return v
}

// route 返回 Route 的请求私有副本；同一个 Route 只复制一次。
//
// 用 CloneRequestScoped（**链式属性表**）而不是 CloneSandboxKeys：请求期真正会被写到的
// 键只有三个 —— container（Router::findRoute 里 $route->setContainer($this->container)）、
// parameters 与 originalParameters（bind() 整体赋值 / setParameter() 就地写元素）。
// 其余键（uri、methods、action、defaults、wheres、router、$compiled…）都是启动期常量，
// 读不到时回落到全局那份 Route 即可（$compiled 也因此与预热结果共享，见 warmRouteCompiles）。
//
// 旧写法 CloneSandboxKeys(ctx, {"parameters","originalParameters"}) 是「请求开始就把整张
// 属性表浅拷一份、再把两个键深拷」——13 条路由每请求一次，pprof 里
// routeCloner.rewrite + OrderedMap.Range/Set 合计约 1%~5% 的 CPU 都在这。
// 链式表把这份固定税换成「按需升级」：首次写到某个键时才把它按值升级（数组递归深拷贝，
// 见 chainStore / promoteContainer），语义与原来的 deepKeys 列表等价。
func (c *routeCloner) route(cv *data.ClassValue) data.Value {
	if cv.ObjectValue == nil {
		return cv
	}
	if got, ok := c.memo[cv.ObjectValue]; ok {
		return got
	}
	clone := cv.CloneRequestScoped(c.ctx)
	if clone == nil {
		return cv
	}
	c.memo[cv.ObjectValue] = clone
	return clone
}

// isRoute 判断是否是 Illuminate\Routing\Route（含子类）。
// 名字命中缓存，避免每个 Route 都走一次父类链。
func (c *routeCloner) isRoute(cv *data.ClassValue) bool {
	cls := cv.Class
	if cls == nil {
		return false
	}
	name := cls.GetName()
	if hit, ok := routeClassCache.Load(name); ok {
		return hit.(bool)
	}
	hit := false
	for cur := cls; cur != nil; {
		if cur.GetName() == fqnRoute {
			hit = true
			break
		}
		ext := cur.GetExtend()
		if ext == nil || *ext == "" {
			break
		}
		next, ctl := c.loadClass(*ext)
		if ctl != nil || next == nil {
			break
		}
		cur = next
	}
	routeClassCache.Store(name, hit)
	return hit
}

func (c *routeCloner) loadClass(name string) (data.ClassStmt, data.Control) {
	if c.ctx == nil {
		return nil, nil
	}
	vm := c.ctx.GetVM()
	if vm == nil {
		return nil, nil
	}
	return vm.GetOrLoadClass(name)
}

// routeClassCache：类名 -> 是否 Route 子类。进程级只读缓存，请求间稳定。
var routeClassCache sync.Map

// warmRouteCompiles 启动期把每条路由编译一次。
//
// Route::compileRoute（Route.php:373）把结果缓存在 **Route 对象自身** 的 $compiled 上：
//
//	if (! $this->compiled) { $this->compiled = $this->toSymfonyRoute()->compile(); }
//
// 它只由 uri / wheres / domain / methods 决定，全是启动期常量；但请求期用的是
// cloneRoutesForRequest 复制出来的 Route 副本（属性按引用拷一份，$compiled 初值为 null），
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
	for _, z := range arr.List {
		if z == nil || z.Value == nil {
			continue
		}
		// 只编译 Route 本体，别的实现类（或未加载的类）交给请求期原路径。
		if !isRouteValue(ctx, z.Value) {
			continue
		}
		_, _ = callObjectMethodInContext(ctx, z.Value, "compileRoute")
	}
}

// isRouteValue 判断数组槽里的值是不是 Route 子类实例（含 ThisValue 包装）。
func isRouteValue(ctx data.Context, v data.Value) bool {
	var cv *data.ClassValue
	switch t := v.(type) {
	case *data.ThisValue:
		cv = t.ClassValue
	case *data.ClassValue:
		cv = t
	}
	if cv == nil || cv.Class == nil {
		return false
	}
	name := cv.Class.GetName()
	if hit, ok := routeClassCache.Load(name); ok {
		return hit.(bool)
	}
	vm := ctx.GetVM()
	if vm == nil {
		return false
	}
	hit := false
	for cur := cv.Class; cur != nil; {
		if cur.GetName() == fqnRoute {
			hit = true
			break
		}
		ext := cur.GetExtend()
		if ext == nil || *ext == "" {
			break
		}
		next, ctl := vm.GetOrLoadClass(*ext)
		if ctl != nil || next == nil {
			break
		}
		cur = next
	}
	routeClassCache.Store(name, hit)
	return hit
}
