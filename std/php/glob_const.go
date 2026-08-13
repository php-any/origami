package php

import "github.com/php-any/origami/data"

// glob() 标志常量，与 PHP 一致。
const (
	GLOB_ERR      = 1
	GLOB_MARK     = 2
	GLOB_NOSORT   = 4
	GLOB_NOCHECK  = 16
	GLOB_NOESCAPE = 64
	GLOB_BRACE    = 1024
	GLOB_ONLYDIR  = 1073741824
)

// InitGlobConstants 注册 glob 相关常量。
func InitGlobConstants(vm data.VM) {
	vm.SetConstant("GLOB_ERR", data.NewIntValue(GLOB_ERR))
	vm.SetConstant("GLOB_MARK", data.NewIntValue(GLOB_MARK))
	vm.SetConstant("GLOB_NOSORT", data.NewIntValue(GLOB_NOSORT))
	vm.SetConstant("GLOB_NOCHECK", data.NewIntValue(GLOB_NOCHECK))
	vm.SetConstant("GLOB_NOESCAPE", data.NewIntValue(GLOB_NOESCAPE))
	vm.SetConstant("GLOB_BRACE", data.NewIntValue(GLOB_BRACE))
	vm.SetConstant("GLOB_ONLYDIR", data.NewIntValue(GLOB_ONLYDIR))
}
