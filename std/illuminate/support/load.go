package support

import "github.com/php-any/origami/data"

// Load 注册 Illuminate\Support 热路径（Arr / helpers）。
// 不 AddClass Collection / Str：Go 实现不完整会挡住 vendor 官方类
//（PackageManifest::flatMap、EncryptionServiceProvider::Str::after 等）。
func Load(vm data.VM) {
	vm.AddClass(NewArrClass())
	registerHelpers(vm)
}
