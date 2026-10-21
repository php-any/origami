package runtime

import "github.com/php-any/origami/internal/pharfs"

func (vm *VM) CurrentPHPFile() string        { return vm.localCall().CurrentFile }
func (vm *RequestVM) CurrentPHPFile() string { return vm.call.CurrentFile }

func enterPHPFile(state *CallState, file string) func() {
	previous := state.CurrentFile
	state.CurrentFile = file
	return func() { state.CurrentFile = previous; releaseAutoCallState(state) }
}

func (vm *VM) PharFS() *pharfs.FS {
	return vm.pharFiles
}
func (vm *RequestVM) PharFS() *pharfs.FS {
	return vm.pharFiles
}
