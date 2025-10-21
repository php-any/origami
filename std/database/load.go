package database

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/database/annotation"
	"github.com/php-any/origami/std/database/sql"
)

func Load(vm data.VM) {
	vm.AddClass(NewDBClass())
	sql.Load(vm)

	// 注册数据库注解类
	annotation.Load(vm)

	// 注册数据库连接管理函数到脚本域
	for _, fun := range []data.FuncStmt{
		NewRegisterConnectionFunction(),
		NewRegisterDefaultConnectionFunction(),
		NewGetConnectionFunction(),
		NewGetDefaultConnectionFunction(),
		NewRemoveConnectionFunction(),
		NewListConnectionsFunction(),
	} {
		vm.AddFunc(fun)
	}
}
