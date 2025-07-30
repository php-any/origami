package system

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	vm.AddClass(&DateTimeClass{})
}
