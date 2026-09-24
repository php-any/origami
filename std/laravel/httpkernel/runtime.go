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

// warmAbstracts 是「请求无关」的重服务。请求沙箱按 appSandboxDeepKeys 深拷贝 instances，
// 这些服务只要在全局 Application 里已经解析过，沙箱的副本就是本地命中；
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
	"session.store",
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

// appSandboxDeepKeys：Container/Application 启动后只读槽共享，仅隔离会按请求改写的表。
//
// 哪些必须隔离：只要请求期可能调用 Container 的写方法，对应的表就得是请求私有的，
// 否则并发请求会同时写同一个数组（实测 fatal error: concurrent map writes），
// 并把 bindings/aliases 改坏（表现为
// "Target [Illuminate\Contracts\Routing\Registrar] is not instantiable"）。
// 请求期会写到的表来自这些 API：
//   bind()/instance()/alias()      -> bindings / aliases / abstractAliases
//   rebinding()/refresh()          -> reboundCallbacks
//   resolving()/afterResolving()   -> *ResolvingCallbacks（延迟注册的 provider 每请求都会走）
//   tag()                          -> tags
//   extend()                       -> extenders
//   registerDeferredProvider()     -> loadedProviders / deferredServices
//   build()                        -> with / buildStack / instances / resolved
var appSandboxDeepKeys = []string{
	"instances",
	"resolved",
	"scopedInstances",
	"buildStack",
	"with",
	"bindings",
	"aliases",
	"abstractAliases",
	"reboundCallbacks",
	"methodBindings",
	"tags",
	"extenders",
	"resolvingCallbacks",
	"afterResolvingCallbacks",
	"globalResolvingCallbacks",
	"globalAfterResolvingCallbacks",
	"beforeResolvingCallbacks",
	"loadedProviders",
	"deferredServices",
	"serviceAliases",
	"checkedForAttributeBindings",
	"checkedForSingletonOrScopedAttributes",
}

// routerSandboxDeepKeys：Router 上随请求变化的少量状态；RouteCollection 启动后只读，共享。
var routerSandboxDeepKeys = []string{
	"current",
	"currentRequest",
}

func cloneApplication(v data.Value, ctx data.Context) data.Value {
	cv, ok := v.(*data.ClassValue)
	if !ok || cv == nil {
		return v
	}
	return cv.CloneSandboxKeys(ctx, appSandboxDeepKeys)
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
	_, _ = callObjectMethodInContext(ctx, app, "setInstance", app)
	// Octane Worker：请求沙箱必须成为 Container::getInstance()。
	// PHP setInstance 若因类型提示/static:: 未写入 overlay，csrf_token() 会落到
	// 进程级 Application（session 未 start），登录页 meta/data-csrf 变成空串。
	if inst := asValue(app); inst != nil {
		storeSandboxContainerInstance(ctx, inst)
	}
	if router != nil {
		_, _ = callObjectMethodInContext(ctx, router, "setContainer", app)
		_, _ = callObjectMethodInContext(ctx, app, "instance", data.NewStringValue("router"), router)
		_, _ = callObjectMethodInContext(ctx, app, "instance", data.NewStringValue("Illuminate\\Routing\\Router"), router)
	}
	vm := ctx.GetVM()
	if vm == nil {
		return
	}
	stmt, ctl := vm.GetOrLoadClass("Illuminate\\Support\\Facades\\Facade")
	if ctl != nil || stmt == nil {
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
	// 先把克隆体写进 instances 表，再逐个 instance() 安装。
	// instance() 会触发 Container::rebound()，而 Laravel 自己的回调会顺着 $app['auth']
	// 去读服务（AuthServiceProvider 的 rebinding('events') -> $app['auth']->guard()）。
	// instances 表是从全局 app 浅拷贝来的（对象按引用共享），安装顺序里 events 在 auth 前，
	// 若边装边触发回调，回调那一刻 $app['auth'] 还是全局单例，
	// guard() 就会把 $guards 写进共享 AuthManager（实测每请求 1 次跨 goroutine 写）。
	for _, p := range pending {
		primeInstances(app, p)
	}
	for _, p := range pending {
		bindService(ctx, app, p)
	}
	isolateViewEngines(ctx, app)
}

// primeInstances 把克隆体先塞进请求 app 的 instances 表，不触发任何容器回调。
// 后续 instance() 安装时，rebound 回调看到的就都是请求级实例了。
func primeInstances(app data.Value, p clonePending) {
	cv, ok := app.(*data.ClassValue)
	if !ok || cv == nil || p.cloned == nil {
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
	arr.SetStringKey(p.abstract, p.cloned)
	for _, alias := range p.aliases {
		arr.SetStringKey(alias, p.cloned)
	}
}

// dropInstances 从请求 app 的 instances 表里删掉继承自全局 app 的条目。
// 用于「请求态」单例（如 auth.driver，其值是一个 guard 对象）：这类服务在 php-fpm
// 下每请求都重新解析，跟着全局 app 一起继承过来等于跨请求共享对象。
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

func cloneAndBind(ctx data.Context, app data.Value, abstract string, aliases ...string) {
	list := appendPending(ctx, app, nil, abstract, aliases...)
	for _, p := range list {
		bindService(ctx, app, p)
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
	raw, ctl := callObjectMethodInContext(ctx, app, "make", data.NewStringValue(abstract))
	if ctl != nil {
		return list
	}
	cv, ok := asValue(raw).(*data.ClassValue)
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

// bindService 把克隆体 instance() 进请求 app。
// instance() 会触发 Container::rebound()，Laravel 注册的回调（例如
// AuthServiceProvider 里的 rebinding('events')）会读 $app['auth']->guard()，
// 所以必须等所有服务都克隆完再统一安装，否则回调拿到的是全局单例。
func bindService(ctx data.Context, app data.Value, p clonePending) {
	_, _ = callObjectMethodInContext(ctx, app, "instance", data.NewStringValue(p.abstract), p.cloned)
	for _, alias := range p.aliases {
		_, _ = callObjectMethodInContext(ctx, app, "instance", data.NewStringValue(alias), p.cloned)
	}
}

// rebindContainer 把克隆体上的「容器/应用」引用改指到请求级 app。
// 各家 setter 名不一致，而且有几个压根没有 setter（Events\Dispatcher 的 $container、
// Manager 系的 $container 只有 Manager 子类才有 setContainer）：
//   Support\Manager 子类(session)    -> setContainer(Container)
//   AuthManager / CacheManager / DatabaseManager -> setApplication($app)
//   Events\Dispatcher                -> 没有 setter，只能写 $container
// 不重绑 = 克隆体仍握着全局 Application，请求期 $this->container->make()、
// $this->app['x'] 都会打到全局 app 的容器表上，直接触发 concurrent map writes。
func rebindContainer(ctx data.Context, cv *data.ClassValue, app data.Value) {
	if cv == nil || app == nil {
		return
	}
	for _, name := range []string{"setContainer", "setApplication"} {
		if _, ok := cv.GetMethod(name); !ok {
			continue
		}
		if _, ctl := callObjectMethodInContext(ctx, cv, name, app); ctl == nil {
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
	raw, ctl := callObjectMethodInContext(ctx, app, "make", data.NewStringValue("view.engine.resolver"))
	if ctl != nil {
		return
	}
	cv, ok := asValue(raw).(*data.ClassValue)
	if !ok || cv == nil {
		return
	}
	cloned := cv.CloneSandbox(ctx)
	_, _ = callObjectMethodInContext(ctx, cloned, "forget", data.NewStringValue("blade"))
	_, _ = callObjectMethodInContext(ctx, cloned, "forget", data.NewStringValue("php"))
	_, _ = callObjectMethodInContext(ctx, app, "instance", data.NewStringValue("view.engine.resolver"), cloned)

	viewRaw, vctl := callObjectMethodInContext(ctx, app, "make", data.NewStringValue("view"))
	if vctl != nil {
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

// ResetViewEngines 对齐 Octane Worker finally：丢掉常驻 Blade/PHP 引擎，
// 避免编译引擎跨请求持有脏缓冲，下一请求欢迎页变成空 body。
func ResetViewEngines(ctx data.Context, kernel *data.ClassValue) {
	if kernel == nil {
		return
	}
	src, _ := kernel.GetSource().(*kernelState)
	if src == nil || src.app == nil {
		return
	}
	resolver, ctl := callObjectMethodInContext(ctx, src.app, "make", data.NewStringValue("view.engine.resolver"))
	obj := asValue(resolver)
	if ctl != nil || obj == nil {
		return
	}
	_, _ = callObjectMethodInContext(ctx, obj, "forget", data.NewStringValue("blade"))
	_, _ = callObjectMethodInContext(ctx, obj, "forget", data.NewStringValue("php"))
}
