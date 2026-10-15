package laravel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework"
	"github.com/php-any/origami/std/laravel/httpkernel"
	"github.com/php-any/origami/std/laravel/prompts"
	"github.com/php-any/origami/std/laravel/sentinel"
	serializableclosure "github.com/php-any/origami/std/laravel/serializable-closure"
	"github.com/php-any/origami/std/laravel/telescope"
	"github.com/php-any/origami/std/laravel/tinker"
)

// Load 直接注册 Laravel 原生组件；每个组件的 Load 管理自己的实现。
// ServeCommand 是独立的 HTTP/进程适配器，由宿主按需加载 serve.Load。
func Load(vm data.VM) {
	framework.Load(vm)
	serializableclosure.Load(vm)
	sentinel.Load(vm)
	tinker.Load(vm)
	prompts.Load(vm)
	telescope.Load(vm)

	httpkernel.Load(vm)
}
