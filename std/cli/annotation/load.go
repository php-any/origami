package annotation

import (
	"github.com/php-any/origami/data"
)

// Load 加载 CLI 注解模块
func Load(vm data.VM) {
	// 注册 CLI 注解类
	vm.AddClass(&CliApplicationClass{})
	vm.AddClass(&CommandClass{})
}
