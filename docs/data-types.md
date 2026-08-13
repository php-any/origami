# 数据类型

折言语言支持多种数据类型，每种类型都有其特定的用途和操作方法。

## 基本类型

### 整数 (int)

表示整数值，支持正数、负数和零。

```php
<?php
int $positive = 42;
int $negative = -10;
int $zero = 0;

// 整数运算
int $sum = $positive + $negative;
int $product = $positive * 3;
```

**特点**:

- 无大小限制（受内存限制）
- 支持所有算术运算
- 可以转换为其他数值类型

### 浮点数 (float)

表示小数，用于需要精度的计算。

```php
<?php
float $pi = 3.14159;
float $price = 19.99;
float $percentage = 0.05;

// 浮点运算
float $total = $price * (1 + $percentage);
```

**特点**:

- 支持科学计数法：`1.23e-4`
- 精度有限，不适合精确计算
- 可以转换为整数（会截断小数部分）

### 字符串 (string)

表示文本数据，支持单引号和双引号。

```php
<?php
string $name = "Alice";
string $message = 'Hello, World!';
string $multiline = "This is a
multi-line string";

// 字符串连接
string $greeting = "Hello, " + $name + "!";

// 字符串插值
string $template = "Name: {$name}, Age: {$age}";
```

**特点**:

- 不可变类型
- 支持转义字符：`\n`, `\t`, `\"`, `\'`
- 丰富的字符串方法

### 布尔值 (bool)

表示真或假，用于条件判断。

```php
<?php
bool $isTrue = true;
bool $isFalse = false;
bool $result = 10 > 5; // true

// 布尔运算
bool $and = true && false; // false
bool $or = true || false;  // true
bool $not = !true;         // false
```

**特点**:

- 只有两个值：`true` 和 `false`
- 用于控制流程和逻辑判断
- 支持逻辑运算符

### 空值 (null)

表示"无值"或"未定义"。

```php
<?php
null $empty = null;
string $name = null;

// 空值检查
if ($name === null) {
    echo "Name is not set";
}

// 空值合并运算符
string $displayName = $name ?? "Anonymous";
```

**特点**:

- 表示缺失或未初始化的值
- 可以与任何类型比较
- 支持空值合并运算符 `??`

## 复合类型

### 数组 (array)

表示有序或关联的数据集合。

```php
<?php
// 索引数组
array $numbers = [1, 2, 3, 4, 5];
array $fruits = ["apple", "banana", "orange"];

// 关联数组
array $person = [
    "name" => "Alice",
    "age" => 25,
    "city" => "Beijing"
];

// 混合数组
array $mixed = [
    "first",
    "second" => "value",
    3 => "third"
];
```

**特点**:

- 动态大小
- 支持索引和关联键
- 丰富的数组方法
- 可以嵌套

### 对象 (object)

表示类的实例，具有属性和方法。

```php
<?php
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
}

// 创建对象
object $person = new Person("Alice", 25);
string $intro = $person->introduce();
```

**特点**:

- 通过类定义创建
- 具有封装性
- 支持继承和多态
- 可以调用方法

## 类型转换

### 自动转换

折言语言支持某些自动类型转换：

```php
<?php
// 字符串到数字
string $numStr = "42";
int $num = $numStr + 10; // 52

// 数字到字符串
int $age = 25;
string $message = "Age: " + $age; // "Age: 25"

// 布尔转换
bool $isValid = "hello"; // true (非空字符串)
bool $isEmpty = "";      // false (空字符串)
bool $isZero = 0;        // false (零值)
```

### 显式转换

使用类型转换函数：

```php
<?php
// 字符串转换
string $str = (string) 42;        // "42"
string $str2 = (string) true;     // "1"
string $str3 = (string) null;     // ""

// 数字转换
int $int = (int) "42";            // 42
int $int2 = (int) "42.5";         // 42 (截断)
int $int3 = (int) "hello";        // 0 (无法转换)

float $float = (float) "3.14";    // 3.14
float $float2 = (float) "42";     // 42.0

// 布尔转换
bool $bool = (bool) "hello";      // true
bool $bool2 = (bool) "";          // false
bool $bool3 = (bool) 0;           // false
bool $bool4 = (bool) 42;          // true
```

## 类型检查

### 类型检查函数

```php
<?php
mixed $value = "Hello";

// 检查类型
bool $isString = is_string($value);    // true
bool $isInt = is_int($value);          // false
bool $isArray = is_array($value);      // false
bool $isObject = is_object($value);    // false
bool $isNull = is_null($value);        // false
```

### 类型断言

```php
<?php
mixed $data = "Hello World";

// 安全类型断言
if (is_string($data)) {
    string $text = $data;
    echo $text->length();
}

// 类型转换后检查
string $text = (string) $data;
if ($text !== "") {
    echo "Valid string: {$text}";
}
```

## 类型声明

### 变量类型声明

```php
<?php
// 基本类型声明
int $age = 25;
string $name = "Alice";
bool $isActive = true;
float $price = 19.99;
array $items = [1, 2, 3];
object $user = new User();

// 混合类型（不推荐）
mixed $data = "anything";
```

### 函数参数和返回值类型

```php
<?php
// 带类型声明的函数
function add(int $a, int $b): int {
    return $a + $b;
}

function greet(string $name = "World"): string {
    return "Hello, {$name}!";
}

function processArray(array $items): array {
    return $items->map(function($item) {
        return $item * 2;
    });
}
```

## 最佳实践

### 1. 使用明确的类型声明

```php
<?php
// 好的做法
int $count = 0;
string $message = "";
array $users = [];

// 避免
mixed $data = 0;
```

### 2. 进行类型检查

```php
<?php
function processUser(mixed $user): string {
    if (!is_object($user)) {
        throw new Exception("Expected object");
    }

    if (!method_exists($user, "getName")) {
        throw new Exception("User must have getName method");
    }

    return $user->getName();
}
```

### 3. 使用类型转换函数

```php
<?php
function safeGetString(mixed $value): string {
    if (is_string($value)) {
        return $value;
    }

    if (is_int($value) || is_float($value)) {
        return (string) $value;
    }

    return "";
}
```

### 4. 避免类型混淆

```php
<?php
// 避免混合类型操作
function calculate(int $a, int $b): int {
    return $a + $b; // 明确返回整数
}

// 而不是
function calculate(mixed $a, mixed $b): mixed {
    return $a + $b; // 可能导致意外结果
}
```

## 总结

折言语言提供了丰富的类型系统：

- **基本类型**: int, float, string, bool, null
- **复合类型**: array, object
- **类型转换**: 自动和显式转换
- **类型检查**: 运行时类型验证
- **类型声明**: 编译时类型检查

合理使用类型系统可以：

- 提高代码可读性
- 减少运行时错误
- 改善开发体验
- 提升程序性能
