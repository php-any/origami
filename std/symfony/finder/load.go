package finder

import "github.com/php-any/origami/data"

// Load 注册 Symfony Finder / SplFileInfo 原生实现。
func Load(vm data.VM) {
	vm.AddClass(NewFinderClass())
	vm.AddClass(NewSplFileInfoClass())
}
