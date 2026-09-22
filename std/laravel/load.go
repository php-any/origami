package laravel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework"
)

// Load 注册 vendor/laravel 原生加速层（仅应由 laravel13 / vendoraccel 调用）。
// 后续 telescope / prompts / serializable-closure 等平级包在此追加。
func Load(vm data.VM) {
	framework.Load(vm)
}
