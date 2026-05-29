package runtime

import "github.com/php-any/origami/data"

// RunHeaderCallbacksFn 由 std/php.Load 注入，执行 header_register_callback 注册的回调。
var RunHeaderCallbacksFn func(vm data.VM)

func runHeaderCallbacks(vm data.VM) {
	if RunHeaderCallbacksFn != nil {
		RunHeaderCallbacksFn(vm)
	}
}
