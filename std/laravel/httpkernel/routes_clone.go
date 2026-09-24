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

// routeSandboxDeepKeys：Route 上请求期会被就地写的数组。
// setParameter() 走 $this->parameters[$name] = $value（就地写数组），
// bind() 走 $this->parameters = <新数组>（整体赋值）；两者都要求
// parameters 是请求私有的，否则并发请求会写同一张数组。
var routeSandboxDeepKeys = []string{
	"parameters",
	"originalParameters",
}

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
func (c *routeCloner) route(cv *data.ClassValue) data.Value {
	if cv.ObjectValue == nil {
		return cv
	}
	if got, ok := c.memo[cv.ObjectValue]; ok {
		return got
	}
	clone := cv.CloneSandboxKeys(c.ctx, routeSandboxDeepKeys)
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
