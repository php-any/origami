package clock

import "github.com/php-any/origami/data"

// Load 注册 NativeClock。不注册 Clock 静态门面，让 PHP 继续加载 vendor 实现。
func Load(vm data.VM) {
	vm.AddClass(NewNativeClockClass())
}
