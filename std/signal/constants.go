package signal

import (
	"syscall"

	"github.com/php-any/origami/data"
)

func initConstants(vm data.VM) {
	vm.SetConstant("SIGINT", data.NewIntValue(int(syscall.SIGINT)))
	vm.SetConstant("SIGTERM", data.NewIntValue(int(syscall.SIGTERM)))
	vm.SetConstant("SIGHUP", data.NewIntValue(int(syscall.SIGHUP)))
	vm.SetConstant("SIGQUIT", data.NewIntValue(int(syscall.SIGQUIT)))
	vm.SetConstant("SIGABRT", data.NewIntValue(int(syscall.SIGABRT)))
	vm.SetConstant("SIGKILL", data.NewIntValue(int(syscall.SIGKILL)))
	initPlatformSignalConstants(vm)
}
