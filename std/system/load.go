package system

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	vm.AddInterface(newDateTimeInterface())
	vm.AddClass(&PhpDateTimeClass{})  // 全局 DateTime
	vm.AddClass(&DateTimeZoneClass{}) // 全局 DateTimeZone
}
