package context

import (
	"github.com/php-any/origami/data"
)

func Load(vm data.VM) {
	// 添加顶级函数
	for _, fun := range []data.FuncStmt{
		NewBackgroundFunction(),
		NewWithCancelFunction(),
		NewWithTimeoutFunction(),
		NewWithValueFunction(),
	} {
		vm.AddFunc(fun)
	}

}
