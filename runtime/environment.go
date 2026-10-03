package runtime

import (
	"os"
	"runtime"
	"strings"
)

func environmentKey[V any](values map[string]V, name string) string {
	if _, found := values[name]; found || runtime.GOOS != "windows" {
		return name
	}
	for key := range values {
		if strings.EqualFold(key, name) {
			return key
		}
	}
	return name
}

func initialPHPEnvironment() *map[string]string {
	values := make(map[string]string)
	for _, entry := range os.Environ() {
		if name, value, ok := strings.Cut(entry, "="); ok && name != "" {
			values[name] = value
		}
	}
	return &values
}

func validPHPEnvironment(name string, value *string) bool {
	return name != "" && !strings.ContainsAny(name, "=\x00") && (value == nil || !strings.ContainsRune(*value, '\x00'))
}

func (vm *VM) PHPEnvironment() map[string]string {
	return copyPHPEnvironment(vm.environment.Load())
}
func copyPHPEnvironment(source *map[string]string) map[string]string {
	values := make(map[string]string)
	if source != nil {
		for name, value := range *source {
			values[name] = value
		}
	}
	return values
}
func (vm *VM) LookupPHPEnvironment(name string) (string, bool) {
	if values := vm.environment.Load(); values != nil {
		value, found := (*values)[environmentKey(*values, name)]
		return value, found
	}
	return "", false
}
func (vm *VM) SetPHPEnvironment(name string, value *string) bool {
	if !validPHPEnvironment(name, value) {
		return false
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	values := vm.PHPEnvironment()
	name = environmentKey(values, name)
	if value == nil {
		delete(values, name)
	} else {
		values[name] = *value
	}
	vm.environment.Store(&values)
	return true
}
func (vm *RequestVM) PHPEnvironment() map[string]string {
	values := copyPHPEnvironment(vm.initialEnvironment)
	for name, value := range vm.environment {
		if value == nil {
			delete(values, name)
		} else {
			values[name] = *value
		}
	}
	return values
}
func (vm *RequestVM) LookupPHPEnvironment(name string) (string, bool) {
	name = environmentKey(vm.environment, name)
	if value, found := vm.environment[name]; found {
		if value == nil {
			return "", false
		}
		return *value, true
	}
	if vm.initialEnvironment != nil {
		value, found := (*vm.initialEnvironment)[environmentKey(*vm.initialEnvironment, name)]
		return value, found
	}
	return "", false
}
func (vm *RequestVM) SetPHPEnvironment(name string, value *string) bool {
	if !validPHPEnvironment(name, value) {
		return false
	}
	if vm.initialEnvironment != nil {
		name = environmentKey(*vm.initialEnvironment, name)
	}
	name = environmentKey(vm.environment, name)
	if vm.environment == nil {
		vm.environment = make(map[string]*string)
	}
	if value != nil {
		owned := *value
		value = &owned
	}
	vm.environment[name] = value
	return true
}
