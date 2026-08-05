package runtime

import (
	"strings"

	"github.com/php-any/origami/data"
)

func (vm *VM) WriteOutput(s string) {
	vm.mu.Lock()
	if n := len(vm.outputBuffers); n > 0 {
		vm.outputBuffers[n-1].WriteString(s)
		vm.mu.Unlock()
		return
	}
	vm.mu.Unlock()
	data.WriteOutput(s)
}

func (vm *VM) StartOutputBuffer() {
	vm.mu.Lock()
	vm.outputBuffers = append(vm.outputBuffers, &strings.Builder{})
	vm.mu.Unlock()
}

func (vm *VM) CleanOutputBuffer() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	content := vm.outputBuffers[n-1].String()
	vm.outputBuffers = vm.outputBuffers[:n-1]
	return content, true
}

func (vm *VM) OutputBufferContents() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	return vm.outputBuffers[n-1].String(), true
}

func (vm *VM) OutputBufferLevel() int {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return len(vm.outputBuffers)
}

func (vm *TempVM) WriteOutput(s string) {
	vm.mu.Lock()
	if n := len(vm.outputBuffers); n > 0 {
		vm.outputBuffers[n-1].WriteString(s)
		vm.mu.Unlock()
		return
	}
	vm.mu.Unlock()
	vm.Base.WriteOutput(s)
}

func (vm *TempVM) StartOutputBuffer() {
	vm.mu.Lock()
	vm.outputBuffers = append(vm.outputBuffers, &strings.Builder{})
	vm.mu.Unlock()
}

func (vm *TempVM) CleanOutputBuffer() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	content := vm.outputBuffers[n-1].String()
	vm.outputBuffers = vm.outputBuffers[:n-1]
	return content, true
}

func (vm *TempVM) OutputBufferContents() (string, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.outputBuffers)
	if n == 0 {
		return "", false
	}
	return vm.outputBuffers[n-1].String(), true
}

func (vm *TempVM) OutputBufferLevel() int {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return len(vm.outputBuffers)
}

var (
	_ data.OutputSink       = (*VM)(nil)
	_ data.OutputBufferHost = (*VM)(nil)
	_ data.OutputSink       = (*TempVM)(nil)
	_ data.OutputBufferHost = (*TempVM)(nil)
)
