package runtime

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestRegisterReflectClass(t *testing.T) {
	// 创建解析器和VM
	p := parser.NewParser()
	vm := NewVM(p)

	// 创建示例结构体
	calculator := NewCalculator("TestCalculator")

	// 注册反射类
	ctl := vm.RegisterReflectClass("Calculator", calculator)
	if ctl != nil {
		t.Errorf("注册Calculator类失败: %v", ctl)
	}

	// 验证类是否注册成功
	if class, ok := vm.GetClass("Calculator"); !ok {
		t.Errorf("Calculator类注册失败")
	} else {
		// 验证类名
		if class.GetName() != "Calculator" {
			t.Errorf("类名不匹配: 期望 Calculator, 实际 %s", class.GetName())
		}

		// 验证方法
		methods := class.GetMethods()
		expectedMethods := []string{"Add", "Divide", "GetName", "Greet", "Multiply"}

		if len(methods) != len(expectedMethods) {
			t.Errorf("方法数量不匹配: 期望 %d, 实际 %d", len(expectedMethods), len(methods))
		}

		// 验证特定方法
		if method, ok := class.GetMethod("Add"); !ok {
			t.Errorf("Add方法不存在")
		} else {
			if method.GetName() != "Add" {
				t.Errorf("方法名不匹配: 期望 Add, 实际 %s", method.GetName())
			}
		}
	}
}

func TestReflectMethodCall(t *testing.T) {
	// 创建解析器和VM
	p := parser.NewParser()
	vm := NewVM(p)

	// 创建示例结构体
	calculator := NewCalculator("TestCalculator")

	// 注册反射类
	vm.RegisterReflectClass("Calculator", calculator)

	// 获取类和方法
	class, _ := vm.GetClass("Calculator")
	method, _ := class.GetMethod("Add")

	// 模拟方法调用（这里只是测试方法存在性，实际调用需要更复杂的上下文设置）
	if method == nil {
		t.Errorf("Add方法获取失败")
	}

	// 验证方法属性
	if method.GetModifier() != data.ModifierPublic {
		t.Errorf("方法修饰符不正确: 期望 public, 实际 %v", method.GetModifier())
	}

	if method.GetIsStatic() {
		t.Errorf("方法不应该是静态的")
	}
}

func TestReflectClassMethods(t *testing.T) {
	// 创建解析器和VM
	p := parser.NewParser()
	vm := NewVM(p)

	// 创建示例结构体
	stringProcessor := &StringProcessor{}

	// 注册反射类
	vm.RegisterReflectClass("StringProcessor", stringProcessor)

	// 验证类注册
	if class, ok := vm.GetClass("StringProcessor"); !ok {
		t.Errorf("StringProcessor类注册失败")
	} else {
		methods := class.GetMethods()
		expectedMethods := []string{"Concat", "Split", "ToLowerCase", "ToUpperCase"}

		if len(methods) != len(expectedMethods) {
			t.Errorf("StringProcessor方法数量不匹配: 期望 %d, 实际 %d", len(expectedMethods), len(methods))
		}

		// 验证特定方法
		for _, expectedMethod := range expectedMethods {
			if _, ok := class.GetMethod(expectedMethod); !ok {
				t.Errorf("方法 %s 不存在", expectedMethod)
			}
		}
	}
}

func TestReflectConstructor(t *testing.T) {
	// 创建解析器和VM
	p := parser.NewParser()
	vm := NewVM(p)

	// 创建示例结构体
	calculator := NewCalculator("TestCalculator")

	// 注册反射类
	vm.RegisterReflectClass("Calculator", calculator)

	// 验证类注册
	if class, ok := vm.GetClass("Calculator"); !ok {
		t.Errorf("Calculator类注册失败")
	} else {
		// 验证构造函数
		constructor := class.GetConstruct()
		if constructor == nil {
			t.Errorf("构造函数不存在")
		} else {
			if constructor.GetName() != "__construct" {
				t.Errorf("构造函数名不匹配: 期望 __construct, 实际 %s", constructor.GetName())
			}
		}
	}
}
