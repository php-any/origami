package hashing

import "github.com/php-any/origami/data"

// Load：BcryptHasher 常开（HashManager 仍走 vendor）。
func Load(vm data.VM) {
	vm.AddClass(NewBcryptHasherClass())
}
