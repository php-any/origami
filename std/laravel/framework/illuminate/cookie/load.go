package cookie

import "github.com/php-any/origami/data"

// Load 注册 illuminate/cookie。
func Load(vm data.VM) {
	vm.AddClass(NewCookieJarClass())
}
