package validation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/validation/annotation"
)

// Load 注册校验注解模块。
func Load(vm data.VM) {
	annotation.Load(vm)
	vm.AddClass(&ValidatorClass{})
}
