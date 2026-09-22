package container

import "github.com/php-any/origami/data"

// Load 注册 illuminate/container。默认关闭；ORIGAMI_STD_CONTAINER=1 启用（spike）。
func Load(vm data.VM) {
	if !Enabled() {
		return
	}
	vm.AddClass(NewContainerClass())
}

const ComposerName = "illuminate/container"
const TargetVersion = "v13.23.0"
