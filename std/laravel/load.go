package laravel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework"
	"github.com/php-any/origami/std/laravel/httpkernel"
	"github.com/php-any/origami/std/laravel/prompts"
	serializableclosure "github.com/php-any/origami/std/laravel/serializable-closure"
	"github.com/php-any/origami/std/laravel/sentinel"
	"github.com/php-any/origami/std/laravel/serve"
	"github.com/php-any/origami/std/laravel/telescope"
	"github.com/php-any/origami/std/laravel/tinker"
)

// Load 注册 Laravel 运行时（仅应由 examples/laravel13 / vendoraccel 调用，禁止加入 zy.go）。
//
// 两层内容：
//   - framework/* 与 serializable-closure 等：vendor 各 Composer 包的原生加速层；
//   - httpkernel / serve：Origami 用 Go 顶替的应用侧 HTTP 桥与开发服务器，
//     不再由示例工程自带。
func Load(vm data.VM) {
	framework.Load(vm)
	serializableclosure.Load(vm)
	sentinel.Load(vm)
	tinker.Load(vm)
	prompts.Load(vm)
	telescope.Load(vm)

	httpkernel.Load(vm)
	serve.Load(vm)
}
