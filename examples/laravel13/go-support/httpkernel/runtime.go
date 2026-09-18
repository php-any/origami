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
	return cv
}

func bindSandboxContainer(ctx data.Context, app, router data.Value) {
	if app == nil {
		return
	}
	_, _ = callObjectMethodInContext(ctx, app, "setInstance", app)
	if router != nil {
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
