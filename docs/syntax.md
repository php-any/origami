# 基础语法

本文档介绍折言语言的基础语法和特性。

## 文件结构

### 文件扩展名

- `.zy` - 折言脚本文件
- `.php` - PHP 兼容脚本文件

### 基本结构

```php
<?php
// 命名空间声明（可选）
namespace myapp;

// 代码内容
echo "Hello World";
```

## 变量和数据类型

### 变量声明

```php
// 基本变量声明
$name = "Alice";
$age = 25;

// 类型声明
string $name = "Alice";
int $age = 25;
bool $isStudent = true;
float $height = 1.75;
```

### 数据类型

#### 基本类型

- `int` - 整数类型
- `string` - 字符串类型
- `bool` - 布尔类型
- `float` - 浮点数类型

#### 复合类型

- `array` - 数组类型
- `object` - 对象类型
- `class` - 类类型

#### 特殊类型

- `null` - 空值
- `void` - 无返回值

#### 可空类型

```php
?string $name = null;  // 可空字符串
?int $age = null;      // 可空整数
```

## 运算符

### 算术运算符

```php
$a = 10;
$b = 3;

$sum = $a + $b;      // 加法: 13
$diff = $a - $b;     // 减法: 7
$product = $a * $b;  // 乘法: 30
$quotient = $a / $b; // 除法: 3.333...
$remainder = $a % $b; // 取余: 1
```

### 比较运算符

```php
$a = 5;
$b = "5";

$a == $b;   // 相等: true
$a === $b;  // 严格相等: false
$a != $b;   // 不等: false
$a !== $b;  // 严格不等: true
$a < $b;    // 小于
$a > $b;    // 大于
$a <= $b;   // 小于等于
$a >= $b;   // 大于等于
```

### 逻辑运算符

```php
$a = true;
$b = false;

$a && $b;   // 逻辑与: false
$a || $b;   // 逻辑或: true
!$a;        // 逻辑非: false
```

### 三元运算符

```php
$age = 18;
$status = $age >= 18 ? "成年" : "未成年";
```

### 空合并运算符

```php
$name = $userName ?? "匿名用户";
$config = $userConfig ?? $defaultConfig;
```

### 赋值运算符

```php
$a = 10;
$a += 5;    // $a = 15
$a -= 3;    // $a = 12
$a *= 2;    // $a = 24
$a /= 4;    // $a = 6
$a %= 4;    // $a = 2
```

## 控制结构

### 条件语句

```php
// if 语句
if ($score >= 90) {
    echo "优秀";
} elseif ($score >= 80) {
    echo "良好";
} else {
    echo "需要努力";
}

// 简写形式
if ($isValid) echo "有效";
```

### 循环语句

#### for 循环

```php
// 传统 for 循环
for (int $i = 0; $i < 10; $i++) {
    echo "Count: {$i}\n";
}

// for...in 循环
array $items = ["apple", "banana", "orange"];
for ($item in $items) {
    echo "Item: {$item}\n";
}

// 无限循环
for (;;) {
    // 循环体
    if ($condition) break;
}
```

#### while 循环

```php
int $count = 0;
while ($count < 5) {
    echo "Count: {$count}\n";
    $count++;
}
```

#### foreach 循环

```php
array $fruits = ["apple", "banana", "orange"];
foreach ($fruits as $fruit) {
    echo "I like {$fruit}\n";
}

// 带键的 foreach
array $person = ["name" => "Alice", "age" => 25];
foreach ($person as $key => $value) {
    echo "{$key}: {$value}\n";
}
```

### 分支语句

#### switch 语句

```php
switch ($status) {
    case 200:
        echo "OK";
        break;
    case 404:
        echo "Not Found";
        break;
    default:
        echo "Unknown";
        break;
}

// 不带括号的 switch
switch $day {
    case "Monday":
        echo "Start of week";
        break;
    case "Friday":
        echo "End of week";
        break;
    default:
        echo "Mid week";
        break;
}
```

#### match 语句

```php
$result = match ($value) {
    0 => "zero",
    1 => "one",
    2 => "two",
    default => "many"
};
```

### 跳转语句

```php
// break - 跳出循环
for (int $i = 0; $i < 10; $i++) {
    if ($i == 5) break;
    echo "{$i}\n";
}

// continue - 跳过当前迭代
for (int $i = 0; $i < 10; $i++) {
    if ($i % 2 == 0) continue;
    echo "{$i}\n";
}

// return - 返回函数值
function getValue(): int {
    return 42;
}
```

## 函数

### 函数定义

```php
// 基本函数
function greet(string $name): string {
    return "Hello, {$name}!";
}

// 带默认值的函数
function greet(string $name = "World"): string {
    return "Hello, {$name}!";
}

// 无返回值的函数
function logMessage(string $message): void {
    echo "[LOG] {$message}\n";
}

// 可变参数函数
function sum(...$numbers): int {
    $total = 0;
    foreach ($numbers as $num) {
        $total += $num;
    }
    return $total;
}
```

### 函数调用

```php
// 基本调用
$message = greet("Alice");

// 命名参数调用
$result = calculate(a: 10, b: 20, operation: "add");

// 可变参数调用
$total = sum(1, 2, 3, 4, 5);
```

## 类和对象

### 类定义

```php
class Person {
    // 属性
    private string $name;
    private int $age;

    // 构造函数
    public function __construct(string $name, int $age) {
        $this->name = $name;
        $this->age = $age;
    }

    // 方法
    public function introduce(): string {
        return "I'm {$this->name}, {$this->age} years old.";
    }

    // Getter
    public function getName(): string {
        return $this->name;
    }

    // Setter
    public function setAge(int $age): void {
        $this->age = $age;
    }
}
```

### 对象创建和使用

```php
// 创建对象
$person = new Person("Alice", 25);

// 调用方法
echo $person->introduce();

// 访问属性
$name = $person->getName();
```

### 继承

```php
class Student extends Person {
    private string $school;

    public function __construct(string $name, int $age, string $school) {
        parent::__construct($name, $age);
        $this->school = $school;
    }

    public function getSchool(): string {
        return $this->school;
    }
}
```

### 接口

```php
interface Animal {
    public function cry(): string;
}

class Dog implements Animal {
    public function cry(): string {
        return "汪汪";
    }
}
```

### 类型检查

```php
$dog = new Dog();

// instanceof 检查
if ($dog instanceof Animal) {
    echo "Dog implements Animal";
}

// like 关键字（鸭子类型）
if ($dog like Animal) {
    echo "Dog has Animal-like structure";
}
```

## 异常处理

### try-catch 语句

```php
try {
    $result = riskyOperation();
    echo "Success: {$result}";
} catch (Exception $e) {
    echo "Error: " . $e->getMessage();
} finally {
    cleanup();
}
```

### 抛出异常

```php
function divide(int $a, int $b): float {
    if ($b == 0) {
        throw new Exception("Division by zero");
    }
    return $a / $b;
}
```

## 字符串

### 字符串字面量

```php
$single = 'Hello World';
$double = "Hello World";
$heredoc = <<<EOT
多行
字符串
EOT;
```

### 字符串插值

```php
$name = "Alice";
$message = "Hello, {$name}!";
$result = "Result: @{calculate(10, 20)}";
```

### 字符串方法

```php
$text = "Hello World";

$text->length();           // 获取长度
$text->toUpperCase();      // 转大写
$text->toLowerCase();      // 转小写
$text->trim();             // 去除空白
$text->indexOf("World");   // 查找子串
$text->substring(0, 5);    // 截取子串
$text->replace("World", "Origami"); // 替换
$text->split(" ");         // 分割
```

## 数组

### 数组字面量

```php
// 索引数组
array $numbers = [1, 2, 3, 4, 5];

// 关联数组
array $person = [
    "name" => "Alice",
    "age" => 25
];

// 混合数组
array $mixed = [1, "hello", true, ["nested"]];
```

### 数组方法

```php
$arr = [1, 2, 3, 4, 5];

// 基础操作
$arr->push(6);             // 添加元素
$arr->pop();               // 移除最后一个
$arr->shift();             // 移除第一个
$arr->unshift(0);          // 添加第一个

// 查找
$arr->indexOf(3);          // 查找索引
$arr->includes(3);         // 检查包含
$arr->find(function($n) { return $n > 3; }); // 查找元素

// 迭代
$arr->forEach(function($item) { echo $item; });
$doubled = $arr->map(function($n) { return $n * 2; });
$evens = $arr->filter(function($n) { return $n % 2 == 0; });
$sum = $arr->reduce(function($acc, $n) { return $acc + $n; }, 0);
```

## 命名空间

### 命名空间声明

```php
namespace myapp\utils;

function helper() {
    return "Helper function";
}
```

### 使用命名空间

```php
namespace myapp;

use myapp\utils\helper;

$result = helper();
```

## 注释

### 单行注释

```php
// 这是单行注释
$value = 42; // 行尾注释
```

### 多行注释

```php
/*
这是多行注释
可以跨越多行
*/
```

## 特殊语法

### HTML 内嵌

```php
$title = "Hello World";
$html = <div class="container">
    <h1>{$title}</h1>
    <p>This is embedded HTML</p>
</div>;

echo $html;
```

### 模板字符串

```php
$name = "Alice";
$template = "Hello ${name}, welcome to Origami!";
```

## 最佳实践

### 代码风格

```php
// 使用有意义的变量名
string $userName = "Alice";
int $userAge = 25;

// 使用类型声明
function calculateArea(float $width, float $height): float {
    return $width * $height;
}

// 使用异常处理
try {
    $result = riskyOperation();
} catch (Exception $e) {
    Log::error("Operation failed: " . $e->getMessage());
}
```

### 性能优化

```php
// 避免在循环中重复计算
int $length = $array->length();
for (int $i = 0; $i < $length; $i++) {
    // 循环体
}

// 使用适当的数据结构
array $map = ["key1" => "value1", "key2" => "value2"];
```

## 下一步

- 学习 [标准库](stdlib.md) 使用内置功能
- 了解 [Go 集成](go-integration.md) 扩展语言能力
- 查看 [API 参考](api-reference.md) 获取详细方法说明
