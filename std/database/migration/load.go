package migration

import "github.com/php-any/origami/data"

// Load 加载数据库迁移模块
func Load(vm data.VM) {
	vm.AddClass(NewMigrateResultClass())
	vm.AddFunc(NewMigrateFunction())
}
