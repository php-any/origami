package serve

import "github.com/php-any/origami/data"

// Load 注册 Go net/http 版的 Illuminate\Foundation\Console\ServeCommand。
//
// 在 Laravel 加载官方 ServeCommand.php 前以同名类抢先注册，保留 Artisan 命令的
// 生命周期与参数解析，同时用 Go HTTP 服务器顶替 php -S。
func Load(vm data.VM) {
	vm.AddClass(NewServeCommandClass())
}
