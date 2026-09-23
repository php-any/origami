package view

import "github.com/php-any/origami/data"

// Load 注册 illuminate/view 加速件。
func Load(vm data.VM) {
	vm.AddClass(NewAppendableClass())
	vm.AddClass(NewAttributeBagClass())
}
