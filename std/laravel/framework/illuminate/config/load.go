package config

import "github.com/php-any/origami/data"

// Load 注册 illuminate/config。
func Load(vm data.VM) {
	vm.AddClass(NewRepositoryClass())
}
