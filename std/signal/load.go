package signal

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

func Load(vm data.VM) {
	runtime.InstallShutdownSignalHandler(vm)

	vm.AddClass(NewSignalChannelClass())

	for _, fun := range []data.FuncStmt{
		NewNotifyFunction(),
		NewStopFunction(),
		NewResetFunction(),
		NewIgnoreFunction(),
		NewWaitFunction(),
	} {
		vm.AddFunc(fun)
	}

	initConstants(vm)
}
