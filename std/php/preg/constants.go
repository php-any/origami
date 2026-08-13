package preg

import (
	"github.com/php-any/origami/data"
)

// InitConstants 注册 preg 相关的 PHP 常量
func InitConstants(vm data.VM) {
	// preg_split 相关常量
	vm.SetConstant("PREG_SPLIT_NO_EMPTY", data.NewIntValue(1))
	vm.SetConstant("PREG_SPLIT_DELIM_CAPTURE", data.NewIntValue(2))
	vm.SetConstant("PREG_SPLIT_OFFSET_CAPTURE", data.NewIntValue(4))

	// preg_match_all 相关常量
	vm.SetConstant("PREG_PATTERN_ORDER", data.NewIntValue(1))
	vm.SetConstant("PREG_SET_ORDER", data.NewIntValue(2))
	vm.SetConstant("PREG_OFFSET_CAPTURE", data.NewIntValue(256))
	vm.SetConstant("PREG_UNMATCHED_AS_NULL", data.NewIntValue(512))

	// preg_grep 相关常量
	vm.SetConstant("PREG_GREP_INVERT", data.NewIntValue(1))

	// preg 错误码常量
	vm.SetConstant("PREG_NO_ERROR", data.NewIntValue(0))
	vm.SetConstant("PREG_INTERNAL_ERROR", data.NewIntValue(1))
	vm.SetConstant("PREG_BACKTRACK_LIMIT_ERROR", data.NewIntValue(2))
	vm.SetConstant("PREG_RECURSION_LIMIT_ERROR", data.NewIntValue(3))
	vm.SetConstant("PREG_BAD_UTF8_ERROR", data.NewIntValue(4))
	vm.SetConstant("PREG_BAD_UTF8_OFFSET_ERROR", data.NewIntValue(5))
	vm.SetConstant("PREG_JIT_STACKLIMIT_ERROR", data.NewIntValue(6))
}
