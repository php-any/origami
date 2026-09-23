package laravel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework"
	"github.com/php-any/origami/std/laravel/prompts"
	serializableclosure "github.com/php-any/origami/std/laravel/serializable-closure"
	"github.com/php-any/origami/std/laravel/sentinel"
	"github.com/php-any/origami/std/laravel/telescope"
	"github.com/php-any/origami/std/laravel/tinker"
)

// Load 注册 vendor/laravel 原生加速层（仅应由 laravel13 / vendoraccel 调用）。
func Load(vm data.VM) {
	framework.Load(vm)
	serializableclosure.Load(vm)
	sentinel.Load(vm)
	tinker.Load(vm)
	prompts.Load(vm)
	telescope.Load(vm)
}
