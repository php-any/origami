package runtime

import (
	"context"
	"os"
)

type uploadedFileState struct {
	paths  map[string]bool
	closed bool
	stop   func() bool
}

func (vm *RequestVM) RegisterUploadedFile(path string) bool {
	vm.mu.Lock()
	if vm.uploads == nil {
		vm.uploads = &uploadedFileState{paths: make(map[string]bool)}
		if request := vm.requestCall().deadline; request != nil && request.Done() != nil {
			vm.uploads.stop = context.AfterFunc(request, vm.cleanupUploadedFiles)
		}
	}
	if vm.uploads.closed {
		vm.mu.Unlock()
		_ = os.Remove(path)
		return false
	}
	vm.uploads.paths[path] = true
	vm.mu.Unlock()
	return true
}
func (vm *RequestVM) IsUploadedFile(path string) bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.uploads != nil && !vm.uploads.closed && vm.uploads.paths[path]
}
func (vm *RequestVM) ForgetUploadedFile(path string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.uploads != nil {
		delete(vm.uploads.paths, path)
	}
}
func (vm *RequestVM) cleanupUploadedFiles() {
	vm.mu.Lock()
	if vm.uploads == nil || vm.uploads.closed {
		vm.mu.Unlock()
		return
	}
	state := vm.uploads
	state.closed = true
	paths, stop := state.paths, state.stop
	state.paths, state.stop = nil, nil
	vm.mu.Unlock()
	if stop != nil {
		stop()
	}
	for path := range paths {
		_ = os.Remove(path)
	}
}
