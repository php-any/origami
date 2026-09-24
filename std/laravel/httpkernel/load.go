package httpkernel

import "github.com/php-any/origami/data"

// Load 注册应用侧 HTTP 桥：App\Http\Kernel。
//
// vendor（Symfony / Illuminate）原生类由 std/vendoraccel + std/laravel/framework 注册；
// 本包只补 Laravel 官方骨架里由应用自己声明、但 Origami 用 Go 顶替的那一层。
// 类名沿用 Laravel 约定（bootstrap/app.php 会把 Contracts\Http\Kernel 绑到它），
// 因此不需要改动应用的 PHP 代码。
func Load(vm data.VM) {
	vm.AddClass(NewClass())
}
