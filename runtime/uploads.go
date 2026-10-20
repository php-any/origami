package runtime

import "os"

type uploadedFileState struct {
	paths  map[string]bool
	closed bool
}

func (vm *RequestVM) RegisterUploadedFile(path string) bool {
	if vm.uploads == nil {
		vm.uploads = &uploadedFileState{paths: make(map[string]bool)}
	}
	if vm.uploads.closed {
		_ = os.Remove(path)
		return false
	}
	vm.uploads.paths[path] = true
	return true
}
func (vm *RequestVM) IsUploadedFile(path string) bool {
	return vm.uploads != nil && !vm.uploads.closed && vm.uploads.paths[path]
}
func (vm *RequestVM) ForgetUploadedFile(path string) {
	if vm.uploads != nil {
		delete(vm.uploads.paths, path)
	}
}
func (vm *RequestVM) cleanupUploadedFiles() {
	if vm.uploads == nil || vm.uploads.closed {
		return
	}
	state := vm.uploads
	state.closed = true
	paths := state.paths
	state.paths = nil
	for path := range paths {
		_ = os.Remove(path)
	}
}
