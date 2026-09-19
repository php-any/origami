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
		app:                cloneIfClass(src.app, ctx),
		router:             cloneIfClass(src.router, ctx),
		bootstrappers:      append([]string(nil), src.bootstrappers...),
		middleware:         append([]string(nil), src.middleware...),
		middlewareGroups:   cloneGroupsMap(src.middlewareGroups),
		middlewareAliases:  cloneStringMap(src.middlewareAliases),
		middlewarePriority: append([]string(nil), src.middlewarePriority...),
		bootstrapped:       src.bootstrapped,
	}
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

// cloneRequestServices 把每请求会往里 append 的共享单例换成 clone。
// PHP clone：数组属性按值拷贝。否则 Dispatcher::$listeners / View composers
// 跨请求膨胀，异常页 VarDumper 能把进程拖到十几秒甚至 OOM。
func cloneRequestServices(ctx data.Context, app data.Value) {
	if app == nil {
		return
	}
	cloneAndBind(ctx, app, "events",
		"Illuminate\\Events\\Dispatcher",
		"Illuminate\\Contracts\\Events\\Dispatcher",
	)
	cloneAndBind(ctx, app, "view",
		"Illuminate\\View\\Factory",
		"Illuminate\\Contracts\\View\\Factory",
	)
	cloneAndBind(ctx, app, "url",
		"Illuminate\\Routing\\UrlGenerator",
		"Illuminate\\Contracts\\Routing\\UrlGenerator",
	)
}

func cloneAndBind(ctx data.Context, app data.Value, abstract string, aliases ...string) {
	raw, ctl := callObjectMethodInContext(ctx, app, "make", data.NewStringValue(abstract))
	if ctl != nil {
		return
	}
	cv, ok := asValue(raw).(*data.ClassValue)
	if !ok || cv == nil {
		return
	}
	cloned := cv.CloneSandbox(ctx)
	_, _ = callObjectMethodInContext(ctx, cloned, "setContainer", app)
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
	_, _ = callObjectMethodInContext(ctx, app, "instance", data.NewStringValue(abstract), cloned)
	for _, alias := range aliases {
		_, _ = callObjectMethodInContext(ctx, app, "instance", data.NewStringValue(alias), cloned)
	}
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
