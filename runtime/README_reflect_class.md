# 反射类功能文档

## 概述

反射类功能允许将 Go 结构体自动转换为 Origami 脚本类，无需手动编写包装代码。系统会通过反射分析 Go 结构体的公开方法，并生成对应的脚本类方法。

## 功能特性

- **自动方法发现**: 通过反射自动发现 Go 结构体的公开方法
- **类型转换**: 自动处理 Go 类型和脚本类型之间的转换
- **简单注册**: 支持单个结构体注册
- **错误处理**: 提供完善的错误处理机制

## 使用方法

### 1. 定义 Go 结构体

```go
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

// GetName 获取计算器名称
func (c *Calculator) GetName() string {
    return c.name
}
```

### 2. 注册反射类

```go
// 创建VM
p := parser.NewParser()
vm := NewVM(p)

// 注册单个类
calculator := NewCalculator("MyCalculator")
vm.RegisterReflectClass("Calculator", calculator)

// 注册多个类
calculator := NewCalculator("MyCalc")
vm.RegisterReflectClass("Calculator", calculator)

stringProcessor := &StringProcessor{}
vm.RegisterReflectClass("StringProcessor", stringProcessor)

timeUtils := &TimeUtils{}
vm.RegisterReflectClass("TimeUtils", timeUtils)
```

### 3. 在脚本中使用

```php
<?php
// 创建类实例
$calc = new Calculator("MyCalculator");

// 调用方法
int $result = $calc->Add(10, 20);
echo "10 + 20 = " + $result + "\n";

int $product = $calc->Multiply(5, 6);
echo "5 * 6 = " + $product + "\n";

string $name = $calc->GetName();
echo "计算器名称: " + $name + "\n";
?>
```

## 支持的类型

### 输入类型（脚本到 Go）

- `string` → `string`
- `int` → `int`, `int64`
- `float` → `float64`
- `bool` → `bool`

### 输出类型（Go 到脚本）

- `string` → `string`
- `int`, `int64` → `int`
- `float64` → `float`
- `bool` → `bool`
- 其他类型 → `string`（通过`fmt.Sprintf`转换）

## 方法发现规则

1. **公开方法**: 只有首字母大写的公开方法会被发现
2. **方法签名**: 支持任意数量的参数和返回值
3. **接收者**: 方法必须绑定到结构体上

### 示例

```go
type MyStruct struct{}

// ✅ 会被发现（公开方法）
func (m *MyStruct) PublicMethod() string {
    return "public"
}

// ❌ 不会被发现（私有方法）
func (m *MyStruct) privateMethod() string {
    return "private"
}

// ✅ 会被发现（带参数和返回值）
func (m *MyStruct) ProcessData(input string, count int) (string, error) {
    return "processed", nil
}
```

## 错误处理

### 注册错误

- 类名冲突
- 结构体类型错误
- 方法分析失败

### 调用错误

- 参数类型转换失败
- 返回值类型转换失败
- 方法不存在

## 最佳实践

### 1. 方法命名

```go
// 推荐：清晰的方法名
func (c *Calculator) Add(a, b int) int { ... }
func (c *Calculator) GetName() string { ... }

// 避免：过于复杂的方法名
func (c *Calculator) CalculateMathematicalOperationOfAddition(a, b int) int { ... }
```

### 2. 参数设计

```go
// 推荐：简单明确的参数
func (c *Calculator) Add(a, b int) int { ... }

// 避免：过于复杂的参数结构
func (c *Calculator) Process(data map[string]interface{}, options []string) { ... }
```

### 3. 返回值设计

```go
// 推荐：简单明确的返回值
func (c *Calculator) GetName() string { ... }

// 避免：复杂的返回值结构
func (c *Calculator) GetInfo() (string, int, map[string]interface{}, error) { ... }
```

## 限制

1. **类型支持**: 目前只支持基本类型，不支持复杂类型（如切片、映射等）
2. **方法复杂度**: 建议保持方法简单，避免过于复杂的参数和返回值
3. **性能**: 反射调用会有一定的性能开销
4. **错误处理**: 复杂的错误处理可能需要额外的包装

## 示例代码

完整的示例代码请参考：

- `runtime/reflect_class_example.go` - 示例结构体定义
- `runtime/reflect_class_test.go` - 单元测试
- `tests/reflect_class_test.zy` - 脚本测试示例

## 测试

运行单元测试：

```bash
go test ./runtime -v -run TestRegisterReflectClass
```

运行所有反射相关测试：

```bash
go test ./runtime -v -run TestReflect
```
