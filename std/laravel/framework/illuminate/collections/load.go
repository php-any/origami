package collections

import "github.com/php-any/origami/data"

// Load 注册 illuminate/collections。
// Arr + collect/data_get/data_set 默认开启。
// Collection 由 collections 阶段按 ORIGAMI_STD_COLLECTION 开关注册（见 enable_collection.go）。
func Load(vm data.VM) {
	vm.AddClass(NewArrClass())
	registerHelpers(vm)
	maybeRegisterCollection(vm)
}

func registerHelpers(vm data.VM) {
	vm.AddFunc(kitHelperCollect())
	vm.AddFunc(kitHelperDataGet())
	vm.AddFunc(kitHelperDataSet())
}
