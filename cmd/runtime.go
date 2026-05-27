package cmd

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

var runtimeLoader func(vm data.VM)

// SetRuntimeLoader 注册标准库加载闭包，由 main 包注入各 Load 调用。
func SetRuntimeLoader(fn func(vm data.VM)) {
	runtimeLoader = fn
}

func getRuntimeVM() (*runtime.VM, *parser.Parser) {
	if runtimeLoader == nil {
		panic("runtime loader not set, call cmd.SetRuntimeLoader from main")
	}
	p := parser.NewParser()
	vm := runtime.NewVM(p)
	runtimeLoader(vm)
	return vm.(*runtime.VM), p
}
