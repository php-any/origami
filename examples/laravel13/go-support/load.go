package gosupport

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/examples/laravel13/go-support/httpfoundation"
	"github.com/php-any/origami/examples/laravel13/go-support/httpkernel"
)

// Load 注册 Laravel 13 示例所需的运行时扩展。
//
// 约定：
//  1. 通用 PHP 语义缺口优先修 Origami 核心（std/ / node/ / data/ / runtime/）。
//  2. 本包只放「Laravel 应用侧」或短期无法进核心的能力：
//     - 覆盖/补充个别内置函数
//     - 实现 Laravel Request 捕获、HTTP Kernel 桥接等适配层
//  3. 禁止在 PHP bootstrap 里用大段匿名 singleton 冒充框架服务；
//     适配逻辑集中在这里或正当的框架层覆盖。
func Load(vm data.VM) {
	vm.AddClass(httpfoundation.NewParameterBagClass())
	vm.AddClass(httpfoundation.NewHeaderBagClass())
	vm.AddClass(httpfoundation.NewInputBagClass())
	vm.AddClass(httpfoundation.NewServerBagClass())
	vm.AddClass(httpfoundation.NewFileBagClass())
	vm.AddClass(httpfoundation.NewResponseHeaderBagClass())
	vm.AddClass(httpfoundation.NewSymfonyRequestClass())
	vm.AddClass(httpfoundation.NewSymfonyResponseClass())
	vm.AddClass(httpfoundation.NewIlluminateRequestClass())
	vm.AddClass(httpfoundation.NewIlluminateResponseClass())
	vm.AddClass(httpkernel.NewClass())
	vm.AddClass(NewServeCommandClass())
	registerBuiltins(vm)
}

func registerBuiltins(vm data.VM) {
	// 随冒烟失败项逐步追加，例如 password_*、debug_backtrace、Request::capture 相关等。
}
