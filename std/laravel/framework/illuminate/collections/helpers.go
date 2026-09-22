package collections

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

func kitHelperCollect() data.FuncStmt {
	return kit.HelperFunc("collect", []string{"value"}, func(ctx data.Context) (data.GetValue, data.Control) {
		v, _ := ctx.GetIndexValue(0)
		return newCollectionInstance(ctx, v)
	})
}

func kitHelperDataGet() data.FuncStmt {
	return kit.HelperFunc("data_get", []string{"target", "key", "default"}, arrGet)
}

func kitHelperDataSet() data.FuncStmt {
	return kit.HelperFuncRef("data_set", []string{"target", "key", "value", "overwrite"}, arrSet)
}
