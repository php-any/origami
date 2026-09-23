package encryption

import "github.com/php-any/origami/data"

// Load 注册 illuminate/encryption（试开；会话异常则撤回占名）。
func Load(vm data.VM) {
	vm.AddClass(NewEncrypterClass())
}
