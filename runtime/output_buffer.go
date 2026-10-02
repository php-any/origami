package runtime

import (
	"github.com/php-any/origami/data"
)

func (vm *VM) HeaderCallbackState() *data.HeaderCallbackState        { return &vm.activeOut().headers }
func (vm *RequestVM) HeaderCallbackState() *data.HeaderCallbackState { return &vm.activeOut().headers }
func (c *Context) HeaderCallbackState() *data.HeaderCallbackState    { return &c.resolveOut().headers }

func (st *OutputState) bindTarget(write func(string) data.Control, flush func()) func() {
	previousSink, previousFlush := st.sink, st.sapiFlush
	st.sink = func(s string) {
		if ctl := write(s); ctl != nil {
			st.setControl(ctl)
		}
	}
	st.sapiFlush = flush
	return func() { st.sink, st.sapiFlush = previousSink, previousFlush }
}

func (vm *VM) BindOutputTarget(write func(string) data.Control, flush func()) func() {
	return vm.activeOut().bindTarget(write, flush)
}

func (vm *RequestVM) BindOutputTarget(write func(string) data.Control, flush func()) func() {
	return vm.activeOut().bindTarget(write, flush)
}

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

func (vm *RequestVM) activeOut() *OutputState {
	if vm.out.local {
		return vm.out
	}
	if st := currentRequestOutput(); st != nil {
		return st
	}
	return vm.out
}

func (vm *RequestVM) fallbackFor(st *OutputState) func(string) {
	if st == vm.out {
		return vm.Base.WriteOutput
	}
	return data.WriteOutput
}

func (vm *RequestVM) WriteOutput(s string) {
	st := vm.activeOut()
	st.write(s, vm.fallbackFor(st))
}

func (vm *RequestVM) StartOutputBuffer() {
	vm.activeOut().start()
}

func (vm *RequestVM) StartOutputBufferSpec(spec data.OutputBufferStartSpec) bool {
	return vm.activeOut().startSpec(spec)
}

func (vm *RequestVM) CleanOutputBuffer() (string, bool) {
	return vm.activeOut().clean()
}

func (vm *RequestVM) FlushOutputBuffer() (string, bool) {
	st := vm.activeOut()
	return st.flushWithFallback(vm.fallbackFor(st))
}

func (vm *RequestVM) FlushCurrentBuffer() (string, bool) {
	st := vm.activeOut()
	return st.flushCurrent(data.PHPOutputHandlerFlush, vm.fallbackFor(st))
}

func (vm *RequestVM) CleanCurrentBuffer() bool {
	return vm.activeOut().cleanCurrent()
}

func (vm *RequestVM) OutputBufferLength() (int, bool) {
	return vm.activeOut().length()
}

func (vm *RequestVM) OutputBufferStatus(full bool) []data.OutputBufferStatusInfo {
	return vm.activeOut().status(full)
}

func (vm *RequestVM) ListOutputHandlers() []string {
	return vm.activeOut().handlers()
}

func (vm *RequestVM) SetImplicitFlush(on bool) {
	vm.activeOut().setImplicitFlush(on)
}

func (vm *RequestVM) IsImplicitFlush() bool {
	return vm.activeOut().isImplicitFlush()
}

func (vm *RequestVM) TakeOutputControl() data.Control {
	return vm.activeOut().takeControl()
}

func (vm *RequestVM) FlushSAPI() {
	vm.activeOut().flushSAPI()
}

func (vm *RequestVM) OutputBufferContents() (string, bool) {
	return vm.activeOut().contents()
}

func (vm *RequestVM) OutputBufferLevel() int {
	return vm.activeOut().level()
}

var (
	_ data.OutputSink       = (*VM)(nil)
	_ data.OutputBufferHost = (*VM)(nil)
	_ data.OutputSink       = (*RequestVM)(nil)
	_ data.OutputBufferHost = (*RequestVM)(nil)
	_ data.OutputSink       = (*Context)(nil)
	_ data.OutputBufferHost = (*Context)(nil)
)
