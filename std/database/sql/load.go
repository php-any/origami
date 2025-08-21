package sql

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	// 添加顶级函数
	for _, fun := range []data.FuncStmt{
		NewOpenFunction(),
	} {
		vm.AddFunc(fun)
	}

	// 添加类
	vm.AddClass(NewConnClass())
	vm.AddClass(NewDBClass())
	vm.AddClass(NewRowClass())
	vm.AddClass(NewRowsClass())
	vm.AddClass(NewStmtClass())
	vm.AddClass(NewTxClass())
	vm.AddClass(NewTxOptionsClass())
}
