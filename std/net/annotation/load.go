package annotation

import (
	"github.com/php-any/origami/data"
)

// Load 加载注解模块
func Load(vm data.VM) {
	// 注册注解类
	vm.AddClass(&ApplicationClass{})
	vm.AddClass(&ControllerClass{})
	vm.AddClass(&RouteClass{})
	vm.AddClass(&InjectClass{})
	vm.AddClass(&GetMappingClass{})
	vm.AddClass(&PostMappingClass{})
	vm.AddFunc(newSpringInlineFunc())
}
