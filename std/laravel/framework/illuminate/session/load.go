package session

import "github.com/php-any/origami/data"

// Load：Store 方法面未齐前不占名。
func Load(vm data.VM) {
	_ = vm
	_ = NewStoreClass
}
