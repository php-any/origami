package std

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/channel"
	"github.com/php-any/origami/std/database"
	"github.com/php-any/origami/std/exception"
	"github.com/php-any/origami/std/log"
	"github.com/php-any/origami/std/loop"
	"github.com/php-any/origami/std/reflect"
	"github.com/php-any/origami/std/system/os"
)

func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewDumpFunction(),
		NewIncludeFunction(),
		NewIntFunction(),
		NewStringFunction(),
		NewBoolFunction(),
		NewFloatFunction(),
	} {
		vm.AddFunc(fun)
	}

	vm.AddClass(log.NewLogClass())
	vm.AddClass(exception.NewExceptionClass())
	vm.AddClass(os.NewOSClass())
	reflect.Load(vm)
	channel.Load(vm)
	loop.Load(vm)
	database.Load(vm)
}
