//go:build windows

package signal

import "github.com/php-any/origami/data"

func initPlatformSignalConstants(vm data.VM) {
	// Windows 无 SIGUSR1/SIGUSR2；占位常量供脚本引用。
	vm.SetConstant("SIGUSR1", data.NewIntValue(10))
	vm.SetConstant("SIGUSR2", data.NewIntValue(12))
}
