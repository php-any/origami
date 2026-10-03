package runtime

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"strings"
)

type autoloadRegistration struct{ original, resolved data.Value }

func (vm *VM) AutoloadExtensions() string {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.autoloadExtensions != nil {
		return *vm.autoloadExtensions
	}
	return ".inc,.php"
}
func (vm *VM) SetAutoloadExtensions(value string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.autoloadExtensions = &value
}
func (vm *RequestVM) AutoloadExtensions() string {
	if vm.autoloadExtensions == nil {
		value := vm.Base.AutoloadExtensions()
		vm.autoloadExtensions = &value
	}
	return *vm.autoloadExtensions
}
func (vm *RequestVM) SetAutoloadExtensions(value string) { vm.autoloadExtensions = &value }

// Callable identity is retained separately from the resolved execution scope.
func (vm *VM) autoloadSnapshot() []autoloadRegistration {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if !vm.autoloadInitialized {
		for _, fn := range parser.GetAutoLoad() {
			vm.autoload = append(vm.autoload, autoloadRegistration{fn, fn})
		}
		vm.autoloadInitialized = true
	}
	return append([]autoloadRegistration(nil), vm.autoload...)
}
func (vm *RequestVM) autoloadSnapshot() []autoloadRegistration {
	if !vm.autoloadInitialized {
		vm.autoloadInitialized = true
		scope := vm.RequestObjectScope()
		for _, entry := range vm.Base.autoloadSnapshot() {
			vm.autoload = append(vm.autoload, autoloadRegistration{scope.Bind(entry.original), scope.Bind(entry.resolved)})
		}
	}
	return append([]autoloadRegistration(nil), vm.autoload...)
}
func (vm *VM) AutoloadCallbacks() []data.Value        { return resolvedAutoload(vm.autoloadSnapshot()) }
func (vm *RequestVM) AutoloadCallbacks() []data.Value { return resolvedAutoload(vm.autoloadSnapshot()) }
func resolvedAutoload(entries []autoloadRegistration) []data.Value {
	out := make([]data.Value, len(entries))
	for i, entry := range entries {
		out[i] = entry.resolved
	}
	return out
}
func originalAutoload(entries []autoloadRegistration) []data.Value {
	out := make([]data.Value, len(entries))
	for i, entry := range entries {
		out[i] = entry.original
	}
	return out
}
func (vm *VM) AutoloadOriginals() []data.Value        { return originalAutoload(vm.autoloadSnapshot()) }
func (vm *RequestVM) AutoloadOriginals() []data.Value { return originalAutoload(vm.autoloadSnapshot()) }
func (vm *VM) RegisterAutoload(original, resolved data.Value, prepend bool) {
	vm.autoloadSnapshot()
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.autoload = registerAutoload(vm, vm.autoload, original, resolved, prepend)
}
func (vm *RequestVM) RegisterAutoload(original, resolved data.Value, prepend bool) {
	vm.autoloadSnapshot()
	vm.autoload = registerAutoload(vm, vm.autoload, original, resolved, prepend)
}
func registerAutoload(vm data.VM, entries []autoloadRegistration, original, resolved data.Value, prepend bool) []autoloadRegistration {
	for _, entry := range entries {
		if sameAutoload(vm, entry.original, original) {
			return entries
		}
	}
	entry := autoloadRegistration{original, resolved}
	if prepend {
		return append([]autoloadRegistration{entry}, entries...)
	}
	return append(entries, entry)
}
func (vm *VM) UnregisterAutoload(original data.Value) bool {
	vm.autoloadSnapshot()
	vm.mu.Lock()
	defer vm.mu.Unlock()
	var found bool
	vm.autoload, found = unregisterAutoload(vm, vm.autoload, original)
	return found
}
func (vm *RequestVM) UnregisterAutoload(original data.Value) bool {
	vm.autoloadSnapshot()
	var found bool
	vm.autoload, found = unregisterAutoload(vm, vm.autoload, original)
	return found
}
func unregisterAutoload(vm data.VM, entries []autoloadRegistration, original data.Value) ([]autoloadRegistration, bool) {
	for i, entry := range entries {
		if sameAutoload(vm, entry.original, original) {
			copy(entries[i:], entries[i+1:])
			entries[len(entries)-1] = autoloadRegistration{}
			return entries[:len(entries)-1], true
		}
	}
	return entries, false
}
func sameAutoload(vm data.VM, a, b data.Value) bool {
	if a == b {
		return true
	}
	if _, ok := a.(*data.FuncValue); ok {
		// First-class function callables are distinct Closure objects even
		// when they wrap the same declaration. Identity was checked above.
		return false
	}
	if x, ok := a.(*data.BoundFuncValue); ok {
		y, ok := b.(*data.BoundFuncValue)
		return ok && x == y
	}
	ax, am, ao := autoloadMethod(vm, a)
	bx, bm, bo := autoloadMethod(vm, b)
	if am != "" || bm != "" {
		return ao == bo && data.TypeNameEqual(ax, bx) && data.TypeNameEqual(am, bm)
	}
	x, xok := a.(*data.StringValue)
	y, yok := b.(*data.StringValue)
	return xok && yok && data.TypeNameEqual(strings.TrimPrefix(x.Value, "\\"), strings.TrimPrefix(y.Value, "\\"))
}
func autoloadMethod(vm data.VM, value data.Value) (string, string, *data.ObjectValue) {
	var receiver data.Value
	method := ""
	switch v := value.(type) {
	case *data.StringValue:
		if name, fn, ok := strings.Cut(v.Value, "::"); ok {
			receiver = data.NewStringValue(name)
			method = fn
		}
	case *data.ArrayValue:
		if v.Len() == 2 {
			if first, _ := v.FindSlotByIntKey(0); first != nil {
				receiver = first.ReadValue()
			}
			if second, _ := v.FindSlotByIntKey(1); second != nil {
				method = second.ReadValue().AsString()
			}
		}
	case *data.ClassValue:
		return "", "__invoke", v.ObjectValue
	case *data.ThisValue:
		return "", "__invoke", v.ObjectValue
	}
	switch v := receiver.(type) {
	case *data.ClassValue:
		return "", method, v.ObjectValue
	case *data.ThisValue:
		return "", method, v.ObjectValue
	case *data.StringValue:
		name := strings.TrimPrefix(v.Value, "\\")
		if class, ok := vm.GetClass(name); ok {
			name = class.GetName()
		}
		return name, method, nil
	}
	return "", "", nil
}

func (vm *VM) AddAutoload(fn *data.FuncValue)           { vm.RegisterAutoload(fn, fn, false) }
func (vm *RequestVM) AddAutoload(fn *data.FuncValue)    { vm.RegisterAutoload(fn, fn, false) }
func (vm *VM) RemoveAutoload(fn *data.FuncValue)        { vm.UnregisterAutoload(fn) }
func (vm *RequestVM) RemoveAutoload(fn *data.FuncValue) { vm.UnregisterAutoload(fn) }
func autoloadFunctions(callbacks []data.Value) []*data.FuncValue {
	out := make([]*data.FuncValue, 0, len(callbacks))
	for _, callback := range callbacks {
		if fn, ok := callback.(*data.FuncValue); ok {
			out = append(out, fn)
		}
	}
	return out
}
func (vm *VM) AutoloadFunctions() []*data.FuncValue { return autoloadFunctions(vm.AutoloadCallbacks()) }
func (vm *RequestVM) AutoloadFunctions() []*data.FuncValue {
	return autoloadFunctions(vm.AutoloadCallbacks())
}
