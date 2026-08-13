package php

import (
	"github.com/php-any/origami/data"
)

// InitJsonConstants 注册 json_encode / json_decode 相关 PHP 常量
func InitJsonConstants(vm data.VM) {
	// 错误码
	vm.SetConstant("JSON_ERROR_NONE", data.NewIntValue(JSON_ERROR_NONE))
	vm.SetConstant("JSON_ERROR_DEPTH", data.NewIntValue(JSON_ERROR_DEPTH))
	vm.SetConstant("JSON_ERROR_STATE_MISMATCH", data.NewIntValue(JSON_ERROR_STATE_MISMATCH))
	vm.SetConstant("JSON_ERROR_CTRL_CHAR", data.NewIntValue(JSON_ERROR_CTRL_CHAR))
	vm.SetConstant("JSON_ERROR_SYNTAX", data.NewIntValue(JSON_ERROR_SYNTAX))
	vm.SetConstant("JSON_ERROR_UTF8", data.NewIntValue(JSON_ERROR_UTF8))
	vm.SetConstant("JSON_ERROR_RECURSION", data.NewIntValue(JSON_ERROR_RECURSION))
	vm.SetConstant("JSON_ERROR_INF_OR_NAN", data.NewIntValue(JSON_ERROR_INF_OR_NAN))
	vm.SetConstant("JSON_ERROR_UNSUPPORTED_TYPE", data.NewIntValue(JSON_ERROR_UNSUPPORTED_TYPE))
	vm.SetConstant("JSON_ERROR_INVALID_PROPERTY_NAME", data.NewIntValue(JSON_ERROR_INVALID_PROPERTY_NAME))
	vm.SetConstant("JSON_ERROR_UTF16", data.NewIntValue(JSON_ERROR_UTF16))

	// json_encode / json_decode 选项（与 PHP 一致）
	vm.SetConstant("JSON_HEX_TAG", data.NewIntValue(1))
	vm.SetConstant("JSON_HEX_AMP", data.NewIntValue(2))
	vm.SetConstant("JSON_HEX_APOS", data.NewIntValue(4))
	vm.SetConstant("JSON_HEX_QUOT", data.NewIntValue(8))
	vm.SetConstant("JSON_FORCE_OBJECT", data.NewIntValue(16))
	vm.SetConstant("JSON_NUMERIC_CHECK", data.NewIntValue(32))
	vm.SetConstant("JSON_UNESCAPED_SLASHES", data.NewIntValue(64))
	vm.SetConstant("JSON_PRETTY_PRINT", data.NewIntValue(128))
	vm.SetConstant("JSON_UNESCAPED_UNICODE", data.NewIntValue(256))
	vm.SetConstant("JSON_PARTIAL_OUTPUT_ON_ERROR", data.NewIntValue(512))
	vm.SetConstant("JSON_PRESERVE_ZERO_FRACTION", data.NewIntValue(1024))
	vm.SetConstant("JSON_UNESCAPED_LINE_TERMINATORS", data.NewIntValue(2048))
	vm.SetConstant("JSON_OBJECT_AS_ARRAY", data.NewIntValue(1))
	vm.SetConstant("JSON_BIGINT_AS_STRING", data.NewIntValue(2))
	vm.SetConstant("JSON_INVALID_UTF8_IGNORE", data.NewIntValue(1048576))
	vm.SetConstant("JSON_INVALID_UTF8_SUBSTITUTE", data.NewIntValue(2097152))
	vm.SetConstant("JSON_THROW_ON_ERROR", data.NewIntValue(4194304))
}
