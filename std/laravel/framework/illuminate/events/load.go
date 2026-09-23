package events

import "github.com/php-any/origami/data"

// Load 注册 illuminate/events。
func Load(vm data.VM) {
	vm.AddClass(NewDispatcherClass())
}
