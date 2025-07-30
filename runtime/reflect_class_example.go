package runtime

import (
	"fmt"
	"strings"
	"time"

	"github.com/php-any/origami/parser"
)

// Calculator 示例结构体
type Calculator struct {
	name string
}

// NewCalculator 创建计算器实例
func NewCalculator(name string) *Calculator {
	return &Calculator{name: name}
}

// Add 加法方法
func (c *Calculator) Add(a, b int) int {
	return a + b
}

// Multiply 乘法方法
func (c *Calculator) Multiply(a, b int) int {
	return a * b
}

// Divide 除法方法
func (c *Calculator) Divide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

// GetName 获取计算器名称
func (c *Calculator) GetName() string {
	return c.name
}

// Greet 问候方法
func (c *Calculator) Greet(name string) string {
	return fmt.Sprintf("Hello %s, I'm %s!", name, c.name)
}

// StringProcessor 字符串处理器
type StringProcessor struct{}

// ToUpperCase 转换为大写
func (sp *StringProcessor) ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase 转换为小写
func (sp *StringProcessor) ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Concat 字符串拼接
func (sp *StringProcessor) Concat(a, b string) string {
	return a + b
}

// Split 字符串分割
func (sp *StringProcessor) Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// TimeUtils 时间工具
type TimeUtils struct{}

// GetCurrentTime 获取当前时间
func (tu *TimeUtils) GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetTimestamp 获取时间戳
func (tu *TimeUtils) GetTimestamp() int64 {
	return time.Now().Unix()
}

// FormatTime 格式化时间
func (tu *TimeUtils) FormatTime(timestamp int64, format string) string {
	t := time.Unix(timestamp, 0)
	return t.Format(format)
}

// RegisterExampleClasses 注册示例类
func RegisterExampleClasses(vm *VM) {
	// 注册计算器类
	calculator := NewCalculator("MyCalculator")
	vm.RegisterReflectClass("Calculator", calculator)

	// 注册字符串处理器类
	stringProcessor := &StringProcessor{}
	vm.RegisterReflectClass("StringProcessor", stringProcessor)

	// 注册时间工具类
	timeUtils := &TimeUtils{}
	vm.RegisterReflectClass("TimeUtils", timeUtils)

	fmt.Println("示例类注册完成")
}

// ExampleUsage 使用示例
func ExampleUsage() {
	// 创建VM
	p := parser.NewParser()
	vm := NewVM(p)

	// 注册示例类
	RegisterExampleClasses(vm.(*VM))

	// 验证注册的类
	expectedClasses := []string{"Calculator", "StringProcessor", "TimeUtils"}

	for _, name := range expectedClasses {
		if class, ok := vm.GetClass(name); ok {
			fmt.Printf("✓ 类 %s 注册成功\n", name)
			methods := class.GetMethods()
			fmt.Printf("  方法数量: %d\n", len(methods))
			for _, method := range methods {
				fmt.Printf("    - %s\n", method.GetName())
			}
		} else {
			fmt.Printf("✗ 类 %s 注册失败\n", name)
		}
	}
}
