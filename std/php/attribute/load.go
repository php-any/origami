package attribute

import (
	"github.com/php-any/origami/data"
)

// Load 注册所有 PHP 原生注解类
func Load(vm data.VM) {
	// 注册 PHP 8.0+ 原生注解类
	vm.AddClass(NewAttributeClass())
	vm.AddClass(NewDeprecatedClass())
	vm.AddClass(NewSensitiveParameterClass())
	vm.AddClass(NewReturnTypeWillChangeClass())
	vm.AddClass(NewAllowDynamicPropertiesClass())
}
