package build

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func AST_Tests_compile_test_lib_test_lib() (data.GetValue, []data.Variable) {
from := node.NewFrom("tests\\compile_test\\lib\\test\\lib.php")

stmts := []data.GetValue{
nil /* TODO: FunctionStatement "test_add" — 太复杂，暂不支持编译期代码生成 */,
}

vars := []data.Variable{
}

return node.NewProgram(from, stmts), vars
}


