package build

import (
	"github.com/php-any/origami/data"
)

// Register 将预编译的 vendor AST 注册到 VM
func Register(vm data.VM) {
	vm.RegisterCompiledFile("tests\\compile_test\\lib\\test\\lib.php", func() (data.GetValue, []data.Variable) {
		return AST_Tests_compile_test_lib_test_lib()
	})
}
