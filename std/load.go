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
		NewIntFunction(),
		NewStringFunction(),
		NewBoolFunction(),
		NewFloatFunction(),
		NewObjectFunction(),
	} {
		vm.AddFunc(fun)
	}

	vm.AddClass(log.NewLogClass())
	// 注册 Throwable / Stringable / JsonSerializable 接口与 Exception 类
	vm.AddInterface(exception.NewThrowableInterface())
	vm.AddInterface(exception.NewStringableInterface())
	vm.AddInterface(exception.NewJsonSerializableInterface())
	vm.AddClass(exception.NewExceptionClass())
	vm.AddClass(exception.NewReflectionExceptionClass())
	vm.AddClass(os.NewOSClass())
	reflect.Load(vm)
	channel.Load(vm)
	loop.Load(vm)
	database.Load(vm)
}
