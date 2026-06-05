package http

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	// 添加接口
	vm.AddInterface(NewMiddlewareInterface())
	// 添加类
	vm.AddClass(NewServerClass())
	vm.AddClass(NewHandlerClass())
	vm.AddClass(NewCookieClass())
	vm.AddClass(NewResponseWriterClass())
	vm.AddClass(NewRequestClass())
	// 添加函数
	vm.AddFunc(NewAppFunction())
	vm.AddFunc(NewAppFlashFunction())
}
