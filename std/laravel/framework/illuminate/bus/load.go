package bus

import "github.com/php-any/origami/data"

// Load：Bus Dispatcher 未覆盖 queue 链前不占名。
func Load(vm data.VM) {
	_ = vm
	_ = NewDispatcherClass
}
