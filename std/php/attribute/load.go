package attribute

import (
	"github.com/php-any/origami/data"
)

// Load 注册所有 PHP 原生注解类
func Load(vm data.VM) {
	vm.AddClass(NewAttributeClass())
	vm.AddClass(NewDeprecatedClass())
	vm.AddClass(NewSensitiveParameterClass())
	vm.AddClass(NewOverrideClass())
	vm.AddClass(NewAllowDynamicPropertiesClass())
	vm.AddClass(NewReturnTypeWillChangeClass())
}
