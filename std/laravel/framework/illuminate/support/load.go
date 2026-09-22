package support

import "github.com/php-any/origami/data"

// Load 注册 illuminate/support：helpers；Str 方法面未齐前不 AddClass。
func Load(vm data.VM) {
	registerHelpers(vm)
	_ = NewStrClass
}
