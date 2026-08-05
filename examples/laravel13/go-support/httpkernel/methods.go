package httpkernel

import (
	"fmt"
	"strconv"

	"github.com/php-any/origami/data"
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

	// 1) bootstrap
	if ctl := doBootstrap(ctx, s); ctl != nil {
		return nil, ctl
	}

	// 2) 绑定 request 到容器
	_, ctl = callObjectMethodInContext(ctx, s.app, "instance", data.NewStringValue("request"), request)
	if ctl != nil {
		return nil, ctl
	}

	// 3) 对齐 Foundation\Http\Kernel::handle：路由分发失败时 report + render，
	// 而不是把 NotFoundHttpException 等冒泡成 Fatal。
	response, ctl := dispatchToRouter(ctx, s, request)
	if ctl == nil {
		if response == nil {
			return data.NewNullValue(), nil
		}
		return response, nil
	}
	return renderException(ctx, s, request, ctl)
}

// dispatchToRouter 完成路由匹配与响应准备（不含全局中间件 Pipeline）。
func dispatchToRouter(ctx data.Context, s *kernelState, request data.Value) (data.GetValue, data.Control) {
	route, ctl := callObjectMethodInContext(ctx, s.router, "findRoute", request)
	if ctl != nil {
		return nil, ctl
	}
	result, ctl := callObjectMethodInContext(ctx, asValue(route), "run")
	if ctl != nil {
		return nil, ctl
	}
	return callObjectMethodInContext(ctx, s.router, "prepareResponse", request, asValue(result))
}

const fqnExceptionHandler = "Illuminate\\Contracts\\Debug\\ExceptionHandler"

// renderException 对齐 Foundation\Http\Kernel::reportException / renderException。
func renderException(ctx data.Context, s *kernelState, request data.Value, thrown data.Control) (data.GetValue, data.Control) {
	exception, ok := exceptionFromControl(thrown)
	if !ok || s.app == nil {
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

// fallbackExceptionResponse 在 Exception Handler 不可用时，尽量按 HttpException 状态码返回。
func fallbackExceptionResponse(ctx data.Context, exception data.Value, original data.Control) (data.GetValue, data.Control) {
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
	// 最小实现：尝试调用 app->terminate()；失败则原样返回 control。
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
