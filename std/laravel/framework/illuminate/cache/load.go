package cache

import "github.com/php-any/origami/data"

// Load：Repository 方法面未覆盖 Store/TaggedCache 前不占名。
func Load(vm data.VM) {
	_ = vm
	_ = NewRepositoryClass
}
