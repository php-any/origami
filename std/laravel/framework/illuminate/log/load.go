package log

import "github.com/php-any/origami/data"

// Load 注册 illuminate/log（LogManager 面过大，暂不 AddClass）。
func Load(vm data.VM) {
	_ = vm
}
