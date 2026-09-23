package translation

import "github.com/php-any/origami/data"

// Load：Translator 未覆盖加载器链前不占名。
func Load(vm data.VM) {
	_ = vm
	_ = NewTranslatorClass
}
