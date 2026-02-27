package preg

import (
	"github.com/php-any/origami/data"
)

// Load 注册所有 preg 相关的函数和常量
func Load(vm data.VM) {
	// 注册函数
	vm.AddFunc(NewPregMatchAllFunction())
	vm.AddFunc(NewPregSplitFunction())
	vm.AddFunc(NewPregReplaceFunction())
	vm.AddFunc(NewPregReplaceCallbackFunction())
	vm.AddFunc(NewPregGrepFunction())

	// 注册常量
	InitConstants(vm)
}
