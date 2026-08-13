package loop

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	vm.AddClass(NewListClass())
	vm.AddClass(NewHashMapClass())
}
