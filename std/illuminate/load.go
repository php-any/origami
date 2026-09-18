package illuminate

import (
	"github.com/php-any/origami/data"
	illuminatehttp "github.com/php-any/origami/std/illuminate/http"
	"github.com/php-any/origami/std/illuminate/support"
)

// Load 注册 Illuminate vendor 原生加速层（仅应由 laravel13 / vendoraccel 调用）。
func Load(vm data.VM) {
	illuminatehttp.Load(vm)
	support.Load(vm)
}
