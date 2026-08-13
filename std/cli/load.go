package cli

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/cli/annotation"
)

// Load 加载 CLI 模块
func Load(vm data.VM) {
	annotation.Load(vm)
}
