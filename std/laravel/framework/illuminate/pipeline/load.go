package pipeline

import "github.com/php-any/origami/data"

// Load 注册 illuminate/pipeline。
func Load(vm data.VM) {
	vm.AddClass(NewPipelineClass())
}
