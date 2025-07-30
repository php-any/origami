# 函数

函数是编程中的基本构建块，用于封装可重用的代码。折言语言提供了强大的函数功能。

## 基本函数定义

### 简单函数

```php
<?php
// 基本函数定义
function greet(): void {
    echo "Hello, World!\n";
}

// 调用函数
greet();
```

### 带参数的函数

```php
<?php
// 单个参数
function greet(string $name): void {
    echo "Hello, {$name}!\n";
}

// 多个参数
function add(int $a, int $b): int {
    return $a + $b;
}

// 调用函数
greet("Alice");
int $result = add(10, 20);
echo "Sum: {$result}\n";
```

### 带返回值的函数

```php
<?php
// 返回字符串
function getGreeting(string $name): string {
    return "Hello, {$name}!";
}

// 返回数字
function calculate(int $a, int $b, string $operation): int {
    switch ($operation) {
        case "add":
            return $a + $b;
        case "subtract":
            return $a - $b;
        case "multiply":
            return $a * $b;
        default:
            return 0;
    }
}

// 调用函数
string $message = getGreeting("Bob");
echo $message + "\n";

int $sum = calculate(10, 5, "add");
echo "Result: {$sum}\n";
```

## 参数类型

### 基本类型参数

```php
<?php
// 整数参数
function increment(int $value): int {
    return $value + 1;
}

// 浮点数参数
function calculateArea(float $radius): float {
    return 3.14159 * $radius * $radius;
}

// 字符串参数
function reverse(string $text): string {
    return strrev($text);
}

// 布尔参数
function isValid(bool $flag): string {
    return $flag ? "valid" : "invalid";
}

// 数组参数
function sumArray(array $numbers): int {
    int $sum = 0;
    foreach ($numbers as $num) {
        $sum += $num;
    }
    return $sum;
}
```

### 默认参数

```php
<?php
// 带默认值的参数
function greet(string $name = "World"): string {
    return "Hello, {$name}!";
}

// 多个默认参数
function createUser(string $name, int $age = 18, string $city = "Unknown"): string {
    return "User: {$name}, Age: {$age}, City: {$city}";
}

// 调用函数
echo greet() + "\n";           // "Hello, World!"
echo greet("Alice") + "\n";    // "Hello, Alice!"

echo createUser("Bob") + "\n";                    // "User: Bob, Age: 18, City: Unknown"
echo createUser("Alice", 25) + "\n";             // "User: Alice, Age: 25, City: Unknown"
echo createUser("Charlie", 30, "Beijing") + "\n"; // "User: Charlie, Age: 30, City: Beijing"
```

### 可变参数

```php
<?php
// 可变参数函数
function sum(...$numbers): int {
    int $total = 0;
    foreach ($numbers as $num) {
        $total += $num;
    }
    return $total;
}

// 调用可变参数函数
int $result1 = sum(1, 2, 3);           // 6
int $result2 = sum(1, 2, 3, 4, 5);     // 15
int $result3 = sum();                   // 0

echo "Sum 1: {$result1}\n";
echo "Sum 2: {$result2}\n";
echo "Sum 3: {$result3}\n";
```

## 返回值类型

### 基本类型返回值

```php
<?php
// 返回整数
function getAge(): int {
    return 25;
}

// 返回浮点数
function getPi(): float {
    return 3.14159;
}

// 返回字符串
function getName(): string {
    return "Alice";
}

// 返回布尔值
function isAdult(int $age): bool {
    return $age >= 18;
}

// 返回数组
function getColors(): array {
    return ["red", "green", "blue"];
}
```

### 多返回值

折言支持函数返回多个值，可以使用逗号分隔的语法直接返回多个值，并使用多变量赋值语法接收返回值。

#### 基本语法

```php
<?php
// 函数返回多个值
function getCoordinates(): array {
    return 10, 20;  // 返回两个值
}

function getUserInfo(): array {
    return "Alice", 25, "Beijing";  // 返回三个值
}

// 多变量赋值接收返回值
$a, $b = getCoordinates();
echo "X: {$a}, Y: {$b}\n";  // 输出: X: 10, Y: 20

$name, $age, $city = getUserInfo();
echo "Name: {$name}, Age: {$age}, City: {$city}\n";
```

#### 多返回值函数定义

```php
<?php
// 返回两个值的函数
function divide(int $a, int $b): array {
    if ($b == 0) {
        return 0, "Division by zero";  // 返回错误信息
    }
    return $a / $b, $a % $b;  // 返回商和余数
}

// 返回三个值的函数
function parseDate(string $date): array {
    $parts = explode("-", $date);
    return $parts[0], $parts[1], $parts[2];  // 返回年、月、日
}
```

#### 多变量赋值

```php
<?php
// 接收两个返回值
$quotient, $remainder = divide(10, 3);
echo "Quotient: {$quotient}, Remainder: {$remainder}\n";

// 接收三个返回值
$year, $month, $day = parseDate("2023-12-25");
echo "Date: {$year}-{$month}-{$day}\n";

// 忽略某些返回值（使用下划线）
$result, $_ = divide(15, 4);  // 只接收第一个返回值
echo "Result: {$result}\n";
```

#### 与数组返回值的对比

```php
<?php
// 传统方式：使用数组返回多个值
function divideArray(int $a, int $b): array {
    if ($b == 0) {
        return ["error" => "Division by zero"];
    }
    return ["quotient" => $a / $b, "remainder" => $a % $b];
}

// 多返回值方式
function divideMultiple(int $a, int $b): array {
    if ($b == 0) {
        return 0, "Division by zero";
    }
    return $a / $b, $a % $b;
}

// 使用方式对比
$result = divideArray(10, 3);
echo "Array: {$result['quotient']}, {$result['remainder']}\n";

$quotient, $remainder = divideMultiple(10, 3);
echo "Multiple: {$quotient}, {$remainder}\n";
```

#### 注意事项

1. **返回值数量匹配**: 接收变量的数量应该与函数返回值的数量匹配
2. **忽略返回值**: 使用下划线 `$_` 可以忽略不需要的返回值
3. **类型安全**: 多返回值支持类型检查，提供更好的开发体验
4. **数组形式**: 多返回值在内部被转换为数组形式处理

### void 返回值

```php
<?php
// 不返回值的函数
function printMessage(string $message): void {
    echo $message . "\n";
}

// 调用 void 函数
printMessage("This is a test message");
```

## 函数类型

### 匿名函数（闭包）

```php
<?php
// 定义匿名函数
$greet = function(string $name): string {
    return "Hello, {$name}!";
};

// 调用匿名函数
string $message = $greet("Alice");
echo $message . "\n";

// 匿名函数作为参数
function processArray(array $items, callable $callback): array {
    array $result = [];
    foreach ($items as $item) {
        $result[] = $callback($item);
    }
    return $result;
}

// 使用匿名函数
array $numbers = [1, 2, 3, 4, 5];
array $doubled = processArray($numbers, function($n) {
    return $n * 2;
});

echo "Doubled: " + implode(", ", $doubled) + "\n";
```

### 箭头函数

```php
<?php
// 箭头函数（简化语法）
$square = fn(int $x): int => $x * $x;

// 调用箭头函数
int $result = $square(5);
echo "Square: {$result}\n";

// 箭头函数在数组方法中的使用
array $numbers = [1, 2, 3, 4, 5];
array $squares = $numbers->map(fn($n) => $n * $n);
echo "Squares: " + $squares->join(", ") + "\n";
```

## 作用域和变量

### 局部变量

```php
<?php
function calculate(int $x, int $y): int {
    // 局部变量
    int $sum = $x + $y;
    int $product = $x * $y;

    return $sum + $product;
}

// 全局变量在函数中不可直接访问
int $globalVar = 100;

function testScope(): void {
    // 这里无法直接访问 $globalVar
    // 需要使用 global 关键字或通过参数传递
    echo "Function scope\n";
}
```

### 静态变量

```php
<?php
function counter(): int {
    static int $count = 0;
    $count++;
    return $count;
}

// 调用函数多次
echo "Count: " + counter() + "\n";  // 1
echo "Count: " + counter() + "\n";  // 2
echo "Count: " + counter() + "\n";  // 3
```

## 递归函数

### 基本递归

```php
<?php
// 计算阶乘
function factorial(int $n): int {
    if ($n <= 1) {
        return 1;
    }
    return $n * factorial($n - 1);
}

// 计算斐波那契数列
function fibonacci(int $n): int {
    if ($n <= 1) {
        return $n;
    }
    return fibonacci($n - 1) + fibonacci($n - 2);
}

// 调用递归函数
echo "Factorial of 5: " + factorial(5) + "\n";
echo "Fibonacci(10): " + fibonacci(10) + "\n";
```

### 尾递归优化

```php
<?php
// 尾递归版本的阶乘
function factorialTail(int $n, int $acc = 1): int {
    if ($n <= 1) {
        return $acc;
    }
    return factorialTail($n - 1, $n * $acc);
}

echo "Tail factorial of 5: " + factorialTail(5) + "\n";
```

## 高阶函数

### 函数作为参数

```php
<?php
// 高阶函数：接受函数作为参数
function applyOperation(int $a, int $b, callable $operation): int {
    return $operation($a, $b);
}

// 定义操作函数
function add(int $a, int $b): int {
    return $a + $b;
}

function multiply(int $a, int $b): int {
    return $a * $b;
}

// 使用高阶函数
int $sum = applyOperation(10, 5, add);
int $product = applyOperation(10, 5, multiply);

echo "Sum: " + $sum + "\n";
echo "Product: " + $product + "\n";
```

### 函数返回函数

```php
<?php
// 返回函数的函数
function createMultiplier(int $factor): callable {
    return function(int $x) use ($factor): int {
        return $x * $factor;
    };
}

// 使用返回的函数
$double = createMultiplier(2);
$triple = createMultiplier(3);

echo "Double of 5: " + $double(5) + "\n";
echo "Triple of 5: " + $triple(5) + "\n";
```

## 错误处理

### 异常处理

```php
<?php
function divide(int $a, int $b): float {
    if ($b == 0) {
        throw new Exception("Division by zero");
    }
    return $a / $b;
}

// 调用可能抛出异常的函数
try {
    float $result = divide(10, 0);
    echo "Result: {$result}\n";
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}
```

### 参数验证

```php
<?php
function validateAge(int $age): string {
    if ($age < 0) {
        throw new Exception("Age cannot be negative");
    }

    if ($age > 150) {
        throw new Exception("Age seems unrealistic");
    }

    if ($age < 18) {
        return "minor";
    } elseif ($age < 65) {
        return "adult";
    } else {
        return "senior";
    }
}

try {
    string $status = validateAge(25);
    echo "Status: {$status}\n";
} catch (Exception $e) {
    echo "Validation error: " . $e->getMessage() . "\n";
}
```

## 最佳实践

### 1. 函数命名

```php
<?php
// 好的命名：清晰描述功能
function calculateTotalPrice(float $price, float $tax): float {
    return $price * (1 + $tax);
}

function isValidEmail(string $email): bool {
    return strpos($email, "@") !== false;
}

// 避免：模糊的命名
function calc(float $p, float $t): float {
    return $p * (1 + $t);
}
```

### 2. 单一职责

```php
<?php
// 好的做法：每个函数只做一件事
function validateEmail(string $email): bool {
    return strpos($email, "@") !== false && strpos($email, ".") !== false;
}

function sendEmail(string $email, string $message): bool {
    // 发送邮件的逻辑
    return true;
}

// 避免：一个函数做太多事
function processUser(string $email, string $message): bool {
    // 验证邮箱
    if (strpos($email, "@") === false) {
        return false;
    }

    // 发送邮件
    // 记录日志
    // 更新数据库
    // 等等...
    return true;
}
```

### 3. 参数验证

```php
<?php
// 好的做法：在函数开始处验证参数
function calculateArea(float $width, float $height): float {
    if ($width <= 0 || $height <= 0) {
        throw new Exception("Width and height must be positive");
    }

    return $width * $height;
}

// 避免：不验证参数
function calculateArea(float $width, float $height): float {
    return $width * $height; // 可能返回负数
}
```

### 4. 返回值一致性

```php
<?php
// 好的做法：一致的返回值类型
function findUser(int $id): mixed {
    if ($id <= 0) {
        return null;
    }

    // 查找用户的逻辑
    if ($id == 1) {
        return ["id" => 1, "name" => "Alice"];
    }

    return null; // 始终返回 null 表示未找到
}

// 避免：不一致的返回值
function findUser(int $id): mixed {
    if ($id <= 0) {
        return false; // 返回 false
    }

    if ($id == 1) {
        return ["id" => 1, "name" => "Alice"];
    }

    return "not found"; // 返回字符串
}
```

## 常见错误

### 1. 忘记返回值

```php
<?php
// 错误：函数声明返回 int 但没有返回值
function add(int $a, int $b): int {
    $result = $a + $b;
    // 忘记 return $result;
}

// 正确：确保有返回值
function add(int $a, int $b): int {
    return $a + $b;
}
```

### 2. 参数类型不匹配

```php
<?php
// 错误：传递错误类型的参数
function greet(string $name): string {
    return "Hello, {$name}!";
}

greet(123); // 传递 int 而不是 string

// 正确：确保参数类型匹配
greet("123"); // 传递字符串
```

### 3. 无限递归

```php
<?php
// 错误：没有递归终止条件
function infiniteRecursion(int $n): int {
    return infiniteRecursion($n + 1); // 永远不会停止
}

// 正确：有终止条件
function factorial(int $n): int {
    if ($n <= 1) {
        return 1; // 终止条件
    }
    return $n * factorial($n - 1);
}
```

## 总结

折言语言的函数系统提供了：

- **基本函数**: 带参数和返回值的函数
- **默认参数**: 为参数提供默认值
- **可变参数**: 接受不定数量的参数
- **匿名函数**: 函数式编程支持
- **递归函数**: 自我调用的函数
- **高阶函数**: 函数作为参数或返回值

合理使用函数可以：

- 提高代码复用性
- 增强代码可读性
- 简化复杂逻辑
- 便于测试和维护
