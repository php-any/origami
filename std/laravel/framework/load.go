package framework

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/illuminate/collections"
	"github.com/php-any/origami/std/laravel/framework/illuminate/config"
	"github.com/php-any/origami/std/laravel/framework/illuminate/container"
	"github.com/php-any/origami/std/laravel/framework/illuminate/events"
	"github.com/php-any/origami/std/laravel/framework/illuminate/foundation"
	illuminatehttp "github.com/php-any/origami/std/laravel/framework/illuminate/http"
	"github.com/php-any/origami/std/laravel/framework/illuminate/routing"
	"github.com/php-any/origami/std/laravel/framework/illuminate/support"
	"github.com/php-any/origami/std/laravel/framework/illuminate/view"
)

// Load 按依赖序注册 laravel/framework 内嵌的 Illuminate 组件。
func Load(vm data.VM) {
	// 叶子工具包优先
	collections.Load(vm)
	support.Load(vm)
	config.Load(vm)
	events.Load(vm)
	container.Load(vm)
	illuminatehttp.Load(vm)
	view.Load(vm)
	routing.Load(vm)
	foundation.Load(vm)
}
