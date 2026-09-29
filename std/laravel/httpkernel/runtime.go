package httpkernel

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// Resolve 从 Laravel 容器解析 Go 实现的 App\Http\Kernel。
// 该路径会正常触发 ApplicationBuilder 注册的 afterResolving 回调。
func Resolve(app *data.ClassValue) (*data.ClassValue, data.Control) {
	if app == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: Application 不可用"))
	}
	value, control := callObjectMethod(
		app,
		"make",
		data.NewStringValue(fqnKernelContract),
	)
	if control != nil {
		return nil, control
	}
	kernel, ok := value.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 容器未返回 App\\Http\\Kernel"))
	}
	return kernel, nil
}

// Bootstrap 在服务启动阶段完成一次 Laravel Kernel bootstrap。
func Bootstrap(kernel *data.ClassValue) data.Control {
	if kernel == nil {
		return data.NewErrorThrow(nil, fmt.Errorf("httpkernel: Kernel 不可用"))
	}
	_, control := callObjectMethod(kernel, "bootstrap")
	return control
}

// warmAbstracts 是「请求无关」的重服务。这些服务只要在全局 Application 上已经解析过，
// 就会留在全局 instances 表里，请求级 Application 首次查到即共享同一个实例；
// 否则每个请求的第一次 $app[$abstract] 都要在解释器里重跑一遍完整 resolve 链。
// 实测（examples/laravel13，/origami-health）：不预热时 Container::instance() 的
// rebound('request') 回调里 $app['auth'] 每请求重解析，占单请求 CPU 的 50% 上下。
//
// 只列请求无关的：request / session 之类带请求状态的不能预热，仍由沙箱隔离。
var warmAbstracts = []string{
	"config",
	"events",
	"log",
	"db",
	"db.factory",
	"cache",
	"cache.store",
	"session",
	// "session.store" 故意不预热：它是 `fn ($app) => $app['session']->driver()`
	// （SessionServiceProvider::registerSessionDriver），解析它就是**造出一个请求级的 Store**，
	// 而预热发生在全局 app 上，等于把一份会话状态挂到全局（详见 resetSessionDrivers）。
	"cookie",
	"encrypter",
	"hash",
	"auth",
	"translator",
	"view",
	"view.engine.resolver",
	"filesystem",
	"blade.compiler",
	"url",
	"redirect",
	"router",
}

// Warm 把请求无关的重服务在全局 Application 上解析一次，供所有请求沙箱共享副本。
// 解析失败不致命：bound 为假或 resolve 抛错都跳过，请求期仍可按原路径懒解析。
func Warm(ctx data.Context, kernel *data.ClassValue) {
	if kernel == nil || ctx == nil {
		return
	}
	st, _ := kernel.GetSource().(*kernelState)
	if st == nil || st.app == nil {
		return
	}
	for _, abstract := range warmAbstracts {
		bound, ctl := callObjectMethodInContext(ctx, st.app, "bound", data.NewStringValue(abstract))
		if ctl != nil || !isTrueValue(bound) {
			continue
		}
		_, _ = callObjectMethodInContext(ctx, st.app, "make", data.NewStringValue(abstract))
	}
	// 路由的 Symfony 编译结果是纯粹的启动期常量，先算出来（详见 warmRouteCompiles）。
	warmRouteCompiles(ctx, st.router)
	// Telescope 的开关、匹配 pattern、实例同样都是启动期常量。在全局 Application 上解析一次，
	// 请求期就只剩一次 $request->is()；放在这里（而非首个请求里）也避免并发首请求在
	// once.Do 上排队，以及把 Telescope 实例建进某个请求沙箱。
	if st.tel != nil {
		st.tel.once.Do(func() { st.tel.resolve(ctx, st.app) })
	}
}

func isTrueValue(v data.GetValue) bool {
	b, ok := asValue(v).(*data.BoolValue)
	return ok && b.Value
}

// Handle 将原生 Request 交给 Go HTTP Kernel。
func Handle(ctx data.Context, kernel *data.ClassValue, request data.Value) (data.GetValue, data.Control) {
	if kernel == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: Kernel 不可用"))
	}
	return callObjectMethodInContext(ctx, kernel, "handle", request)
}

// Sandbox 为当前请求复制 Kernel / Application / Router，并绑到 goroutine 本地的
// Container::$instance 与 Facade，使并发 Handle 各自处理自己的 Request。
func Sandbox(ctx data.Context, kernel *data.ClassValue) *data.ClassValue {
	if kernel == nil {
		return nil
	}
	src, _ := kernel.GetSource().(*kernelState)
	if src == nil {
		return kernel.CloneSandbox(ctx)
	}
	st := &kernelState{
		app:                cloneApplication(src.app, ctx),
		router:             cloneRouter(src.router, ctx),
		bootstrappers:      append([]string(nil), src.bootstrappers...),
		middleware:         append([]string(nil), src.middleware...),
		middlewareGroups:   cloneGroupsMap(src.middlewareGroups),
		middlewareAliases:  cloneStringMap(src.middlewareAliases),
		middlewarePriority: append([]string(nil), src.middlewarePriority...),
		bootstrapped:       src.bootstrapped,
		// Telescope 判定是启动期常量：沙箱共享常驻 state 的那一份，不每请求重算。
		tel: src.tel,
	}
	cloneRoutesForRequest(ctx, st.router)
	kc := newKernelClass(st)
	cv := data.NewProxyValue(kc, ctx)
	if kernel.ObjectValue != nil {
		cv.ObjectValue = data.DeepCloneObjectValue(kernel.ObjectValue)
	}
	if ctx != nil {
		if vm := ctx.GetVM(); vm != nil {
			cv.SetVM(vm)
		}
	}
	syncProperties(cv, st)
	bindSandboxContainer(ctx, st.app, st.router)
	cloneRequestServices(ctx, st.app)
	return cv
}

// routerSandboxDeepKeys：Router 上随请求变化的少量状态；RouteCollection 启动后只读，共享。
var routerSandboxDeepKeys = []string{
	"current",
	"currentRequest",
}

// cloneApplication 造请求级 Application：它是全局 Application 的一个空壳，
// 未命中的属性读会回落到全局（见 data.CloneRequestScoped / data.chainStore）。
//
// 这里曾经是 CloneSandboxKeys(ctx, appSandboxDeepKeys)——每请求把 22 张容器表全量深拷贝。
// 那份键表本身是对的（这些表请求期确实可能被 Container 的写方法改到）：
//   bind()/instance()/alias()      -> bindings / aliases / abstractAliases
//   rebinding()/refresh()          -> reboundCallbacks
//   resolving()/afterResolving()   -> *ResolvingCallbacks（延迟注册的 provider 每请求都会走）
//   tag()                          -> tags
//   extend()                       -> extenders
//   registerDeferredProvider()     -> loadedProviders / deferredServices
//   build()                        -> with / buildStack / instances / resolved
// 问题在于「每请求拷 22 张表」是一笔固定税：暖请求真正写到的通常只有 instances 一两个键，
// bindings/aliases 这种几百项的表连内层数组都要递归拷一遍，纯属白做。
// 现在改成按需：链式属性表只在该键第一次被访问时把全局的值升级到请求级，
// 「拷多少」由「请求真的碰过多少」决定，而不是由「表有多大」决定。
func cloneApplication(v data.Value, ctx data.Context) data.Value {
	cv, ok := v.(*data.ClassValue)
	if !ok || cv == nil {
		return v
	}
	return cv.CloneRequestScoped(ctx)
}

func cloneRouter(v data.Value, ctx data.Context) data.Value {
	cv, ok := v.(*data.ClassValue)
	if !ok || cv == nil {
		return v
	}
	// Router 体量小于 Application，但全量深拷贝仍贵；按键隔离即可。
	return cv.CloneSandboxKeys(ctx, routerSandboxDeepKeys)
}

func bindSandboxContainer(ctx data.Context, app, router data.Value) {
	if app == nil {
		return
	}
	// Octane Worker：请求沙箱必须成为 Container::getInstance()。
	// Container::setInstance($app) 的方法体只有 `static::$instance = $app` 一句，
	// 而 storeSandboxContainerInstance 已经把同一个实例按 PHP 读该属性时用的类名
	// 写进请求 overlay，所以那次 PHP 调用省掉。
	if inst := asValue(app); inst != nil {
		storeSandboxContainerInstance(ctx, inst)
	}
	if router != nil {
		// Router::setContainer($app) 的方法体是 `$this->container = $container;`。
		if rcv, ok := asValue(router).(*data.ClassValue); ok && rcv != nil {
			if _, declared := rcv.GetPropertyStmt("container"); declared {
				_ = rcv.SetProperty("container", app)
			} else {
				_, _ = callObjectMethodInContext(ctx, router, "setContainer", app)
			}
		}
		// Container::instance($abstract, $router) 落在容器上的效果就是
		// `instances[$abstract] = $router`：别名簿记对这两个 abstract 恒为空转
		// （removeAbstractAlias 要求 isset(aliases[$abstract])，而 aliases 的键一律是
		// 别名本身即类名，取自 registerCoreContainerAliases 的 `aliases[类名] = 抽象名`；
		// unset($this->aliases[$abstract]) 同样落空），rebound 也没有注册过 router 的回调。
		setRequestInstance(app, "router", router)
		setRequestInstance(app, "Illuminate\\Routing\\Router", router)
	}
	setFacadeStatics(ctx, app)
}

// setFacadeStatics 对齐 Facade::clearResolvedInstances() 与 Facade::setFacadeApplication($app)。
// 两句都是给 Facade 的静态属性整体赋值（`static::$resolvedInstance = []` / `static::$app = $app`），
// 调用点建的是 Facade 基类实例，late static binding 解析到 Facade，所以 overlay 的键就是基类名。
// 写不进 overlay（非请求路径）时才回落 PHP 调用。
func setFacadeStatics(ctx data.Context, app data.Value) {
	if ctx == nil || app == nil {
		return
	}
	vm := ctx.GetVM()
	if vm == nil {
		return
	}
	stmt, ctl := vm.GetOrLoadClass("Illuminate\\Support\\Facades\\Facade")
	if ctl != nil || stmt == nil {
		return
	}
	if data.StoreRequestStatic(stmt.GetName(), "resolvedInstance", data.NewArrayValue(nil)) {
		data.StoreRequestStatic(stmt.GetName(), "app", app)
		return
	}
	facade := data.NewClassValue(stmt, ctx)
	_, _ = callObjectMethodInContext(ctx, facade, "clearResolvedInstances")
	_, _ = callObjectMethodInContext(ctx, facade, "setFacadeApplication", app)
}

func storeSandboxContainerInstance(ctx data.Context, inst data.Value) {
	names := []string{
		"Illuminate\\Container\\Container",
		"Illuminate\\Foundation\\Application",
	}
	var vm data.VM
	if ctx != nil {
		vm = ctx.GetVM()
	}
	for _, n := range names {
		actual := n
		if vm != nil {
			if stmt, ctl := vm.GetOrLoadClass(n); ctl == nil && stmt != nil {
				actual = stmt.GetName()
			}
		}
		data.StoreRequestStatic(actual, "instance", inst)
	}
}

// cloneRequestServices 把每请求会往自己数组里写的共享单例换成 clone。
// PHP clone：数组属性按值拷贝。否则 Dispatcher::$listeners / View composers
// 跨请求膨胀，异常页 VarDumper 能把进程拖到十几秒甚至 OOM。
//
// Manager 系（auth/cache/session）是同一类问题的另一种表现：
// 它们的 $drivers / $guards / $stores 是**共享对象上的数组**，请求期 $this->driver()
// 一类调用会往里写；更糟的是它们启动时注入的 $this->app 是**全局 Application**，
// 请求期一句 $this->app['xxx'] 就会去写全局 app 的 instances/resolved/buildStack
// （压测实测：跨 goroutine 写全局 app 容器表的 1500+ 次里，很大一部分来自这些服务）。
// 所以克隆它们时**必须**把容器引用重绑到请求级 app。
func cloneRequestServices(ctx data.Context, app data.Value) {
	if app == nil {
		return
	}
	// 顺序很关键：instance() 会触发 Container::rebound()，而 Laravel 自己的回调
	// （AuthServiceProvider::register 里的 rebinding('events')）会去读 $app['auth']->guard()。
	// 边克隆边安装时，回调跑在 auth 还是全局单例的那一刻，
	// 于是把 guard 写进了共享 AuthManager 的 $guards（实测每请求 1 次跨 goroutine 写）。
	pending := make([]clonePending, 0, 6)
	pending = appendPending(ctx, app, pending, "events",
		"Illuminate\\Events\\Dispatcher",
		"Illuminate\\Contracts\\Events\\Dispatcher",
	)
	pending = appendPending(ctx, app, pending, "view",
		"Illuminate\\View\\Factory",
		"Illuminate\\Contracts\\View\\Factory",
	)
	pending = appendPending(ctx, app, pending, "url",
		"Illuminate\\Routing\\UrlGenerator",
		"Illuminate\\Contracts\\Routing\\UrlGenerator",
	)
	pending = appendPending(ctx, app, pending, "auth",
		"Illuminate\\Contracts\\Auth\\Factory",
		"Illuminate\\Auth\\AuthManager",
	)
	pending = appendPending(ctx, app, pending, "cache",
		"Illuminate\\Contracts\\Cache\\Factory",
		"Illuminate\\Cache\\CacheManager",
	)
	pending = appendPending(ctx, app, pending, "session",
		"Illuminate\\Session\\SessionManager",
	)
	// auth.driver 是 `fn ($app) => $app['auth']->guard()`（authserviceprovider.php:39）。
	// 全局 app 上只要解析过一次（预热请求就会），instances 里就留着一个用全局 app 造出来的
	// SessionGuard，clone 时对象按引用共享 —— 所有请求就会共用同一个 guard（user/session 全串）。
	// 删掉让它按请求 app 重建，与 php-fpm 每请求全新容器一致。
	dropInstances(app, "auth.driver", "Illuminate\\Contracts\\Auth\\Guard")
	// session.store 同理，但共享的是**整份会话状态**：`fn ($app) => $app['session']->driver()`
	// 在全局 app 上被 Warm 解析过一次，全局 instances 里就留着那个 Store 对象。
	// 删掉（本地 instances 数组是深拷贝，删本地即屏蔽全局）后走 resolve 重建，
	// 配合 resetSessionDrivers 得到请求私有的 Store。
	dropInstances(app, "session.store")
	// 安装：直接把克隆体写进请求 app 的 instances 表，不再逐个走 PHP 的
	// Container::instance()（每请求 6 个 abstract + 10 个别名 = 16 次 PHP 调用，
	// 外加匹配到的 rebound 回调一次）。
	//
	// 等价性依据 —— 本应用注册的 rebound 回调只有三处：
	//   routes  RoutingServiceProvider::registerUrlGenerator 的 extend('url') 闭包里
	//   request AuthServiceProvider::registerRequestRebindHandler
	//   events  AuthServiceProvider::registerEventRebindHandler
	// pending 里只有 events 会命中，而它第一句就是
	// `! $app->resolved('auth') || $app['auth']->hasResolvedGuards() === false` 时 return；
	// appendPending 里的 forgetGuards() 保证后者为 false（且请求刚开头 auth 尚未 resolved）
	// —— 回调恒为空转。
	//
	// instance() 余下的别名簿记（removeAbstractAlias / unset($this->aliases[$abstract])）
	// 对这批 abstract 同样是 no-op：core alias 的方向是 aliases[契约] = 'events'，
	// 而这里写的是 instances[abstract] 与 instances[契约]，getAlias 两条路径都指回克隆体。
	// 顺序仍是「先全部 prime，再装别人」：isolateViewEngines 里的 view 直写要能看见克隆体。
	for _, p := range pending {
		primeInstances(app, p)
	}
	isolateViewEngines(ctx, app)
}

// primeInstances 把克隆体先塞进请求 app 的 instances 表，不触发任何容器回调。
// 后续 instance() 安装时，rebound 回调看到的就都是请求级实例了。
//
// GetProperty 返回的是「请求级」的那张表：请求级 Application 的属性表是链式的，
// 容器值第一次被读到就会从全局升级（分离一层）到请求级，所以这里的原地写只落在本请求，
// 不会改到全局 app 的 instances。改动本函数时不要绕过 GetProperty 直接拿全局的表。
func primeInstances(app data.Value, p clonePending) {
	if p.cloned == nil {
		return
	}
	setRequestInstance(app, p.abstract, p.cloned)
	for _, alias := range p.aliases {
		setRequestInstance(app, alias, p.cloned)
	}
}

// setRequestInstance 直写请求级 instances 表，等价于 Container::instance() 里的
// `$this->instances[$abstract] = $instance` 那一句：不走 PHP 调用、不触发 rebound。
// 调用方必须已经确认该 abstract 没有生效的 rebound 回调（见 cloneRequestServices）。
//
// 同 primeInstances：GetProperty 返回的是请求级那张表，写入只落在本请求，
// 不会改到全局 app 的 instances。改动本函数时不要绕过 GetProperty 直接拿全局的表。
func setRequestInstance(app data.Value, abstract string, val data.Value) {
	cv, ok := app.(*data.ClassValue)
	if !ok || cv == nil || val == nil {
		return
	}
	raw, ctl := cv.GetProperty("instances")
	if ctl != nil || raw == nil {
		return
	}
	arr, ok := raw.(*data.ArrayValue)
	if !ok || arr == nil {
		return
	}
	arr.SetStringKey(abstract, val)
}

// requestInstance 取请求 app 上的某个已解析实例：优先直读 instances 表（Go 侧），
// 取不到才回落 PHP 的 $app->make($abstract)。
//
// 等价性：Warm 已经用**同一个 abstract 字符串**把请求无关的重服务解析进全局 instances，
// 请求级表链式回落到它；而 Container::resolve 对已存在的实例就是一句
// `if (isset($this->instances[$abstract])) return $this->instances[$abstract];`。
// 之前那句 make() 多出来的只有三件事，对本函数用到的这批名字都是空转：
//   getAlias($abstract)              —— aliases 的键是别名（类名），不是抽象名
//   fireBeforeResolvingCallbacks()   —— 只有 resolving() 为这些名字注册过回调才有意义
//   instance 未命中时的 build 链      —— 命中就不会走到，真未命中时本函数已回落 make()
// 所以只要 instances 里命中，直读与 make() 返回同一个对象。
//
// 注意 GetProperty 返回的是请求级那张表（链式属性第一次读会从全局升级成本地副本），
// 读它不会改到全局 app。
func requestInstance(ctx data.Context, app data.Value, abstract string) data.Value {
	if cv, ok := app.(*data.ClassValue); ok && cv != nil {
		if raw, ctl := cv.GetProperty("instances"); ctl == nil && raw != nil {
			if arr, ok := raw.(*data.ArrayValue); ok && arr != nil {
				if z, ok := arr.LookupZValByStringKey(abstract); ok && z != nil && z.Value != nil {
					return z.Value
				}
			}
		}
	}
	raw, ctl := callObjectMethodInContext(ctx, app, "make", data.NewStringValue(abstract))
	if ctl != nil {
		return nil
	}
	return asValue(raw)
}

// dropInstances 从请求 app 的 instances 表里删掉继承自全局 app 的条目。
// 用于「请求态」单例（如 auth.driver，其值是一个 guard 对象）：这类服务在 php-fpm
// 下每请求都重新解析，跟着全局 app 一起继承过来等于跨请求共享对象。
//
// 同 primeInstances：GetProperty 给的是请求级那张表，删只删本请求这一份。
func dropInstances(app data.Value, abstracts ...string) {
	cv, ok := app.(*data.ClassValue)
	if !ok || cv == nil {
		return
	}
	raw, ctl := cv.GetProperty("instances")
	if ctl != nil || raw == nil {
		return
	}
	arr, ok := raw.(*data.ArrayValue)
	if !ok || arr == nil {
		return
	}
	for _, abstract := range abstracts {
		arr.UnsetKey(data.NewStringValue(abstract))
	}
}

// clonePending 是一个已克隆好、待安装到请求 app 的服务。
type clonePending struct {
	abstract string
	aliases  []string
	cloned   *data.ClassValue
}

// appendPending 克隆 abstract 对应的单例（不改容器），返回追加后的列表。
func appendPending(ctx data.Context, app data.Value, list []clonePending, abstract string, aliases ...string) []clonePending {
	raw := requestInstance(ctx, app, abstract)
	cv, ok := raw.(*data.ClassValue)
	if !ok || cv == nil {
		return list
	}
	cloned := cv.CloneSandbox(ctx)
	rebindContainer(ctx, cloned, app)
	if abstract == "auth" {
		// AuthManager::$guards 缓存的是 guard 实例（SessionGuard 持有 user / session store）。
		// 全局 app 上解析过 auth.driver 就会在 $guards 里留下一个用全局 app 造的 guard，
		// clone 时对象按引用共享 → 所有请求共用同一个 guard。forgetGuards() 后
		// 每请求首次 guard() 才按请求 app 重建；同时让 rebinding('events') 回调
		// 因 hasResolvedGuards()===false 提前返回，省掉每请求一次 guard 构造。
		if _, ok := cloned.GetMethod("forgetGuards"); ok {
			_, _ = callObjectMethodInContext(ctx, cloned, "forgetGuards")
		}
	}
	if abstract == "session" {
		resetSessionDrivers(cloned)
	}
	if abstract == "view" {
		// PHP clone 不会改 shared['__env']=$this。Blade 编译视图用 $__env，
		// View::render 用 View::$factory。两者必须是同一实例，否则嵌套 table
		// 在旧 Factory 上 flushStateIfDoneRendering 会清掉 page 的 componentStack（View []）。
		_, _ = callObjectMethodInContext(ctx, cloned, "share", data.NewStringValue("__env"), cloned)
		// 源 Factory 可能正被别的请求渲染（renderCount/componentStack 非空）。
		// PHP clone 会把这些状态拷过来；不 flush 就会 flushStateIfDoneRendering 清错栈，
		// Livewire 得到注释/空 HTML → RootTagMissing，异常页再被 HtmlDumper 拖死。
		_, _ = callObjectMethodInContext(ctx, cloned, "flushState")
	}
	return append(list, clonePending{abstract: abstract, aliases: aliases, cloned: cloned})
}

// resetSessionDrivers 清掉克隆出来的 SessionManager 的 $drivers 缓存，
// 逼每个请求自己 driver() 出一个 Store（Manager::driver 是 `$this->drivers[$d] ??= createDriver()`）。
//
// 为什么必须清：Warm 在**全局 app** 上解析过一次 'session.store'，那会走
// `$app['session']->driver()`，于是全局 SessionManager 的 $drivers 里留下一个 Store。
// 请求级 manager 是 clone 出来的，$drivers 数组按值带过来，但数组里的 **Store 是同一个对象**
// （PHP 数组按值 / 对象按引用），StartSession 的 `$this->manager->driver()` 又直接命中缓存，
// 于是所有请求共用同一个 Store：$attributes（会话数据）/ $id / handler 全部共享。
// 实测（examples/laravel13 的探针路由）：两个**全新 cookie**、session id 完全不同的请求，
// 后一个进 handler 时能看到前一个写的键与 _previous.url；SESSION_DRIVER=array 时
// ArraySessionHandler 里那个共享数组被并发写，直接
// `fatal error: concurrent map writes`（value_array.go:SetStringKey ← IndexExpression）。
//
// 清掉后 `$drivers[$name] ??= createDriver()` 重新建 Store / handler，等于 php-fpm 每请求
// 一份的状态；成本只有一次 Store + handler 构造（driver 名以外的配置都还在）。
func resetSessionDrivers(cv *data.ClassValue) {
	if cv == nil {
		return
	}
	// $drivers 声明在 Support\Manager 上，GetPropertyStmt 会沿 extend 链找到；
	// 找不到（改名/换实现）就保持原样，行为回落到旧实现。
	if _, declared := cv.GetPropertyStmt("drivers"); !declared {
		return
	}
	_ = cv.SetProperty("drivers", data.NewArrayValue(nil))
}

// rebindSetters 是「容器引用 setter → 它写的属性」的对照表。
// 这些方法体在 vendor 里都只有一句赋值，可以直接写属性：
//   Support\Manager::setContainer            -> $this->container = $container
//   AuthManager/CacheManager::setApplication -> $this->app = $app
// （本函数只在这两个类的实例上被调用，见 cloneRequestServices 的 pending 名单。）
var rebindSetters = []struct {
	method   string
	property string
}{
	{"setContainer", "container"},
	{"setApplication", "app"},
}

// rebindContainer 把克隆体上的「容器/应用」引用改指到请求级 app。
// 各家 setter 名不一致，而且有几个压根没有 setter（Events\Dispatcher 的 $container、
// UrlGenerator 连容器属性都没有），只能写属性。
// 不重绑 = 克隆体仍握着全局 Application，请求期 $this->container->make()、
// $this->app['x'] 都会打到全局 app 的容器表上，直接触发 concurrent map writes。
//
// 改 Go 直写（原为逐个 callObjectMethodInContext）：上面两个 setter 的方法体只有一句赋值，
// 而 GetPropertyStmt 会沿 extend 链找到父类声明（session 的 $container 声明在 Support\Manager）。
// 万一属性声明找不到（改名 / 换实现），才回落 PHP 调用，行为与旧实现一致。
func rebindContainer(ctx data.Context, cv *data.ClassValue, app data.Value) {
	if cv == nil || app == nil {
		return
	}
	for _, s := range rebindSetters {
		if _, ok := cv.GetMethod(s.method); !ok {
			continue
		}
		if _, declared := cv.GetPropertyStmt(s.property); declared {
			_ = cv.SetProperty(s.property, app)
			return
		}
		if _, ctl := callObjectMethodInContext(ctx, cv, s.method, app); ctl == nil {
			return
		}
	}
	for _, name := range []string{"container", "app"} {
		if _, ok := cv.GetPropertyStmt(name); ok {
			_ = cv.SetProperty(name, app)
		}
	}
}

// isolateViewEngines 每请求克隆 EngineResolver 并丢掉已解析的 blade/php 引擎。
// 并发 livewire/update 若共用 CompilerEngine::$lastCompiled / 输出缓冲，
// 会得到空 HTML → Livewire RootTagMissing（仪表盘 widget 一直 Loading）。
func isolateViewEngines(ctx data.Context, app data.Value) {
	cv, ok := requestInstance(ctx, app, "view.engine.resolver").(*data.ClassValue)
	if !ok || cv == nil {
		return
	}
	cloned := cv.CloneSandbox(ctx)
	_, _ = callObjectMethodInContext(ctx, cloned, "forget", data.NewStringValue("blade"))
	_, _ = callObjectMethodInContext(ctx, cloned, "forget", data.NewStringValue("php"))
	// 同 cloneRequestServices：view.engine.resolver 没有注册过 rebound 回调，直写 instances 即可。
	setRequestInstance(app, "view.engine.resolver", cloned)

	viewRaw := requestInstance(ctx, app, "view")
	if viewRaw == nil {
		return
	}
	view, ok := asValue(viewRaw).(*data.ClassValue)
	if !ok || view == nil {
		return
	}
	_ = view.SetProperty("engines", cloned)
}

// Terminate 执行 Laravel 请求结束生命周期。
func Terminate(ctx data.Context, kernel *data.ClassValue, request, response data.Value) data.Control {
	if kernel == nil {
		return nil
	}
	_, control := callObjectMethodInContext(ctx, kernel, "terminate", request, response)
	return control
}

// 注：Octane Worker 的 finally 里有一句「丢掉常驻 Blade/PHP 引擎」，这里**不需要**——
// 那个 reset 是为「resolver 是跨请求共享单例」准备的，而本实现每请求都在
// isolateViewEngines 里从常驻 resolver 克隆出一份私有的、并且已经 forget 掉 blade/php 的副本，
// 请求结束整份副本就丢了，引擎不可能把脏缓冲带到下一请求。
// 曾经的 ResetViewEngines 作用在请求 app 上（make 命中的正是这份私有副本），
// 每请求白打 3 次 PHP 调用（make + forget×2）后把副本再 forget 一遍，已删除。

