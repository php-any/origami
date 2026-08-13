package annotation

import (
	"github.com/php-any/origami/data"
)

// Load 加载数据库注解模块
func Load(vm data.VM) {
	// 注册数据库注解类
	vm.AddClass(NewTableClass())
	vm.AddClass(NewColumnClass())
	vm.AddClass(NewIdClass())
	vm.AddClass(NewGeneratedValueClass())
}
