# 反射函数注册功能

本功能允许使用反射将 Go 函数注册为脚本函数，无需手动实现`data.FuncStmt`接口。

## 功能特性

- **自动参数分析**: 使用反射自动分析 Go 函数的参数
- **类型转换**: 自动处理脚本类型和 Go 类型之间的转换
- **批量注册**: 支持批量注册多个函数
- **错误处理**: 完善的错误处理和类型检查

## 基本用法

### 1. 单个函数注册

```go
package main

import (
    "github.com/php-any/origami/parser"
    "github.com/php-any/origami/runtime"
)

// 定义Go函数
func Add(a, b int) int {
    return a + b
}

func Concat(a, b string) string {
    return a + b
}

func main() {
    // 创建VM
    p := parser.NewParser()
    vm := runtime.NewVM(p)

    // 注册函数
    vm.RegisterFunction("add", Add)
    vm.RegisterFunction("concat", Concat)

    // 或者直接注册匿名函数
    vm.RegisterFunction("multiply", func(a, b int) int {
        return a * b
    })
}
```

### 2. 批量函数注册

```go
// 逐个注册函数
vm.RegisterFunction("add", Add)
vm.RegisterFunction("multiply", Multiply)
vm.RegisterFunction("concat", Concat)
vm.RegisterFunction("isEven", IsEven)

// 或者使用辅助函数批量注册
functions := map[string]interface{}{
    "add":      Add,
    "multiply": Multiply,
    "concat":   Concat,
    "isEven":   IsEven,
}

// 批量注册
for name, fn := range functions {
    vm.RegisterFunction(name, fn)
}
```

## 支持的 Go 类型

### 输入参数类型

- `string` → 脚本字符串
- `int` → 脚本整数
- `float64` → 脚本浮点数
- `bool` → 脚本布尔值

### 返回值类型

- `string` → 脚本字符串
- `int` → 脚本整数
- `float64` → 脚本浮点数
- `bool` → 脚本布尔值
- `int64` → 脚本整数
- 其他类型 → 转换为字符串

## 示例函数

### 数学函数

```go
func Add(a, b int) int {
    return a + b
}

func Multiply(a, b int) int {
    return a * b
}

func Divide(a, b float64) float64 {
    if b == 0 {
        return 0
    }
    return a / b
}
```

### 字符串函数

```go
func ToUpperCase(s string) string {
    return strings.ToUpper(s)
}

func ToLowerCase(s string) string {
    return strings.ToLower(s)
}

func Concat(a, b string) string {
    return a + b
}
```

### 时间函数

```go
func GetCurrentTime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

func GetTimestamp() int64 {
    return time.Now().Unix()
}
```

### 逻辑函数

```go
func IsEven(n int) bool {
    return n%2 == 0
}

func IsOdd(n int) bool {
    return n%2 != 0
}

func IsEmpty(s string) bool {
    return len(strings.TrimSpace(s)) == 0
}
```

## 在脚本中使用

注册后的函数可以在脚本中直接调用：

```php
<?php
// 数学运算
int $result = add(10, 20);
echo "加法结果: " + $result + "\n";

int $product = multiply(5, 6);
echo "乘法结果: " + $product + "\n";

// 字符串操作
string $upper = toUpperCase("hello world");
echo "大写: " + $upper + "\n";

string $combined = concat("Hello", " World");
echo "拼接: " + $combined + "\n";

// 时间函数
string $time = getCurrentTime();
echo "当前时间: " + $time + "\n";

int $timestamp = getTimestamp();
echo "时间戳: " + $timestamp + "\n";

// 逻辑判断
bool $isEven = isEven(10);
echo "10是偶数: " + $isEven + "\n";

bool $isEmpty = isEmpty("");
echo "空字符串: " + $isEmpty + "\n";
?>
```

## 错误处理

### 类型转换错误

如果脚本传递的参数类型无法转换为 Go 函数期望的类型，会抛出错误：

```php
<?php
// 错误：传递字符串给期望int的函数
int $result = add("hello", 20); // 会抛出类型转换错误
?>
```

### 函数不存在错误

如果调用未注册的函数，会抛出函数未定义错误：

```php
<?php
// 错误：函数未注册
unknownFunction(); // 会抛出函数未定义错误
?>
```

## 最佳实践

### 1. 函数命名

- 使用清晰的函数名
- 避免与内置函数冲突
- 使用小写字母和下划线命名

### 2. 参数设计

- 保持参数数量合理（建议不超过 5 个）
- 使用基本类型作为参数
- 提供合理的默认值处理

### 3. 返回值处理

- 确保返回值类型一致
- 处理异常情况（如除零）
- 返回有意义的错误信息

### 4. 性能考虑

- 避免在函数中进行大量计算
- 对于复杂操作，考虑异步处理
- 合理使用缓存机制

## 测试

运行测试来验证功能：

```bash
go test ./runtime -v
```

## 注意事项

1. **类型安全**: 反射注册的函数在运行时进行类型检查，确保类型安全
2. **性能**: 反射调用比直接调用稍慢，但对于脚本函数来说是可以接受的
3. **错误处理**: 所有类型转换错误都会被捕获并转换为脚本异常
4. **线程安全**: 注册的函数是线程安全的，可以在多个 goroutine 中使用

## 扩展

如果需要支持更多类型，可以在`convertToGoValue`和`convertToScriptValue`方法中添加相应的类型转换逻辑。

## 注册示例

```go
// 在main函数中注册示例函数
vm.RegisterFunction("add", Add)
vm.RegisterFunction("multiply", Multiply)
vm.RegisterFunction("divide", Divide)
vm.RegisterFunction("toUpperCase", ToUpperCase)
vm.RegisterFunction("toLowerCase", ToLowerCase)
vm.RegisterFunction("concat", Concat)
vm.RegisterFunction("getCurrentTime", GetCurrentTime)
vm.RegisterFunction("getTimestamp", GetTimestamp)
vm.RegisterFunction("isEven", IsEven)
vm.RegisterFunction("isOdd", IsOdd)
vm.RegisterFunction("isEmpty", IsEmpty)
```
