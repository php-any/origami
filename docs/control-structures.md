# 控制结构

控制结构是编程语言的核心，用于控制程序的执行流程。折言语言提供了丰富的控制结构。

## 条件语句

### if 语句

最基本的条件判断语句。

```php
<?php
int $age = 18;

// 简单 if 语句
if ($age >= 18) {
    echo "You are an adult.\n";
}

// if-else 语句
if ($age >= 18) {
    echo "You are an adult.\n";
} else {
    echo "You are a minor.\n";
}

// if-elseif-else 语句
if ($age < 13) {
    echo "You are a child.\n";
} elseif ($age < 18) {
    echo "You are a teenager.\n";
} else {
    echo "You are an adult.\n";
}
```

### 嵌套 if 语句

```php
<?php
int $age = 25;
bool $hasLicense = true;

if ($age >= 18) {
    if ($hasLicense) {
        echo "You can drive.\n";
    } else {
        echo "You need a license to drive.\n";
    }
} else {
    echo "You are too young to drive.\n";
}
```

### 条件表达式

```php
<?php
int $score = 85;
string $grade = $score >= 90 ? "A" :
                ($score >= 80 ? "B" :
                ($score >= 70 ? "C" : "D"));

echo "Grade: {$grade}\n";
```

## 循环语句

### for 循环

用于执行固定次数的循环。

```php
<?php
// 基本 for 循环
for (int $i = 0; $i < 5; $i++) {
    echo "Count: {$i}\n";
}

// 递减循环
for (int $i = 10; $i > 0; $i--) {
    echo "Countdown: {$i}\n";
}

// 步长循环
for (int $i = 0; $i <= 20; $i += 2) {
    echo "Even: {$i}\n";
}
```

### while 循环

当条件为真时重复执行。

```php
<?php
int $count = 0;

// 基本 while 循环
while ($count < 5) {
    echo "Count: {$count}\n";
    $count++;
}

// 条件循环
int $number = 1;
while ($number <= 100) {
    if ($number % 2 == 0) {
        echo "Even: {$number}\n";
    }
    $number++;
}
```

### do-while 循环

至少执行一次，然后根据条件决定是否继续。

```php
<?php
int $attempts = 0;

do {
    echo "Attempt: " . ($attempts + 1) . "\n";
    $attempts++;
} while ($attempts < 3);
```

### foreach 循环

遍历数组或集合。

```php
<?php
// 遍历索引数组
array $fruits = ["apple", "banana", "orange"];
foreach ($fruits as $fruit) {
    echo "I like {$fruit}\n";
}

// 遍历关联数组
array $person = [
    "name" => "Alice",
    "age" => 25,
    "city" => "Beijing"
];

foreach ($person as $key => $value) {
    echo "{$key}: {$value}\n";
}

// 遍历带索引的数组
array $colors = ["red", "green", "blue"];
foreach ($colors as $index => $color) {
    echo "Color {$index}: {$color}\n";
}

// 遍历实现Iterator接口的对象
// 更多关于Iterator接口的信息，请参阅[Iterator与foreach](iterator-foreach.md)文档
```

## 分支语句

### switch 语句

多路分支选择。

```php
<?php
string $day = "Monday";

switch ($day) {
    case "Monday":
        echo "Start of work week\n";
        break;
    case "Tuesday":
    case "Wednesday":
    case "Thursday":
        echo "Mid week\n";
        break;
    case "Friday":
        echo "TGIF!\n";
        break;
    case "Saturday":
    case "Sunday":
        echo "Weekend!\n";
        break;
    default:
        echo "Invalid day\n";
}
```

### match 语句

更简洁的多路分支（类似 switch 的表达式版本）。

```php
<?php
int $age = 18;

string $status = match ($age) {
    0, 1, 2 => "baby",
    3, 4, 5 => "toddler",
    6, 7, 8, 9, 10, 11, 12 => "child",
    13, 14, 15, 16, 17 => "teenager",
    18, 19, 20, 21, 22, 23, 24, 25 => "young adult",
    default => "adult"
};

echo "Status: {$status}\n";
```

## 跳转语句

### break 语句

跳出循环或 switch 语句。

```php
<?php
// 跳出循环
for (int $i = 0; $i < 10; $i++) {
    if ($i == 5) {
        break; // 跳出循环
    }
    echo "Number: {$i}\n";
}

// 跳出 switch
string $grade = "B";
switch ($grade) {
    case "A":
        echo "Excellent\n";
        break;
    case "B":
        echo "Good\n";
        break;
    default:
        echo "Need improvement\n";
}
```

### continue 语句

跳过当前循环迭代，继续下一次。

```php
<?php
// 跳过偶数
for (int $i = 1; $i <= 10; $i++) {
    if ($i % 2 == 0) {
        continue; // 跳过偶数
    }
    echo "Odd: {$i}\n";
}

// 跳过特定值
array $numbers = [1, 2, 3, 4, 5, 6];
foreach ($numbers as $num) {
    if ($num == 3) {
        continue; // 跳过 3
    }
    echo "Number: {$num}\n";
}
```

### return 语句

从函数中返回。

```php
<?php
function checkAge(int $age): string {
    if ($age < 0) {
        return "Invalid age";
    }

    if ($age < 18) {
        return "Minor";
    }

    return "Adult";
}

string $result = checkAge(25);
echo "Result: {$result}\n";
```

## 异常处理

### try-catch 语句

处理程序中的异常。

```php
<?php
function divide(int $a, int $b): float {
    if ($b == 0) {
        throw new Exception("Division by zero");
    }
    return $a / $b;
}

// 基本异常处理
try {
    float $result = divide(10, 0);
    echo "Result: {$result}\n";
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}

// 多个 catch 块
try {
    // 可能抛出异常的代码
    float $result = divide(10, 2);
    echo "Result: {$result}\n";
} catch (Exception $e) {
    echo "General error: " . $e->getMessage() . "\n";
} catch (Error $e) {
    echo "System error: " . $e->getMessage() . "\n";
}
```

### finally 块

无论是否发生异常都会执行的代码。

```php
<?php
function processFile(string $filename): string {
    try {
        // 模拟文件操作
        if ($filename == "") {
            throw new Exception("Empty filename");
        }
        return "File processed: {$filename}";
    } catch (Exception $e) {
        echo "Error: " . $e->getMessage() . "\n";
        return "";
    } finally {
        echo "Cleanup completed\n";
    }
}

string $result = processFile("test.txt");
echo "Result: {$result}\n";
```

## 高级控制结构

### 嵌套循环

```php
<?php
// 打印乘法表
for (int $i = 1; $i <= 9; $i++) {
    for (int $j = 1; $j <= $i; $j++) {
        echo "{$j}×{$i}=" . ($i * $j) . " ";
    }
    echo "\n";
}
```

### 循环标签

```php
<?php
// 使用标签跳出外层循环
outer: for (int $i = 0; $i < 3; $i++) {
    for (int $j = 0; $j < 3; $j++) {
        if ($i == 1 && $j == 1) {
            break outer; // 跳出外层循环
        }
        echo "({$i}, {$j}) ";
    }
    echo "\n";
}
```

### 条件循环

```php
<?php
// 用户输入验证
string $input = "";
do {
    echo "Enter 'yes' to continue: ";
    $input = "yes"; // 模拟输入
} while ($input !== "yes");

echo "Thank you!\n";
```

## 最佳实践

### 1. 使用适当的条件结构

```php
<?php
// 好的做法：使用 elseif 而不是嵌套 if
int $score = 85;
if ($score >= 90) {
    echo "A";
} elseif ($score >= 80) {
    echo "B";
} elseif ($score >= 70) {
    echo "C";
} else {
    echo "D";
}

// 避免：过度嵌套
if ($score >= 90) {
    echo "A";
} else {
    if ($score >= 80) {
        echo "B";
    } else {
        if ($score >= 70) {
            echo "C";
        } else {
            echo "D";
        }
    }
}
```

### 2. 循环优化

```php
<?php
// 好的做法：使用 foreach 遍历数组
array $items = [1, 2, 3, 4, 5];
foreach ($items as $item) {
    echo "Item: {$item}\n";
}

// 避免：使用 for 循环遍历数组
for (int $i = 0; $i < count($items); $i++) {
    echo "Item: {$items[$i]}\n";
}
```

### 3. 异常处理

```php
<?php
// 好的做法：具体的异常处理
try {
    // 可能出错的代码
    float $result = divide(10, 0);
} catch (Exception $e) {
    // 记录错误
    echo "Error occurred: " . $e->getMessage() . "\n";
    // 提供默认值或恢复策略
    $result = 0;
}

// 避免：忽略异常
try {
    float $result = divide(10, 0);
} catch (Exception $e) {
    // 空的 catch 块
}
```

### 4. 循环控制

```php
<?php
// 好的做法：使用 break 提前退出
array $numbers = [1, 2, 3, 4, 5];
bool $found = false;

foreach ($numbers as $num) {
    if ($num == 3) {
        $found = true;
        break; // 找到后立即退出
    }
}

// 避免：不必要的循环
foreach ($numbers as $num) {
    if ($num == 3) {
        $found = true;
    }
    // 继续循环，即使已经找到
}
```

## 常见错误

### 1. 无限循环

```php
<?php
// 错误：缺少循环变量更新
int $i = 0;
while ($i < 10) {
    echo "Count: {$i}\n";
    // 忘记 $i++
}

// 正确：确保循环条件会改变
int $i = 0;
while ($i < 10) {
    echo "Count: {$i}\n";
    $i++;
}
```

### 2. 条件判断错误

```php
<?php
// 错误：使用赋值而不是比较
int $age = 18;
if ($age = 20) { // 这是赋值，不是比较
    echo "Age is 20\n";
}

// 正确：使用比较运算符
if ($age == 20) {
    echo "Age is 20\n";
}
```

### 3. 循环边界错误

```php
<?php
// 错误：数组越界
array $items = [1, 2, 3];
for (int $i = 0; $i <= count($items); $i++) { // 应该是 < 而不是 <=
    echo "Item: {$items[$i]}\n";
}

// 正确：使用正确的边界
for (int $i = 0; $i < count($items); $i++) {
    echo "Item: {$items[$i]}\n";
}
```

## 总结

折言语言的控制结构提供了：

- **条件语句**: if, if-else, if-elseif-else
- **循环语句**: for, while, do-while, foreach
- **分支语句**: switch, match
- **跳转语句**: break, continue, return
- **异常处理**: try-catch-finally

合理使用控制结构可以：

- 编写清晰易读的代码
- 避免无限循环和死锁
- 提高程序执行效率
- 增强错误处理能力
