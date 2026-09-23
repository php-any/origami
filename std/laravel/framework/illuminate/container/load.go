package container

import "github.com/php-any/origami/data"

// Load 注册 illuminate/container。
// Container 方法面（Application extends、完整 autowire）未齐前不 AddClass，走 vendor PHP。
func Load(vm data.VM) {
	_ = vm
	_ = NewContainerClass // 保留实现，待方法面对齐后注册
}

const ComposerName = "illuminate/container"
const TargetVersion = "v13.23.0"
