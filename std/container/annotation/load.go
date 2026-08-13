package annotation

import "github.com/php-any/origami/data"

// Load 加载容器注解模块
func Load(vm data.VM) {
	vm.AddClass(NewComponentClass())
	vm.AddClass(NewSingletonAnnotationClass())
	vm.AddClass(NewScopedAnnotationClass())
	vm.AddClass(NewBindClass())
	vm.AddClass(NewInjectClass())
	vm.AddClass(NewNamedClass())
}
