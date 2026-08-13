//go:build !windows

package signal

import (
	"syscall"

	"github.com/php-any/origami/data"
)

func initPlatformSignalConstants(vm data.VM) {
	vm.SetConstant("SIGUSR1", data.NewIntValue(int(syscall.SIGUSR1)))
	vm.SetConstant("SIGUSR2", data.NewIntValue(int(syscall.SIGUSR2)))
}
