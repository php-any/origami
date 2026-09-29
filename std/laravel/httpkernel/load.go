package httpkernel

import "github.com/php-any/origami/data"

// Load 在 Composer autoload 执行前注册 Go 版 Illuminate\Foundation\Http\Kernel。
//
// vendor（Symfony / Illuminate）原生类由 std/vendoraccel + std/laravel/framework 注册；
// ApplicationBuilder::withKernels() 会将 Contracts\Http\Kernel 绑到这个官方 FQN，
// 因此 bootstrap/app.php 保持 Laravel 13 官方内容，无需应用侧重新绑定。
func Load(vm data.VM) {
	vm.AddClass(NewClass())
}
