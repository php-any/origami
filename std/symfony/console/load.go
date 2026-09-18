package console

import "github.com/php-any/origami/data"

// Load 暂不注册不完整类，避免挡住 vendor 完整 ArgvInput/Output。
// 增量补齐后再 AddClass。
func Load(vm data.VM) {
	_ = vm
}
