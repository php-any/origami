package events

import (
	"os"

	"github.com/php-any/origami/data"
)

// Load 注册 illuminate/events。默认关闭；ORIGAMI_STD_EVENTS=1 启用。
func Load(vm data.VM) {
	v := os.Getenv("ORIGAMI_STD_EVENTS")
	if v != "1" && v != "true" && v != "yes" {
		return
	}
	vm.AddClass(NewDispatcherClass())
}
