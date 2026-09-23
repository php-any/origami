package reflection

import "github.com/php-any/origami/data"

// Load：getClassAttributes / 参数反射与容器联调未稳前不占名；实现保留。
func Load(vm data.VM) {
	_ = vm
	_ = NewReflectorClass
}
