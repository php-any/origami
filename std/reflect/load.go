package reflect

import (
	"github.com/php-any/origami/data"
)

// Load 加载反射模块到VM
func Load(vm data.VM) {
	vm.AddClass(&ReflectClass{})
}
