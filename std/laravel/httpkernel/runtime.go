package httpkernel

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// Resolve 从 Laravel 容器解析 Go 实现的 Illuminate\Foundation\Http\Kernel。
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
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 容器未返回 Illuminate\\Foundation\\Http\\Kernel"))
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
	// 而预热发生在全局 app 上，等于把一份会话状态挂到全局。
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
	var app, router data.Value
	if st != nil {
		app, router = st.app, st.router
	} else {
		app, _ = kernel.GetProperty("app")
		router, _ = kernel.GetProperty("router")
	}
	if app == nil {
		return
	}
	for _, abstract := range warmAbstracts {
		bound, ctl := callObjectMethodInContext(ctx, app, "bound", data.NewStringValue(abstract))
		if ctl != nil || !isTrueValue(bound) {
			continue
		}
		_, _ = callObjectMethodInContext(ctx, app, "make", data.NewStringValue(abstract))
	}
	// 路由的 Symfony 编译结果是纯粹的启动期常量，先算出来（详见 warmRouteCompiles）。
	warmRouteCompiles(ctx, router)
	// Telescope 的开关、匹配 pattern、实例同样都是启动期常量。在全局 Application 上解析一次，
	// 请求期就只剩一次 $request->is()；放在这里（而非首个请求里）也避免并发首请求在
	// once.Do 上排队，以及把 Telescope 实例建进某个请求沙箱。
	if st != nil && st.tel != nil {
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
	var scope *data.RequestObjectScope
	if provider, ok := ctx.GetVM().(data.RequestScopeProvider); ok {
		scope = provider.RequestObjectScope()
	} else {
		scope = data.NewRequestObjectScope(ctx)
	}
	request := scope.Object(kernel)
	state, _ := request.GetSource().(*kernelState)
	var app, router data.Value
	if state == nil {
		app, _ = request.GetProperty("app")
		router, _ = request.GetProperty("router")
	} else {
		syncProperties(request, state)
		app, router = state.app, state.router
	}
	bindSandboxContainer(ctx, app, router)
	setRequestInstance(app, fqnKernel, request)
	setRequestInstance(app, fqnKernelContract, request)
	if _, exists := methodExists(app, "forgetScopedInstances"); exists {
		if _, ctl := callObjectMethodInContext(ctx, app, "forgetScopedInstances"); ctl != nil {
			ctx.GetVM().ThrowControl(ctl)
		}
	}
	return request
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

// setRequestInstance 直写请求级 instances 表，等价于 Container::instance() 里的
// `$this->instances[$abstract] = $instance` 那一句：不走 PHP 调用、不触发 rebound。
// 此入口仅用于安装请求的 Kernel 和 Router 身份。
//
// GetProperty 返回的是请求级那张表，写入只落在本请求，
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
//
//	getAlias($abstract)              —— aliases 的键是别名（类名），不是抽象名
//	fireBeforeResolvingCallbacks()   —— 只有 resolving() 为这些名字注册过回调才有意义
//	instance 未命中时的 build 链      —— 命中就不会走到，真未命中时本函数已回落 make()
//
// 所以只要 instances 里命中，直读与 make() 返回同一个对象。
//
// 注意 GetProperty 返回的是请求级那张表（链式属性第一次读会从全局升级成本地副本），
// 读它不会改到全局 app。
func requestInstance(ctx data.Context, app data.Value, abstract string) data.Value {
	if cv, ok := app.(*data.ClassValue); ok && cv != nil {
		if raw, ctl := cv.GetProperty("instances"); ctl == nil && raw != nil {
			if arr, ok := raw.(*data.ArrayValue); ok && arr != nil {
				if z, ok := arr.LookupZValByStringKey(abstract); ok && z != nil && z.ReadValue() != nil {
					return z.ReadValue()
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

// Terminate 执行 Laravel 请求结束生命周期。
func Terminate(ctx data.Context, kernel *data.ClassValue, request, response data.Value) data.Control {
	if kernel == nil {
		return nil
	}
	_, control := callObjectMethodInContext(ctx, kernel, "terminate", request, response)
	return control
}
