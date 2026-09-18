package routing

import "github.com/php-any/origami/data"

// Load 注册 Symfony Routing 热路径原生类（对齐 v8.1 / Laravel 13）。
// 不注册 Router / Loader / DI：那些依赖 Config 包，继续走 vendor PHP。
func Load(vm data.VM) {
	vm.AddInterface(newExceptionInterface())
	for _, iface := range routingInterfaces() {
		vm.AddInterface(iface)
	}
	for _, c := range routingExceptionClasses() {
		vm.AddClass(c)
	}
	vm.AddClass(newAliasClass())
	vm.AddClass(newCompiledRouteClass())
	vm.AddClass(newRouteCompilerClass())
	vm.AddClass(newRouteClass())
	vm.AddClass(newRouteCollectionClass())
	vm.AddClass(newRequestContextClass())
	vm.AddClass(newMatcherDumperClass())
	vm.AddClass(newUrlMatcherClass())
	vm.AddClass(newCompiledUrlMatcherClass())
	vm.AddClass(newCompiledUrlMatcherDumperClass())
	vm.AddClass(newUrlGeneratorClass())
	vm.AddClass(newAttributeRouteClass())
}
