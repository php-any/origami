package annotation

import "github.com/php-any/origami/data"

// Load 注册校验约束注解。
func Load(vm data.VM) {
	for _, spec := range constraintSpecs {
		vm.AddClass(newConstraintClass(spec))
	}
}
