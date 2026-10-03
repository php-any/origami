package httpkernel

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/illuminate/foundation"
)

func kernelMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("app", 0, nil, data.NewBaseType("Illuminate\\Foundation\\Application")),
				param("router", 1, nil, data.NewBaseType("Illuminate\\Routing\\Router")),
			},
			[]data.Variable{
				variable("app", 0, data.NewBaseType("Illuminate\\Foundation\\Application")),
				variable("router", 1, data.NewBaseType("Illuminate\\Routing\\Router")),
			},
			nil, kernelConstruct),

		pubMethod("bootstrap", nil, nil, nil, kernelBootstrap),

		pubMethod("handle",
			[]data.GetValue{param("request", 0, nil, nil)},
			[]data.Variable{variable("request", 0, nil)},
			nil, kernelHandle),

		pubMethod("terminate",
			[]data.GetValue{
				param("request", 0, nil, nil),
				param("response", 1, nil, nil),
			},
			[]data.Variable{
				variable("request", 0, nil),
				variable("response", 1, nil),
			},
			nil, kernelTerminate),

		pubMethod("getApplication", nil, nil, nil, kernelGetApplication),

		pubMethod("hasMiddleware",
			[]data.GetValue{param("middleware", 0, nil, nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			data.NewBaseType("bool"), kernelHasMiddleware),

		pubMethod("prependMiddleware",
			[]data.GetValue{param("middleware", 0, nil, nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			nil, kernelPrependMiddleware),

		pubMethod("pushMiddleware",
			[]data.GetValue{param("middleware", 0, nil, nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			nil, kernelPushMiddleware),

		pubMethod("setGlobalMiddleware",
			[]data.GetValue{param("middleware", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			nil, kernelSetGlobalMiddleware),

		// Laravel: prependMiddlewareToGroup / appendMiddlewareToGroup
		// 同时提供用户约定的 prependToMiddlewareGroup / appendToMiddlewareGroup 别名。
		pubMethod("prependMiddlewareToGroup",
			[]data.GetValue{
				param("group", 0, nil, nil),
				param("middleware", 1, nil, nil),
			},
			[]data.Variable{
				variable("group", 0, nil),
				variable("middleware", 1, nil),
			},
			nil, kernelPrependToMiddlewareGroup),
		pubMethod("appendMiddlewareToGroup",
			[]data.GetValue{
				param("group", 0, nil, nil),
				param("middleware", 1, nil, nil),
			},
			[]data.Variable{
				variable("group", 0, nil),
				variable("middleware", 1, nil),
			},
			nil, kernelAppendToMiddlewareGroup),
		pubMethod("prependToMiddlewareGroup",
			[]data.GetValue{
				param("group", 0, nil, nil),
				param("middleware", 1, nil, nil),
			},
			[]data.Variable{
				variable("group", 0, nil),
				variable("middleware", 1, nil),
			},
			nil, kernelPrependToMiddlewareGroup),
		pubMethod("appendToMiddlewareGroup",
			[]data.GetValue{
				param("group", 0, nil, nil),
				param("middleware", 1, nil, nil),
			},
			[]data.Variable{
				variable("group", 0, nil),
				variable("middleware", 1, nil),
			},
			nil, kernelAppendToMiddlewareGroup),

		pubMethod("setMiddlewareGroups",
			[]data.GetValue{param("groups", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("groups", 0, nil)},
			nil, kernelSetMiddlewareGroups),

		pubMethod("addToMiddlewarePriorityBefore",
			[]data.GetValue{
				param("before", 0, nil, nil),
				param("middleware", 1, nil, nil),
			},
			[]data.Variable{
				variable("before", 0, nil),
				variable("middleware", 1, nil),
			},
			nil, kernelAddToMiddlewarePriorityBefore),

		pubMethod("addToMiddlewarePriorityAfter",
			[]data.GetValue{
				param("after", 0, nil, nil),
				param("middleware", 1, nil, nil),
			},
			[]data.Variable{
				variable("after", 0, nil),
				variable("middleware", 1, nil),
			},
			nil, kernelAddToMiddlewarePriorityAfter),

		// Laravel: prependToMiddlewarePriority / appendToMiddlewarePriority
		// 同时提供 *PriorityList 别名。
		pubMethod("prependToMiddlewarePriority",
			[]data.GetValue{param("middleware", 0, nil, nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			nil, kernelPrependToMiddlewarePriority),
		pubMethod("appendToMiddlewarePriority",
			[]data.GetValue{param("middleware", 0, nil, nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			nil, kernelAppendToMiddlewarePriority),
		pubMethod("prependToMiddlewarePriorityList",
			[]data.GetValue{param("middleware", 0, nil, nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			nil, kernelPrependToMiddlewarePriority),
		pubMethod("appendToMiddlewarePriorityList",
			[]data.GetValue{param("middleware", 0, nil, nil)},
			[]data.Variable{variable("middleware", 0, nil)},
			nil, kernelAppendToMiddlewarePriority),

		pubMethod("setMiddlewarePriority",
			[]data.GetValue{param("priority", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("priority", 0, nil)},
			nil, kernelSetMiddlewarePriority),

		pubMethod("setMiddlewareAliases",
			[]data.GetValue{param("aliases", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("aliases", 0, nil)},
			nil, kernelSetMiddlewareAliases),
	}

	m := make(map[string]data.Method, len(list))
	for _, method := range list {
		m[method.GetName()] = method
	}
	return m, list
}

func kernelConstruct(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	app := indexValue(ctx, 0)
	router := indexValue(ctx, 1)
	if app == nil || router == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: __construct 需要 Application 与 Router"))
	}
	// $this is an evaluation wrapper around the same object, not a distinct
	// application kind. Worker isolation operates on the canonical identity.
	if value, ok := app.(*data.ThisValue); ok {
		app = value.ClassValue.InstanceIdentity()
	}
	if value, ok := router.(*data.ThisValue); ok {
		router = value.ClassValue.InstanceIdentity()
	}
	s.app = app
	s.router = router

	cv := thisClassValue(ctx)
	syncProperties(cv, s)
	if ctl := syncMiddlewareToRouter(s); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func kernelBootstrap(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if ctl := doBootstrap(ctx, s); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func doBootstrap(ctx data.Context, s *kernelState) data.Control {
	if s == nil || s.app == nil {
		return data.NewErrorThrow(nil, fmt.Errorf("httpkernel: Application 未设置"))
	}
	if s.bootstrapped {
		return nil
	}
	ret, ctl := callObjectMethodInContext(ctx, s.app, "hasBeenBootstrapped")
	if ctl != nil {
		return ctl
	}
	if isTruthy(asValue(ret)) {
		s.bootstrapped = true
		return nil
	}
	_, ctl = callObjectMethodInContext(ctx, s.app, "bootstrapWith", stringsToArrayValue(s.bootstrappers))
	if ctl != nil {
		return ctl
	}
	s.bootstrapped = true
	return nil
}

func kernelHandle(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	request := indexValue(ctx, 0)
	if request == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: handle 需要 Request"))
	}
	if _, exists := methodExists(request, "enableHttpMethodParameterOverride"); exists {
		if _, ctl := callObjectMethodInContext(ctx, request, "enableHttpMethodParameterOverride"); ctl != nil {
			return nil, ctl
		}
	}

	// 1) bootstrap
	if ctl := doBootstrap(ctx, s); ctl != nil {
		return nil, ctl
	}

	// 1.5) 常驻 Application 对齐 Octane：每请求重置 Livewire 前端资源标记。
	// 否则 hasRenderedScripts 会跨请求残留，登录页第二次起不再输出 livewire.js，
	// 表单退化成浏览器默认 GET。
	flushLivewireState(ctx, s)
	flushOnceHelper(ctx, s)

	// 2) 绑定 request 到容器（框架语义在 std/laravel/foundation）
	if appCV, ok := s.app.(*data.ClassValue); ok {
		if ctl := foundation.BindRequest(appCV, request); ctl != nil {
			return nil, ctl
		}
	} else {
		_, ctl = callObjectMethodInContext(ctx, s.app, "instance", data.NewStringValue("request"), request)
		if ctl != nil {
			return nil, ctl
		}
	}

	// 2.5) 每次请求重置 Telescope 记录状态：telescope/ignore 路径不记录自身请求
	syncTelescopeRecording(ctx, s, request)

	// 3) 对齐 Foundation\Http\Kernel::handle：路由分发失败时 report + render，
	// 而不是把 NotFoundHttpException 等冒泡成 Fatal。
	response, ctl := dispatchToRouter(ctx, s, request)
	if ctl != nil {
		if response, ctl = renderException(ctx, s, request, ctl); ctl != nil {
			return nil, ctl
		}
	}
	if response == nil {
		return data.NewNullValue(), nil
	}

	// 4) 对齐 Foundation\Http\Kernel::handle：请求处理完成后派发 RequestHandled 事件，
	// 供 Telescope RequestWatcher 等监听器记录请求。
	if ctl := dispatchRequestHandled(ctx, s, request, asValue(response)); ctl != nil {
		return nil, ctl
	}

	return response, nil
}

func flushLivewireState(ctx data.Context, s *kernelState) {
	if s == nil || s.app == nil {
		return
	}
	bound, ctl := callObjectMethodInContext(ctx, s.app, "bound", data.NewStringValue("livewire"))
	if ctl != nil || bound == nil {
		return
	}
	if b, ok := bound.(data.AsBool); ok {
		okv, err := b.AsBool()
		if err != nil || !okv {
			return
		}
	}
	lw, ctl := callObjectMethodInContext(ctx, s.app, "make", data.NewStringValue("livewire"))
	if ctl != nil || lw == nil {
		return
	}
	if obj := asValue(lw); obj != nil {
		_, _ = callObjectMethodInContext(ctx, obj, "flushState")
	}
}

// flushOnceHelper 对齐 Octane FlushOnce：Illuminate\Support\Once 用 WeakMap 挂在单例上，
// 不 flush 则 Filament/异常页 frames() 等 once() 结果跨请求串态。
func flushOnceHelper(ctx data.Context, s *kernelState) {
	if ctx == nil {
		return
	}
	vm := ctx.GetVM()
	if vm == nil {
		return
	}
	stmt, ctl := vm.GetOrLoadClass("Illuminate\\Support\\Once")
	if ctl != nil || stmt == nil {
		return
	}
	once := data.NewClassValue(stmt, ctx)
	_, _ = callObjectMethodInContext(ctx, once, "flush")
	if s == nil || s.app == nil {
		return
	}
	vite, ctl := callObjectMethodInContext(ctx, s.app, "make", data.NewStringValue("Illuminate\\Foundation\\Vite"))
	if ctl != nil {
		return
	}
	if obj := asValue(vite); obj != nil {
		_, _ = callObjectMethodInContext(ctx, obj, "flush")
	}
}

const fqnRequestHandled = "Illuminate\\Foundation\\Http\\Events\\RequestHandled"

// dispatchRequestHandled 派发 Illuminate\Foundation\Http\Events\RequestHandled 事件，
// 对齐 Laravel Kernel::handle 在请求处理完成后的尾部行为。
func dispatchRequestHandled(ctx data.Context, s *kernelState, request, response data.Value) data.Control {
	if s == nil || s.app == nil || request == nil || response == nil {
		return nil
	}
	// 无监听器快路径：整段（make 事件 + make dispatcher + dispatch）等价于
	// 「构造一个没人接的事件对象再丢掉」，没有监听器时直接不构造。
	//
	// 语义依据 Events\Dispatcher::dispatch（Dispatcher.php:171）：
	//   [$event, $payload] = $this->parseEventAndPayload(...);   // 事件对象已是实参，这里只是取名
	//   if ($this->shouldBroadcast($payload)) { ... }            // 要求 $payload[0] instanceof ShouldBroadcast
	//   foreach ($this->getListeners($event) as $listener) { ... }  // 空 ⇒ $responses 空
	// RequestHandled（Foundation/Http/Events/RequestHandled.php）是普通类，不实现 ShouldBroadcast，
	// 所以监听器为空时 dispatch 恒返回空数组、无任何副作用。
	//
	// 监听器表读的是**请求级** dispatcher（通用对象图策略生成的副本，其 listeners
	// 未命中时回落到全局）：本请求中途才 listen 的事件也能看见，不依赖启动期快照。
	// Livewire 的资源自动注入和 Telescope RequestWatcher 都依赖这个事件。
	dispatcher := requestInstance(ctx, s.app, "events")
	if !hasEventListeners(ctx, dispatcher, fqnRequestHandled) {
		return nil
	}
	// $app->make(RequestHandled::class, ['request' => $request, 'response' => $response])
	params := data.NewArrayValueFromSlots([]*data.ZVal{
		data.NewNamedZVal("request", request),
		data.NewNamedZVal("response", response),
	})
	event, ctl := callObjectMethodInContext(ctx, s.app, "make",
		data.NewStringValue("Illuminate\\Foundation\\Http\\Events\\RequestHandled"), params)
	if ctl != nil || event == nil {
		return ctl
	}
	// $app['events']->dispatch($event)
	_, ctl = callObjectMethodInContext(ctx, dispatcher, "dispatch", asValue(event))
	return ctl
}

// hasEventListeners 通过 Dispatcher 的公开 API 判断事件是否有监听器。
// 不能直接读取 PHP 的 $listeners 属性：Origami 预加载的 Go Dispatcher 为避免热路径
// 上反复转换 PHP 数组，把监听器保存在 dispatcherState 中，声明属性只是兼容反射的空壳。
// 直接读属性会把 Livewire 注册的 RequestHandled 监听器误判为空，导致其前端资源未注入。
// 调用失败时保守返回 true，确保只损失快路径而不改变事件语义。
func hasEventListeners(ctx data.Context, dispatcher data.Value, name string) bool {
	ret, ctl := callObjectMethodInContext(ctx, dispatcher, "hasListeners", data.NewStringValue(name))
	if ctl != nil || ret == nil {
		return true
	}
	return isTruthy(asValue(ret))
}

// dispatchToRouter 对齐 Foundation\Http\Kernel::dispatchToRouter：
// 将请求交给 Router::dispatch，由其在 Pipeline 中执行路由级中间件
// （web 组的 StartSession / ShareErrorsFromSession 等），再运行路由。
func dispatchToRouter(ctx data.Context, s *kernelState, request data.Value) (data.GetValue, data.Control) {
	skipped, ctl := callObjectMethodInContext(ctx, s.app, "shouldSkipMiddleware")
	if ctl != nil {
		return nil, ctl
	}
	destination := &routerDestination{state: s}
	if isTruthy(asValue(skipped)) || len(s.middleware) == 0 {
		return destination.dispatch(ctx, request)
	}
	class, ctl := ctx.GetVM().GetOrLoadClass("Illuminate\\Routing\\Pipeline")
	if ctl != nil {
		return nil, ctl
	}
	pipeline := data.NewClassValue(class, ctx.CreateBaseContext())
	for _, call := range []struct {
		method string
		arg    data.Value
	}{
		{"__construct", s.app},
		{"send", request},
		{"through", stringsToArrayValue(s.middleware)},
	} {
		if _, ctl := callObjectMethodInContext(ctx, pipeline, call.method, call.arg); ctl != nil {
			return nil, ctl
		}
	}
	return callObjectMethodInContext(ctx, pipeline, "then", data.NewFuncValue(destination))
}

type routerDestination struct{ state *kernelState }

func (*routerDestination) GetName() string { return "kernel_dispatch" }
func (*routerDestination) GetParams() []data.GetValue {
	return []data.GetValue{param("request", 0, nil, nil)}
}
func (*routerDestination) GetVariables() []data.Variable {
	return []data.Variable{variable("request", 0, nil)}
}
func (d *routerDestination) Call(ctx data.Context) (data.GetValue, data.Control) {
	return d.dispatch(ctx, indexValue(ctx, 0))
}
func (d *routerDestination) dispatch(ctx data.Context, request data.Value) (data.GetValue, data.Control) {
	if _, ctl := callObjectMethodInContext(ctx, d.state.app, "instance", data.NewStringValue("request"), request); ctl != nil {
		return nil, ctl
	}
	return callObjectMethodInContext(ctx, d.state.router, "dispatch", request)
}

const fqnExceptionHandler = "Illuminate\\Contracts\\Debug\\ExceptionHandler"

// renderException 对齐 Foundation\Http\Kernel::reportException / renderException。
func renderException(ctx data.Context, s *kernelState, request data.Value, thrown data.Control) (data.GetValue, data.Control) {
	if exit, ok := thrown.(data.ExitControl); ok && exit.IsExit() {
		return nil, thrown
	}
	exception, ok := exceptionFromControl(thrown)
	if !ok || s.app == nil {
		// This control is returned to the SAPI caller, which owns diagnostics.
		// Logging here as well would report the same failure twice.
		return nil, thrown
	}

	handlerRet, ctl := callObjectMethodInContext(ctx, s.app, "make", data.NewStringValue(fqnExceptionHandler))
	if ctl != nil {
		return fallbackExceptionResponse(ctx, exception, thrown)
	}
	handler := asValue(handlerRet)
	if handler == nil {
		return fallbackExceptionResponse(ctx, exception, thrown)
	}

	// report 失败不阻断 render（HttpException 通常本就不会上报）。
	_, _ = callObjectMethodInContext(ctx, handler, "report", exception)

	response, ctl := callObjectMethodInContext(ctx, handler, "render", request, exception)
	if ctl != nil || response == nil {
		return fallbackExceptionResponse(ctx, exception, thrown)
	}
	return response, nil
}

func exceptionFromControl(ctl data.Control) (data.Value, bool) {
	tv, ok := ctl.(*data.ThrowValue)
	if !ok || tv == nil {
		return nil, false
	}
	if tv.Object != nil {
		return tv.Object, true
	}
	return nil, false
}

func logHandleException(thrown data.Control) {
	if thrown == nil {
		return
	}
	msg := thrown.AsString()
	if len(msg) > 400 {
		msg = msg[:400]
	}
	fmt.Fprintf(os.Stderr, "origami Handle exception: %s\n", msg)
}

// fallbackExceptionResponse 在 Exception Handler 不可用时，尽量按 HttpException 状态码返回。
func fallbackExceptionResponse(ctx data.Context, exception data.Value, original data.Control) (data.GetValue, data.Control) {
	logHandleException(original)
	status := 500
	message := "Server Error"
	if exception != nil {
		if ret, ctl := callObjectMethodInContext(ctx, exception, "getStatusCode"); ctl == nil {
			if v := asValue(ret); v != nil {
				if n, err := strconv.Atoi(v.AsString()); err == nil {
					status = n
				}
			}
		}
		if ret, ctl := callObjectMethodInContext(ctx, exception, "getMessage"); ctl == nil {
			if v := asValue(ret); v != nil && v.AsString() != "" {
				message = v.AsString()
			}
		}
	}
	if status < 100 || status >= 600 {
		status = 500
	}

	class, ctl := ctx.GetVM().GetOrLoadClass("Illuminate\\Http\\Response")
	if ctl != nil {
		return nil, original
	}
	cv := data.NewClassValue(class, ctx.CreateBaseContext())
	construct := class.GetConstruct()
	if construct == nil {
		return nil, original
	}
	fnCtx := cv.CreateContext(construct.GetVariables())
	vars := construct.GetVariables()
	args := []data.Value{
		data.NewStringValue(message),
		data.NewIntValue(status),
	}
	for i, arg := range args {
		if i >= len(vars) {
			break
		}
		if acl := fnCtx.SetVariableValue(vars[i], arg); acl != nil {
			return nil, original
		}
	}
	if _, ctl = construct.Call(fnCtx); ctl != nil {
		return nil, original
	}
	return cv, nil
}

func kernelTerminate(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	request, response := indexValue(ctx, 0), indexValue(ctx, 1)
	event, ctl := callObjectMethodInContext(ctx, s.app, "make", data.NewStringValue("Illuminate\\Foundation\\Events\\Terminating"))
	if ctl != nil {
		return nil, ctl
	}
	dispatcher, ctl := callObjectMethodInContext(ctx, s.app, "make", data.NewStringValue("events"))
	if ctl != nil {
		return nil, ctl
	}
	if _, ctl := callObjectMethodInContext(ctx, asValue(dispatcher), "dispatch", asValue(event)); ctl != nil {
		return nil, ctl
	}
	// Laravel gathers route middleware first, followed by global middleware, and
	// resolves a new instance for terminate (unless the container binding is shared).
	skipped, ctl := callObjectMethodInContext(ctx, s.app, "shouldSkipMiddleware")
	if ctl != nil {
		return nil, ctl
	}
	if !isTruthy(asValue(skipped)) {
		middlewares := []string{}
		route, ctl := callObjectMethodInContext(ctx, request, "route")
		if ctl != nil {
			return nil, ctl
		}
		if rv := asValue(route); rv != nil && isTruthy(rv) {
			gathered, ctl := callObjectMethodInContext(ctx, s.router, "gatherRouteMiddleware", rv)
			if ctl != nil {
				return nil, ctl
			}
			middlewares = stringListFromValue(asValue(gathered))
		}
		middlewares = append(middlewares, s.middleware...)
		for _, middleware := range middlewares {
			name, _, _ := strings.Cut(middleware, ":")
			made, ctl := callObjectMethodInContext(ctx, s.app, "make", data.NewStringValue(name))
			if ctl != nil {
				return nil, ctl
			}
			instance := asValue(made)
			if _, exists := methodExists(instance, "terminate"); exists {
				if _, ctl := callObjectMethodInContext(ctx, instance, "terminate", request, response); ctl != nil {
					return nil, ctl
				}
			}
		}
	}
	if s.app != nil {
		if _, exists := methodExists(s.app, "terminate"); exists {
			_, ctl = callObjectMethodInContext(ctx, s.app, "terminate")
			if ctl != nil {
				return nil, ctl
			}
		}
	}
	return data.NewNullValue(), nil
}

func kernelGetApplication(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if s.app == nil {
		return data.NewNullValue(), nil
	}
	return s.app, nil
}

func kernelHasMiddleware(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	mw := valueAsString(indexValue(ctx, 0))
	return data.NewBoolValue(containsString(s.middleware, mw)), nil
}

func kernelPrependMiddleware(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	mw := valueAsString(indexValue(ctx, 0))
	if mw != "" && !containsString(s.middleware, mw) {
		s.middleware = append([]string{mw}, s.middleware...)
		syncProperties(thisClassValue(ctx), s)
	}
	return thisClassValue(ctx), nil
}

func kernelPushMiddleware(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	mw := valueAsString(indexValue(ctx, 0))
	if mw != "" && !containsString(s.middleware, mw) {
		s.middleware = append(s.middleware, mw)
		syncProperties(thisClassValue(ctx), s)
	}
	return thisClassValue(ctx), nil
}

func kernelSetGlobalMiddleware(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s.middleware = stringListFromValue(indexValue(ctx, 0))
	syncProperties(thisClassValue(ctx), s)
	if ctl := syncMiddlewareToRouter(s); ctl != nil {
		return nil, ctl
	}
	return thisClassValue(ctx), nil
}

func kernelPrependToMiddlewareGroup(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	group := valueAsString(indexValue(ctx, 0))
	mw := valueAsString(indexValue(ctx, 1))
	if _, ok := s.middlewareGroups[group]; !ok {
		return nil, data.NewErrorThrowByName(nil,
			fmt.Errorf("The [%s] middleware group has not been defined.", group),
			"InvalidArgumentException")
	}
	if mw != "" && !containsString(s.middlewareGroups[group], mw) {
		s.middlewareGroups[group] = append([]string{mw}, s.middlewareGroups[group]...)
		syncProperties(thisClassValue(ctx), s)
		if ctl := syncMiddlewareToRouter(s); ctl != nil {
			return nil, ctl
		}
	}
	return thisClassValue(ctx), nil
}

func kernelAppendToMiddlewareGroup(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	group := valueAsString(indexValue(ctx, 0))
	mw := valueAsString(indexValue(ctx, 1))
	if _, ok := s.middlewareGroups[group]; !ok {
		return nil, data.NewErrorThrowByName(nil,
			fmt.Errorf("The [%s] middleware group has not been defined.", group),
			"InvalidArgumentException")
	}
	if mw != "" && !containsString(s.middlewareGroups[group], mw) {
		s.middlewareGroups[group] = append(s.middlewareGroups[group], mw)
		syncProperties(thisClassValue(ctx), s)
		if ctl := syncMiddlewareToRouter(s); ctl != nil {
			return nil, ctl
		}
	}
	return thisClassValue(ctx), nil
}

func kernelSetMiddlewareGroups(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s.middlewareGroups = groupsFromValue(indexValue(ctx, 0))
	syncProperties(thisClassValue(ctx), s)
	if ctl := syncMiddlewareToRouter(s); ctl != nil {
		return nil, ctl
	}
	return thisClassValue(ctx), nil
}

func kernelAddToMiddlewarePriorityBefore(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	before := existingNames(indexValue(ctx, 0))
	mw := valueAsString(indexValue(ctx, 1))
	addToMiddlewarePriorityRelative(s, before, mw, false)
	syncProperties(thisClassValue(ctx), s)
	if ctl := syncMiddlewareToRouter(s); ctl != nil {
		return nil, ctl
	}
	return thisClassValue(ctx), nil
}

func kernelAddToMiddlewarePriorityAfter(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	after := existingNames(indexValue(ctx, 0))
	mw := valueAsString(indexValue(ctx, 1))
	addToMiddlewarePriorityRelative(s, after, mw, true)
	syncProperties(thisClassValue(ctx), s)
	if ctl := syncMiddlewareToRouter(s); ctl != nil {
		return nil, ctl
	}
	return thisClassValue(ctx), nil
}

func kernelPrependToMiddlewarePriority(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	mw := valueAsString(indexValue(ctx, 0))
	if mw != "" && !containsString(s.middlewarePriority, mw) {
		s.middlewarePriority = append([]string{mw}, s.middlewarePriority...)
		syncProperties(thisClassValue(ctx), s)
		if ctl := syncMiddlewareToRouter(s); ctl != nil {
			return nil, ctl
		}
	}
	return thisClassValue(ctx), nil
}

func kernelAppendToMiddlewarePriority(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	mw := valueAsString(indexValue(ctx, 0))
	if mw != "" && !containsString(s.middlewarePriority, mw) {
		s.middlewarePriority = append(s.middlewarePriority, mw)
		syncProperties(thisClassValue(ctx), s)
		if ctl := syncMiddlewareToRouter(s); ctl != nil {
			return nil, ctl
		}
	}
	return thisClassValue(ctx), nil
}

func kernelSetMiddlewarePriority(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s.middlewarePriority = stringListFromValue(indexValue(ctx, 0))
	syncProperties(thisClassValue(ctx), s)
	if ctl := syncMiddlewareToRouter(s); ctl != nil {
		return nil, ctl
	}
	return thisClassValue(ctx), nil
}

func kernelSetMiddlewareAliases(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := getState(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s.middlewareAliases = stringMapFromValue(indexValue(ctx, 0))
	syncProperties(thisClassValue(ctx), s)
	if ctl := syncMiddlewareToRouter(s); ctl != nil {
		return nil, ctl
	}
	return thisClassValue(ctx), nil
}

// addToMiddlewarePriorityRelative 对齐 Foundation\Http\Kernel::addToMiddlewarePriorityRelative。
func addToMiddlewarePriorityRelative(s *kernelState, existing []string, middleware string, after bool) {
	if s == nil || middleware == "" || containsString(s.middlewarePriority, middleware) {
		return
	}
	index := 0
	if !after {
		index = len(s.middlewarePriority)
	}
	for _, existingMiddleware := range existing {
		if containsString(s.middlewarePriority, existingMiddleware) {
			middlewareIndex := indexOfString(s.middlewarePriority, existingMiddleware)
			if after && middlewareIndex > index {
				index = middlewareIndex + 1
			} else if !after && middlewareIndex < index {
				index = middlewareIndex
			}
		}
	}
	switch {
	case index == 0 && !after:
		s.middlewarePriority = append([]string{middleware}, s.middlewarePriority...)
	case (after && index == 0) || index == len(s.middlewarePriority):
		s.middlewarePriority = append(s.middlewarePriority, middleware)
	default:
		s.middlewarePriority = append(s.middlewarePriority[:index], append([]string{middleware}, s.middlewarePriority[index:]...)...)
	}
}

func syncMiddlewareToRouter(s *kernelState) data.Control {
	if s == nil || s.router == nil {
		return nil
	}
	router, ok := s.router.(*data.ClassValue)
	if !ok || router == nil {
		return nil
	}
	_ = router.SetProperty("middlewarePriority", stringsToArrayValue(s.middlewarePriority))

	for key, middleware := range s.middlewareGroups {
		_, ctl := callObjectMethod(s.router, "middlewareGroup", data.NewStringValue(key), stringsToArrayValue(middleware))
		if ctl != nil {
			return ctl
		}
	}
	for key, middleware := range s.middlewareAliases {
		_, ctl := callObjectMethod(s.router, "aliasMiddleware", data.NewStringValue(key), data.NewStringValue(middleware))
		if ctl != nil {
			return ctl
		}
	}
	return nil
}

func asValue(v data.GetValue) data.Value {
	if v == nil {
		return nil
	}
	if val, ok := v.(data.Value); ok {
		return val
	}
	return nil
}

func methodExists(obj data.Value, name string) (data.Method, bool) {
	cv, ok := obj.(*data.ClassValue)
	if !ok || cv == nil {
		return nil, false
	}
	return cv.GetMethod(name)
}
