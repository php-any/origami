package support

import "github.com/php-any/origami/data"

func Load(vm data.VM) {
	registerHelpers(vm)
	vm.AddClass(NewStrClass())
	vm.AddClass(NewHtmlStringClass())
	vm.AddClass(NewStringableClass())
}
