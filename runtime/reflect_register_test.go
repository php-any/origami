package runtime

import (
	"testing"

	"github.com/php-any/origami/parser"
)

// 测试用的Go函数
func testAdd(a, b int) int {
	return a + b
}

func testConcat(a, b string) string {
	return a + b
}

func testIsEven(n int) bool {
	return n%2 == 0
}

func testNoReturn() {
	// 无返回值的函数
}

func TestRegisterReflectFunction(t *testing.T) {
	// 创建解析器和VM
	p := parser.NewParser()
	vm := NewVM(p)

	// 使用新的RegisterFunction方法注册函数
	vm.RegisterFunction("add", testAdd)
	vm.RegisterFunction("concat", testConcat)
	vm.RegisterFunction("isEven", testIsEven)
	vm.RegisterFunction("noReturn", testNoReturn)

	// 验证函数是否注册成功
	expectedFunctions := []string{"add", "concat", "isEven", "noReturn"}
	for _, name := range expectedFunctions {
		if fn, ok := vm.GetFunc(name); !ok {
			t.Errorf("函数 %s 注册失败", name)
		} else if fn.GetName() != name {
			t.Errorf("函数名不匹配: 期望 %s, 实际 %s", name, fn.GetName())
		}
	}
}

func TestReflectFunctionCall(t *testing.T) {
	// 创建解析器和VM
	p := parser.NewParser()
	vm := NewVM(p)

	// 使用新的RegisterFunction方法注册测试函数
	vm.RegisterFunction("add", testAdd)
	vm.RegisterFunction("concat", testConcat)
	vm.RegisterFunction("isEven", testIsEven)

	// 获取函数
	addFn, ok := vm.GetFunc("add")
	if !ok {
		t.Fatal("无法获取add函数")
	}

	// 验证函数参数
	params := addFn.GetParams()
	if len(params) != 2 {
		t.Errorf("期望2个参数，实际有%d个", len(params))
	}

	variables := addFn.GetVariables()
	if len(variables) != 2 {
		t.Errorf("期望2个变量，实际有%d个", len(variables))
	}
}
