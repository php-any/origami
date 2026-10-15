package httpkernel

import "github.com/php-any/origami/data"

// Load 注册 Go 原生 Illuminate\Foundation\Http\Kernel，替换 Composer 的官方实现。
//
// 请求生命周期由本包承担：__construct / bootstrap / handle / terminate 以及
// ApplicationBuilder 在 afterResolving 里调用的全部 setter（setGlobalMiddleware /
// setMiddlewareGroups / setMiddlewareAliases / setMiddlewarePriority /
// addToMiddlewarePriority{Before,After}）都已注册，handle 后的
// sendRequestThroughRouter 语义（bootstrap → 全局中间件 Pipeline → Router::dispatch
// → RequestHandled → report/renderException）在 kernelHandle/dispatchToRouter 里逐条对齐。
func Load(vm data.VM) { vm.AddClass(NewClass()) }
