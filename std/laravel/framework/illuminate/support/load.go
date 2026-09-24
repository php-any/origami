package support

import "github.com/php-any/origami/data"

func Load(vm data.VM) {
	registerHelpers(vm)
	// Ramsey\Uuid 的接口用 Go 原生注册：Str::uuid() 家族的返回值要能通过 instanceof，
	// 且不依赖 vendor/ramsey/uuid 的 PHP 文件是否被解析。
	// 注意只注册接口、不注册 Ramsey\Uuid\Uuid 类，避免顶掉 vendor 里 UuidV4 extends Uuid 的继承链。
	vm.AddInterface(NewRamseyUuidInterface())
	vm.AddInterface(NewRamseyRfc4122UuidInterface())
	vm.AddClass(NewStrClass())
	vm.AddClass(NewHtmlStringClass())
	vm.AddClass(NewStringableClass())
}
