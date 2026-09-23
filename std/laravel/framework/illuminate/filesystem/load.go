package filesystem

import "github.com/php-any/origami/data"

// Load 注册 illuminate/filesystem。
func Load(vm data.VM) {
	vm.AddClass(NewFilesystemClass())
}
