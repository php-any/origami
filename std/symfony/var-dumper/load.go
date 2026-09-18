package vardumper

import "github.com/php-any/origami/data"

// Load 暂不注册不完整 VarDumper，避免挡住 vendor。
func Load(vm data.VM) {
	_ = vm
}
