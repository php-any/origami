# 运算符

折言语言提供了丰富的运算符，用于执行各种操作和计算。

## 算术运算符

### 基本算术运算符

```php
<?php
int $a = 10;
int $b = 3;

// 加法
int $sum = $a + $b;        // 13

// 减法
int $diff = $a - $b;       // 7

// 乘法
int $product = $a * $b;    // 30

// 除法
float $quotient = $a / $b; // 3.333...

// 取模（余数）
int $remainder = $a % $b;  // 1

// 幂运算
int $power = $a ** $b;     // 1000
```

### 自增自减运算符

```php
<?php
int $x = 5;

// 前缀自增
int $preInc = ++$x;        // $x = 6, $preInc = 6

// 后缀自增
int $postInc = $x++;       // $postInc = 6, $x = 7

// 前缀自减
int $preDec = --$x;        // $x = 6, $preDec = 6

// 后缀自减
int $postDec = $x--;       // $postDec = 6, $x = 5
```

## 比较运算符

### 基本比较运算符

```php
<?php
int $a = 10;
int $b = 5;

// 等于
bool $equal = $a == $b;        // false

// 不等于
bool $notEqual = $a != $b;     // true

// 严格等于（类型和值都相同）
bool $strictEqual = $a === $b; // false

// 严格不等于
bool $strictNotEqual = $a !== $b; // true

// 大于
bool $greater = $a > $b;       // true

// 小于
bool $less = $a < $b;          // false

// 大于等于
bool $greaterEqual = $a >= $b; // true

// 小于等于
bool $lessEqual = $a <= $b;    // false
```

### 字符串比较

```php
<?php
string $str1 = "hello";
string $str2 = "world";

// 字符串比较
bool $strEqual = $str1 == $str2;     // false
bool $strLess = $str1 < $str2;       // true (按字典序)

// 字符串连接
string $combined = $str1 + " " + $str2; // "hello world"
```

## 逻辑运算符

### 基本逻辑运算符

```php
<?php
bool $a = true;
bool $b = false;

// 逻辑与
bool $and = $a && $b;        // false

// 逻辑或
bool $or = $a || $b;         // true

// 逻辑非
bool $not = !$a;             // false

// 逻辑异或
bool $xor = $a xor $b;       // true
```

### 短路求值

```php
<?php
// 短路与：如果第一个为假，不计算第二个
bool $result1 = false && someExpensiveFunction(); // 不会调用函数

// 短路或：如果第一个为真，不计算第二个
bool $result2 = true || someExpensiveFunction();  // 不会调用函数

function someExpensiveFunction(): bool {
    echo "This function is expensive!\n";
    return true;
}
```

## 赋值运算符

### 基本赋值

```php
<?php
int $x = 10;        // 基本赋值

// 复合赋值运算符
$x += 5;            // $x = $x + 5 (15)
$x -= 3;            // $x = $x - 3 (12)
$x *= 2;            // $x = $x * 2 (24)
$x /= 4;            // $x = $x / 4 (6)
$x %= 4;            // $x = $x % 4 (2)
$x **= 3;           // $x = $x ** 3 (8)

// 字符串连接赋值
string $text = "Hello";
$text += " World";  // $text = "Hello World"
```

### 链式赋值

```php
<?php
int $a = $b = $c = 10; // 所有变量都赋值为 10
```

## 位运算符

### 基本位运算

```php
<?php
int $a = 60;  // 二进制: 00111100
int $b = 13;  // 二进制: 00001101

// 按位与
int $and = $a & $b;     // 12 (00001100)

// 按位或
int $or = $a | $b;      // 61 (00111101)

// 按位异或
int $xor = $a ^ $b;     // 49 (00110001)

// 按位非
int $not = ~$a;         // -61 (11000011)

// 左移
int $leftShift = $a << 2;   // 240 (11110000)

// 右移
int $rightShift = $a >> 2;  // 15 (00001111)
```

### 复合位赋值

```php
<?php
int $x = 60;

$x &= 13;    // $x = $x & 13
$x |= 13;    // $x = $x | 13
$x ^= 13;    // $x = $x ^ 13
$x <<= 2;    // $x = $x << 2
$x >>= 2;    // $x = $x >> 2
```

## 空值合并运算符

```php
<?php
// 空值合并运算符 ??
string $name = null;
string $displayName = $name ?? "Anonymous"; // "Anonymous"

string $user = "Alice";
string $greeting = $user ?? "Guest";        // "Alice"

// 链式空值合并
string $result = $a ?? $b ?? $c ?? "default";
```

## 三元运算符

```php
<?php
int $age = 18;

// 基本三元运算符
string $status = $age >= 18 ? "adult" : "minor";

// 嵌套三元运算符
string $category = $age < 13 ? "child" :
                  ($age < 18 ? "teenager" : "adult");

// 与空值合并结合
string $name = $user ?? "Guest";
string $message = $name ? "Hello, {$name}!" : "Hello, Guest!";
```

## 类型运算符

### instanceof 运算符

```php
<?php
class Animal {}
class Dog extends Animal {}

object $animal = new Animal();
object $dog = new Dog();

bool $isAnimal = $animal instanceof Animal;  // true
bool $isDog = $dog instanceof Dog;          // true
bool $isAnimal2 = $dog instanceof Animal;   // true (继承)
bool $isDog2 = $animal instanceof Dog;      // false
```

## 运算符优先级

运算符按以下优先级执行（从高到低）：

1. **最高优先级**

   - `()` 括号
   - `[]` 数组访问
   - `->` 对象成员访问
   - `::` 静态成员访问
   - `.` 对象方法调用（与 `->` 相同）

2. **一元运算符**

   - `++` `--` 自增自减
   - `+` `-` 正负号
   - `!` 逻辑非
   - `~` 按位非
   - `(type)` 类型转换

3. **算术运算符**

   - `**` 幂运算
   - `*` `/` `%` 乘除模
   - `+` `-` 加减
   - `+` 字符串连接

4. **位运算符**

   - `<<` `>>` 位移
   - `&` 按位与
   - `^` 按位异或
   - `|` 按位或

5. **比较运算符**

   - `<` `<=` `>` `>=`
   - `==` `!=` `===` `!==`
   - `instanceof`

6. **逻辑运算符**

   - `&&` 逻辑与
   - `||` 逻辑或
   - `xor` 逻辑异或

7. **赋值运算符**

   - `=` `+=` `-=` `*=` `/=` `%=` `**=`
   - `+=` `&=` `|=` `^=` `<<=` `>>=`

8. **最低优先级**
   - `??` 空值合并
   - `? :` 三元运算符

### 优先级示例

```php
<?php
int $result = 2 + 3 * 4;        // 14 (不是 20)
bool $logical = true && false || true;  // true
int $bitwise = 1 << 2 + 3;      // 32 (不是 8)

// 使用括号明确优先级
int $explicit = (2 + 3) * 4;    // 20
bool $explicit2 = (true && false) || true;  // true
```

### 重要说明

在折言语言中，`.` 符号有两种用途：

1. **字符串连接**：使用 `+` 符号

   ```php
   string $result = "Hello" + " " + "World"; // "Hello World"
   ```

2. **对象方法调用**：使用 `.` 或 `->` 符号（两者相同）
   ```php
   string $length = $text->length();  // 使用 ->
   string $length = $text.length();   // 使用 . （新语法）
   ```

这种设计使得字符串连接和对象方法调用更加清晰和一致。

## 最佳实践

### 1. 使用括号明确优先级

```php
<?php
// 好的做法：使用括号明确意图
int $result = (2 + 3) * 4;

// 避免：依赖隐式优先级
int $result = 2 + 3 * 4; // 可能造成混淆
```

### 2. 避免复杂的表达式

```php
<?php
// 好的做法：分解复杂表达式
int $a = 10;
int $b = 5;
int $c = 3;

int $result = ($a + $b) * $c;

// 避免：一行写太多操作
int $result = $a + $b * $c - $d / $e % $f;
```

### 3. 使用适当的比较运算符

```php
<?php
// 字符串比较
string $name = "Alice";
bool $isAlice = $name === "Alice";  // 严格比较

// 数字比较
int $age = 25;
bool $isAdult = $age >= 18;         // 数值比较

// 空值检查
mixed $data = null;
bool $hasData = $data !== null;     // 明确检查
```

### 4. 注意类型转换

```php
<?php
// 字符串和数字混合运算
string $numStr = "42";
int $num = 10;
int $result = (int) $numStr + $num;  // 52

// 布尔运算
bool $isValid = "hello" && 42;       // true
bool $isEmpty = "" || 0;             // false
```

## 常见错误

### 1. 浮点数精度问题

```php
<?php
// 浮点数比较
float $a = 0.1 + 0.2;
float $b = 0.3;

// 错误的比较方式
bool $equal = $a == $b;  // 可能为 false

// 正确的比较方式
bool $equal = abs($a - $b) < 0.0001;  // 使用容差
```

### 2. 字符串和数字比较

```php
<?php
string $str = "42";
int $num = 42;

// 松散比较
bool $loose = $str == $num;   // true

// 严格比较
bool $strict = $str === $num; // false
```

### 3. 空值处理

```php
<?php
mixed $value = null;

// 错误的检查方式
bool $hasValue = $value == null;  // 可能有问题

// 正确的检查方式
bool $hasValue = $value === null; // 严格检查
```

## 总结

折言语言的运算符系统提供了：

- **算术运算符**: 基本数学运算
- **比较运算符**: 值比较和类型比较
- **逻辑运算符**: 布尔逻辑运算
- **赋值运算符**: 变量赋值和复合赋值
- **位运算符**: 位级操作
- **特殊运算符**: 空值合并、三元运算符等

合理使用运算符可以：

- 编写简洁高效的代码
- 避免类型转换错误
- 提高代码可读性
- 减少运行时错误
