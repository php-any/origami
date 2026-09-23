package collections

import "github.com/php-any/origami/data"

// Load 注册 illuminate/collections：Arr + Collection + HigherOrderCollectionProxy + helpers。
func Load(vm data.VM) {
	vm.AddClass(NewArrClass())
	vm.AddInterface(NewEnumerableInterface())
	vm.AddClass(NewCollectionClass())
	vm.AddClass(NewHigherOrderProxyClass())
	vm.AddFunc(kitHelperCollect())
	vm.AddFunc(kitHelperDataGet())
	vm.AddFunc(kitHelperDataSet())
}
