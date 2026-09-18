package gosupport

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/examples/laravel13/go-support/httpkernel"
)

// Load 注册 Laravel 13 示例侧适配（非 vendor 加速）。
//
// vendor（Symfony / Illuminate）原生类由 std/vendoraccel.Load 注册。
// 本包只保留：
//   - App\Http\Kernel（应用侧 HTTP 桥）
//   - ServeCommand（Go net/http 替换 php -S）
func Load(vm data.VM) {
	vm.AddClass(httpkernel.NewClass())
	vm.AddClass(NewServeCommandClass())
}
