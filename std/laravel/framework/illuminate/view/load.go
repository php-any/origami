package view

import (
	"os"

	"github.com/php-any/origami/data"
)

// Load 注册 illuminate/view 热类。默认关闭（方法面仍在补齐）；ORIGAMI_STD_VIEW=1 启用。
func Load(vm data.VM) {
	v := os.Getenv("ORIGAMI_STD_VIEW")
	if v != "1" && v != "true" && v != "yes" {
		return
	}
	vm.AddClass(NewAppendableClass())
	vm.AddClass(NewAttributeBagClass())
}
