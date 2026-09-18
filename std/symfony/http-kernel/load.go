package httpkernel

import "github.com/php-any/origami/data"

// Load 暂不注册不完整 HttpKernel 事件类，避免挡住 vendor。
func Load(vm data.VM) {
	_ = vm
}
