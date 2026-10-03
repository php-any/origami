package runtime

import "strings"

func (vm *VM) LookupPHPIni(name string) (string, bool) {
	value, found := vm.phpIni.Load(strings.ToLower(name))
	if !found {
		return "", false
	}
	return value.(string), true
}
func (vm *VM) StorePHPIni(name, value string) {
	vm.phpIni.Store(strings.ToLower(name), value)
	applyRequestIni(currentRequestCallState(), name, value)
}
func (vm *RequestVM) LookupPHPIni(name string) (string, bool) {
	if value, found := vm.phpIni[strings.ToLower(name)]; found {
		return value, true
	}
	return vm.Base.LookupPHPIni(name)
}
func (vm *RequestVM) StorePHPIni(name, value string) {
	if vm.phpIni == nil {
		vm.phpIni = make(map[string]string)
	}
	vm.phpIni[strings.ToLower(name)] = value
	applyRequestIni(vm.requestCall(), name, value)
}

func applyRequestIni(state *CallState, name, value string) {
	if state == nil {
		return
	}
	if strings.EqualFold(name, "ignore_user_abort") {
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "", "0", "off", "false", "no", "none":
			setIgnoreUserAbort(state, false)
		default:
			setIgnoreUserAbort(state, true)
		}
	}
}
