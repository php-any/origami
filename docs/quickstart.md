# 快速开始

本指南将帮助您在 5 分钟内创建并运行第一个折言程序。

## 创建第一个程序

### 1. 创建脚本文件

创建一个名为 `hello.cjp` 的文件：

```php
<?php
// 第一个折言程序
echo "Hello, Origami!";

// 变量声明
string $name = "World";
echo "Hello, {$name}!";

// 简单函数
function greet(string $name): string {
    return "Hello, " . $name . "!";
}

$message = greet("Origami");
echo $message;
```

### 2. 运行程序

```bash
./origami hello.cjp
```

输出：

```
Hello, Origami!
Hello, World!
Hello, Origami!
```

## 基础示例

### 变量和数据类型

```php
<?php
// 基本类型
int $age = 25;
string $name = "Alice";
bool $isStudent = true;
float $height = 1.75;

// 数组
array $colors = ["red", "green", "blue"];
array $person = [
    "name" => "Bob",
    "age" => 30
];

// 输出
echo "Name: {$name}, Age: {$age}";
echo "Colors: " + $colors->join(", ");
```

### 控制结构

```php
<?php
// 条件语句
int $score = 85;

if ($score >= 90) {
    echo "优秀";
} elseif ($score >= 80) {
    echo "良好";
} else {
    echo "需要努力";
}

// 循环
for (int $i = 1; $i <= 5; $i++) {
    echo "Count: {$i}\n";
}

// foreach 循环
array $fruits = ["apple", "banana", "orange"];
foreach ($fruits as $fruit) {
    echo "I like {$fruit}\n";
}
```

### 函数定义

```php
<?php
// 基本函数
function add(int $a, int $b): int {
    return $a + $b;
}

// 带默认值的函数
function greet(string $name = "World"): string {
    return "Hello, {$name}!";
}

// 函数调用
$result = add(10, 20);
echo "Sum: " + $result + "\n";

$message = greet("Alice");
echo $message;
```

### 类和对象

```php
<?php
// 类定义
class Person {
    private string $name;
    private int $age;

    public function __construct(string $name, int $age) {
        $this->name = $name;
        $this->age = $age;
    }

    public function introduce(): string {
        return "I'm {$this->name}, {$this->age} years old.";
    }

    public function getAge(): int {
        return $this->age;
    }
}

// 创建对象
$person = new Person("Alice", 25);
echo $person->introduce();
```

### 字符串操作

```php
<?php
string $text = "Hello World";

// 字符串方法
echo "Length: " + $text->length() + "\n";
echo "Uppercase: " + $text->toUpperCase() + "\n";
echo "Lowercase: " + $text->toLowerCase() + "\n";
echo "Contains 'World': " + ($text->indexOf("World") >= 0 ? "Yes" : "No") + "\n";
```

### 数组操作

```php
<?php
array $numbers = [1, 2, 3, 4, 5];

// 数组方法
echo "Original: " + $numbers->join(", ") + "\n";

$doubled = $numbers->map(function($n) {
    return $n * 2;
});
echo "Doubled: " + $doubled->join(", ") + "\n";

$evens = $numbers->filter(function($n) {
    return $n % 2 == 0;
});
echo "Evens: " + $evens->join(", ") + "\n";
```

## 实用示例

### 简单计算器

```php
<?php
function calculate(string $operation, float $a, float $b): float {
    switch ($operation) {
        case "add":
            return $a + $b;
        case "subtract":
            return $a - $b;
        case "multiply":
            return $a * $b;
        case "divide":
            if ($b == 0) {
                throw new Exception("Division by zero");
            }
            return $a / $b;
        default:
            throw new Exception("Unknown operation");
    }
}

// 测试计算器
try {
    echo "10 + 5 = " + calculate("add", 10, 5) + "\n";
echo "10 - 5 = " + calculate("subtract", 10, 5) + "\n";
echo "10 * 5 = " + calculate("multiply", 10, 5) + "\n";
echo "10 / 5 = " + calculate("divide", 10, 5) + "\n";
} catch (Exception $e) {
    echo "Error: " + $e->getMessage() + "\n";
}
```

### 文件处理

```php
<?php
// 创建简单的日志系统
function logMessage(string $level, string $message): void {
    $timestamp = date("Y-m-d H:i:s");
    $logEntry = "[{$timestamp}] [{$level}] {$message}\n";
    echo $logEntry;
}

logMessage("INFO", "Application started");
logMessage("WARN", "This is a warning message");
logMessage("ERROR", "Something went wrong");
```

## 运行测试

验证语言功能：

```bash
./origami tests/run_tests.cjp
```

## 下一步

现在您已经了解了基础用法，建议：

1. **学习更多语法**: 阅读 [基础语法](syntax.md)
2. **探索标准库**: 查看 [标准库](stdlib.md) 文档
3. **集成 Go 代码**: 学习 [Go 集成](go-integration.md)
4. **查看示例**: 浏览 `tests/` 目录中的示例代码

## 常见问题

### Q: 如何调试程序？

A: 使用 `Log::debug()` 或 `echo` 输出调试信息。

### Q: 支持哪些文件扩展名？

A: 支持 `.cjp` 和 `.php` 文件。

### Q: 如何获取帮助？

A: 运行 `./origami` 查看命令行帮助。

### Q: 如何报告问题？

A: 在 GitHub 上提交 Issue 或加入讨论群。
