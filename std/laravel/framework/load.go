package framework

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/illuminate/auth"
	"github.com/php-any/origami/std/laravel/framework/illuminate/broadcasting"
	"github.com/php-any/origami/std/laravel/framework/illuminate/bus"
	"github.com/php-any/origami/std/laravel/framework/illuminate/cache"
	"github.com/php-any/origami/std/laravel/framework/illuminate/collections"
	"github.com/php-any/origami/std/laravel/framework/illuminate/concurrency"
	"github.com/php-any/origami/std/laravel/framework/illuminate/conditionable"
	"github.com/php-any/origami/std/laravel/framework/illuminate/config"
	"github.com/php-any/origami/std/laravel/framework/illuminate/console"
	"github.com/php-any/origami/std/laravel/framework/illuminate/container"
	"github.com/php-any/origami/std/laravel/framework/illuminate/contracts"
	"github.com/php-any/origami/std/laravel/framework/illuminate/cookie"
	"github.com/php-any/origami/std/laravel/framework/illuminate/database"
	"github.com/php-any/origami/std/laravel/framework/illuminate/encryption"
	"github.com/php-any/origami/std/laravel/framework/illuminate/events"
	"github.com/php-any/origami/std/laravel/framework/illuminate/filesystem"
	"github.com/php-any/origami/std/laravel/framework/illuminate/foundation"
	"github.com/php-any/origami/std/laravel/framework/illuminate/hashing"
	illuminatehttp "github.com/php-any/origami/std/laravel/framework/illuminate/http"
	"github.com/php-any/origami/std/laravel/framework/illuminate/image"
	jsonschema "github.com/php-any/origami/std/laravel/framework/illuminate/json-schema"
	"github.com/php-any/origami/std/laravel/framework/illuminate/log"
	"github.com/php-any/origami/std/laravel/framework/illuminate/macroable"
	"github.com/php-any/origami/std/laravel/framework/illuminate/mail"
	"github.com/php-any/origami/std/laravel/framework/illuminate/notifications"
	"github.com/php-any/origami/std/laravel/framework/illuminate/pagination"
	"github.com/php-any/origami/std/laravel/framework/illuminate/pipeline"
	"github.com/php-any/origami/std/laravel/framework/illuminate/process"
	"github.com/php-any/origami/std/laravel/framework/illuminate/queue"
	"github.com/php-any/origami/std/laravel/framework/illuminate/redis"
	"github.com/php-any/origami/std/laravel/framework/illuminate/reflection"
	"github.com/php-any/origami/std/laravel/framework/illuminate/routing"
	"github.com/php-any/origami/std/laravel/framework/illuminate/session"
	"github.com/php-any/origami/std/laravel/framework/illuminate/support"
	"github.com/php-any/origami/std/laravel/framework/illuminate/testing"
	"github.com/php-any/origami/std/laravel/framework/illuminate/translation"
	"github.com/php-any/origami/std/laravel/framework/illuminate/validation"
	"github.com/php-any/origami/std/laravel/framework/illuminate/view"
)

// Load 按依赖序注册 laravel/framework 内嵌的 Illuminate 组件。
// 各子包仅在公开方法面对齐后 AddClass；未齐则 Load 为空（无 ORIGAMI_STD_* 开关）。
func Load(vm data.VM) {
	macroable.Load(vm)
	conditionable.Load(vm)
	pipeline.Load(vm)
	reflection.Load(vm)
	encryption.Load(vm)
	cookie.Load(vm)
	hashing.Load(vm)
	concurrency.Load(vm)
	collections.Load(vm)
	support.Load(vm)
	filesystem.Load(vm)
	log.Load(vm)
	config.Load(vm)
	events.Load(vm)
	container.Load(vm)
	translation.Load(vm)
	bus.Load(vm)
	process.Load(vm)
	contracts.Load(vm)
	cache.Load(vm)
	session.Load(vm)
	pagination.Load(vm)
	redis.Load(vm)
	broadcasting.Load(vm)
	notifications.Load(vm)
	auth.Load(vm)
	validation.Load(vm)
	jsonschema.Load(vm)
	console.Load(vm)
	database.Load(vm)
	queue.Load(vm)
	mail.Load(vm)
	image.Load(vm)
	testing.Load(vm)
	illuminatehttp.Load(vm)
	view.Load(vm)
	routing.Load(vm)
	foundation.Load(vm)
}
