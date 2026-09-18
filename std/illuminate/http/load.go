package http

import "github.com/php-any/origami/data"

// Load 注册 Illuminate\Http\*（依赖 std/symfony/http-foundation）。
func Load(vm data.VM) {
	vm.AddClass(NewIlluminateRequestClass())
	vm.AddClass(NewIlluminateResponseClass())
}
