package php

import "github.com/php-any/origami/data"

// Pathinfo 常量
const (
	PATHINFO_DIRNAME   = 1
	PATHINFO_BASENAME  = 2
	PATHINFO_EXTENSION = 4
	PATHINFO_FILENAME  = 8
)

// InitPathinfoConstants 初始化 pathinfo 相关常量
func InitPathinfoConstants(vm data.VM) {
	vm.SetConstant("PATHINFO_DIRNAME", data.NewIntValue(PATHINFO_DIRNAME))
	vm.SetConstant("PATHINFO_BASENAME", data.NewIntValue(PATHINFO_BASENAME))
	vm.SetConstant("PATHINFO_EXTENSION", data.NewIntValue(PATHINFO_EXTENSION))
	vm.SetConstant("PATHINFO_FILENAME", data.NewIntValue(PATHINFO_FILENAME))
}
