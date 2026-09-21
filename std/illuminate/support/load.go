package support

import "github.com/php-any/origami/data"

// Load 把 Illuminate 热包以 Go 实现提前注册进 VM。
// class_exists / GetOrLoadClass 命中后不再解析 vendor 里对应的 PHP 文件。
func Load(vm data.VM) {
	vm.AddClass(NewArrClass())
	vm.AddClass(NewStrClass())
	vm.AddClass(NewStringableClass())
	vm.AddClass(NewHtmlStringClass())
	vm.AddClass(NewAttributeBagClass())
	vm.AddClass(NewAppendableClass())
	registerHelpers(vm)
}
