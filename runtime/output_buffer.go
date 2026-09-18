package runtime

import (
	"github.com/php-any/origami/data"
)

func (vm *VM) activeOut() *OutputState {
	if st := currentRequestOutput(); st != nil {
		return st
	}
	return vm.out
}

func (vm *VM) WriteOutput(s string) {
	vm.activeOut().write(s, data.WriteOutput)
}

func (vm *VM) StartOutputBuffer() {
	vm.activeOut().start()
}

func (vm *VM) StartOutputBufferSpec(spec data.OutputBufferStartSpec) bool {
	return vm.activeOut().startSpec(spec)
}

func (vm *VM) CleanOutputBuffer() (string, bool) {
	return vm.activeOut().clean()
}

func (vm *VM) FlushOutputBuffer() (string, bool) {
	return vm.activeOut().flushWithFallback(data.WriteOutput)
}

func (vm *VM) FlushCurrentBuffer() (string, bool) {
	return vm.activeOut().flushCurrent(data.PHPOutputHandlerFlush, data.WriteOutput)
}

func (vm *VM) CleanCurrentBuffer() bool {
	return vm.activeOut().cleanCurrent()
}

func (vm *VM) OutputBufferLength() (int, bool) {
	return vm.activeOut().length()
}

func (vm *VM) OutputBufferStatus(full bool) []data.OutputBufferStatusInfo {
	return vm.activeOut().status(full)
}

func (vm *VM) ListOutputHandlers() []string {
	return vm.activeOut().handlers()
}

func (vm *VM) SetImplicitFlush(on bool) {
	vm.activeOut().setImplicitFlush(on)
}

func (vm *VM) IsImplicitFlush() bool {
	return vm.activeOut().isImplicitFlush()
}

func (vm *VM) TakeOutputControl() data.Control {
	return vm.activeOut().takeControl()
}

func (vm *VM) FlushSAPI() {
	vm.activeOut().flushSAPI()
}

func (vm *VM) OutputBufferContents() (string, bool) {
	return vm.activeOut().contents()
}

func (vm *VM) OutputBufferLevel() int {
	return vm.activeOut().level()
}

func (vm *TempVM) activeOut() *OutputState {
	if st := currentRequestOutput(); st != nil {
		return st
	}
	return vm.out
}

func (vm *TempVM) fallbackFor(st *OutputState) func(string) {
	if st == vm.out {
		return vm.Base.WriteOutput
	}
	return data.WriteOutput
}

func (vm *TempVM) WriteOutput(s string) {
	st := vm.activeOut()
	st.write(s, vm.fallbackFor(st))
}

func (vm *TempVM) StartOutputBuffer() {
	vm.activeOut().start()
}

func (vm *TempVM) StartOutputBufferSpec(spec data.OutputBufferStartSpec) bool {
	return vm.activeOut().startSpec(spec)
}

func (vm *TempVM) CleanOutputBuffer() (string, bool) {
	return vm.activeOut().clean()
}

func (vm *TempVM) FlushOutputBuffer() (string, bool) {
	st := vm.activeOut()
	return st.flushWithFallback(vm.fallbackFor(st))
}

func (vm *TempVM) FlushCurrentBuffer() (string, bool) {
	st := vm.activeOut()
	return st.flushCurrent(data.PHPOutputHandlerFlush, vm.fallbackFor(st))
}

func (vm *TempVM) CleanCurrentBuffer() bool {
	return vm.activeOut().cleanCurrent()
}

func (vm *TempVM) OutputBufferLength() (int, bool) {
	return vm.activeOut().length()
}

func (vm *TempVM) OutputBufferStatus(full bool) []data.OutputBufferStatusInfo {
	return vm.activeOut().status(full)
}

func (vm *TempVM) ListOutputHandlers() []string {
	return vm.activeOut().handlers()
}

func (vm *TempVM) SetImplicitFlush(on bool) {
	vm.activeOut().setImplicitFlush(on)
}

func (vm *TempVM) IsImplicitFlush() bool {
	return vm.activeOut().isImplicitFlush()
}

func (vm *TempVM) TakeOutputControl() data.Control {
	return vm.activeOut().takeControl()
}

func (vm *TempVM) FlushSAPI() {
	vm.activeOut().flushSAPI()
}

func (vm *TempVM) OutputBufferContents() (string, bool) {
	return vm.activeOut().contents()
}

func (vm *TempVM) OutputBufferLevel() int {
	return vm.activeOut().level()
}

var (
	_ data.OutputSink       = (*VM)(nil)
	_ data.OutputBufferHost = (*VM)(nil)
	_ data.OutputSink       = (*TempVM)(nil)
	_ data.OutputBufferHost = (*TempVM)(nil)
	_ data.OutputSink       = (*Context)(nil)
	_ data.OutputBufferHost = (*Context)(nil)
)
