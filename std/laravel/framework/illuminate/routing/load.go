package routing

import "github.com/php-any/origami/data"

// Load 注册 illuminate/routing。默认关闭；ORIGAMI_STD_ROUTING=1 启用 UrlGenerator。
func Load(vm data.VM) {
	if !routingEnabled() {
		return
	}
	vm.AddClass(NewUrlGeneratorClass())
}

const ComposerName = "illuminate/routing"
const TargetVersion = "v13.23.0"
