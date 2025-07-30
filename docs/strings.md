# 字符串操作

折言语言提供了丰富的字符串操作方法，支持各种字符串处理需求。

## 字符串字面量

### 基本字符串

```php
<?php
// 单引号字符串
string $single = 'Hello World';

// 双引号字符串
string $double = "Hello World";

// 字符串连接（使用 + 符号）
string $name = "Alice";
string $greeting = "Hello, " + $name + "!";
```

### 字符串插值

```php
<?php
string $name = "Alice";
int $age = 25;

// 变量插值
string $message = "Name: {$name}, Age: {$age}";

// 表达式插值
string $result = "Sum: @{10 + 20}";
```

### 转义字符

```php
<?php
// 常用转义字符
string $newline = "Line 1\nLine 2";
string $tab = "Column1\tColumn2";
string $quote = "He said \"Hello\"";
string $backslash = "Path: C:\\Users\\Name";
```

## 字符串方法

### 长度和查找

#### length()

获取字符串的长度。

```php
<?php
string $text = "Hello World";
int $len = $text->length();  // 11
echo "Length: " + $len + "\n";
```

#### indexOf(search)

查找子字符串在字符串中的位置，如果找不到返回 -1。

```php
<?php
string $text = "Hello World";

int $pos1 = $text->indexOf("World");  // 6
int $pos2 = $text->indexOf("xyz");    // -1
int $pos3 = $text->indexOf("");       // 0

echo "Position of 'World': " + $pos1 + "\n";
echo "Position of 'xyz': " + $pos2 + "\n";
```

### 大小写转换

#### toUpperCase()

将字符串转换为大写。

```php
<?php
string $text = "Hello World";
string $upper = $text->toUpperCase();  // "HELLO WORLD"
echo "Uppercase: " + $upper + "\n";
```

#### toLowerCase()

将字符串转换为小写。

```php
<?php
string $text = "Hello World";
string $lower = $text->toLowerCase();  // "hello world"
echo "Lowercase: " + $lower + "\n";
```

### 字符串截取

#### substring(start, end?)

从字符串中提取子字符串。

```php
<?php
string $text = "Hello World";

// 指定开始和结束位置
string $sub1 = $text->substring(0, 5);   // "Hello"
string $sub2 = $text->substring(6);       // "World"
string $sub3 = $text->substring(0, 0);    // ""

echo "Substring (0, 5): " + $sub1 + "\n";
echo "Substring (6): " + $sub2 + "\n";
```

### 字符串替换

#### replace(search, replace)

替换字符串中的子字符串。

```php
<?php
string $text = "Hello World";

// 替换子字符串
string $result1 = $text->replace("World", "Universe");  // "Hello Universe"
string $result2 = $text->replace("o", "0");            // "Hell0 W0rld"

echo "Replace 'World' with 'Universe': " + $result1 + "\n";
echo "Replace 'o' with '0': " + $result2 + "\n";
```

### 字符串分割

#### split(separator?)

将字符串分割为数组。

```php
<?php
string $text = "Hello World";

// 使用空格分割
array $parts1 = $text->split(" ");  // ["Hello", "World"]

// 使用字符分割
array $parts2 = $text->split("o");  // ["Hell", " W", "rld"]

// 默认分割（按空格）
array $parts3 = $text->split();     // ["Hello", "World"]

echo "Split by space: " + $parts1->join(", ") + "\n";
echo "Split by 'o': " + $parts2->join(", ") + "\n";
```

### 字符串检查

#### startsWith(search)

检查字符串是否以指定前缀开始。

```php
<?php
string $text = "Hello World";

bool $starts1 = $text->startsWith("Hello");  // true
bool $starts2 = $text->startsWith("World");  // false
bool $starts3 = $text->startsWith("");       // true

echo "Starts with 'Hello': " + ($starts1 ? "true" : "false") + "\n";
echo "Starts with 'World': " + ($starts2 ? "true" : "false") + "\n";
```

#### endsWith(search)

检查字符串是否以指定后缀结束。

```php
<?php
string $text = "Hello World";

bool $ends1 = $text->endsWith("World");  // true
bool $ends2 = $text->endsWith("Hello");  // false
bool $ends3 = $text->endsWith("");       // true

echo "Ends with 'World': " + ($ends1 ? "true" : "false") + "\n";
echo "Ends with 'Hello': " + ($ends2 ? "true" : "false") + "\n";
```

### 字符串清理

#### trim()

去除字符串首尾的空白字符。

```php
<?php
string $text = "  Hello World  ";
string $trimmed = $text->trim();  // "Hello World"

echo "Original: '" + $text + "'\n";
echo "Trimmed: '" + $trimmed + "'\n";
```

## 字符串属性

### length

获取字符串的长度（只读属性）。

```php
<?php
string $text = "Hello World";
int $len = $text->length;  // 11

echo "String length: " + $len + "\n";
```

## 实用示例

### 字符串验证

```php
<?php
function validateEmail(string $email): bool {
    // 检查是否包含 @ 符号
    if ($email->indexOf("@") == -1) {
        return false;
    }

    // 检查是否以 .com 结尾
    if (!$email->endsWith(".com")) {
        return false;
    }

    return true;
}

// 测试邮箱验证
string $email1 = "user@example.com";
string $email2 = "invalid-email";

echo "Email 1 valid: " + (validateEmail($email1) ? "true" : "false") + "\n";
echo "Email 2 valid: " + (validateEmail($email2) ? "true" : "false") + "\n";
```

### 字符串格式化

```php
<?php
function formatName(string $firstName, string $lastName): string {
    // 首字母大写
    string $first = $firstName->substring(0, 1)->toUpperCase() +
                   $firstName->substring(1)->toLowerCase();
    string $last = $lastName->substring(0, 1)->toUpperCase() +
                  $lastName->substring(1)->toLowerCase();

    return $first + " " + $last;
}

// 格式化姓名
string $formatted = formatName("john", "doe");
echo "Formatted name: " + $formatted + "\n";  // "John Doe"
```

### 字符串解析

```php
<?php
function parseCSV(string $line): array {
    // 按逗号分割
    array $parts = $line->split(",");

    // 去除每个部分的首尾空白
    array $cleaned = [];
    foreach ($parts as $part) {
        $cleaned[] = $part->trim();
    }

    return $cleaned;
}

// 解析 CSV 行
string $csvLine = "Alice, 25, Beijing ";
array $data = parseCSV($csvLine);
echo "Parsed data: " + $data->join(" | ") + "\n";
```

### 字符串搜索

```php
<?php
function findWords(string $text, string $keyword): array {
    array $words = $text->toLowerCase()->split(" ");
    array $matches = [];

    foreach ($words as $word) {
        if ($word->startsWith($keyword->toLowerCase())) {
            $matches[] = $word;
        }
    }

    return $matches;
}

// 搜索包含特定前缀的单词
string $text = "Hello World How are you";
array $matches = findWords($text, "h");
echo "Words starting with 'h': " + $matches->join(", ") + "\n";
```

## 最佳实践

### 1. 字符串连接

```php
<?php
// 好的做法：使用 + 符号连接字符串
string $name = "Alice";
string $message = "Hello, " + $name + "!";

// 避免：使用 . 符号（这是对象方法调用）
string $message = "Hello, " . $name . "!";  // 错误用法
```

### 2. 字符串比较

```php
<?php
// 好的做法：使用严格比较
string $name = "Alice";
bool $isAlice = $name === "Alice";

// 避免：使用松散比较
bool $isAlice = $name == "Alice";  // 可能有问题
```

### 3. 字符串检查

```php
<?php
// 好的做法：检查字符串是否为空
string $text = "";
if ($text->length() == 0) {
    echo "String is empty\n";
}

// 或者使用空值检查
if ($text === "") {
    echo "String is empty\n";
}
```

### 4. 字符串方法链

```php
<?php
// 好的做法：链式调用方法
string $text = "  Hello World  ";
string $result = $text->trim()->toLowerCase()->replace("world", "universe");
echo "Result: " + $result + "\n";
```

## 常见错误

### 1. 字符串连接错误

```php
<?php
// 错误：使用 . 进行字符串连接
string $result = "Hello" . "World";  // 这是对象方法调用

// 正确：使用 + 进行字符串连接
string $result = "Hello" + "World";
```

### 2. 字符串方法调用错误

```php
<?php
// 错误：调用不存在的方法
string $text = "Hello";
string $result = $text->reverse();  // 方法不存在

// 正确：使用存在的方法
string $result = $text->toUpperCase();
```

### 3. 字符串索引错误

```php
<?php
// 错误：直接使用索引访问字符
string $text = "Hello";
string $char = $text[0];  // 不支持

// 正确：使用 substring 方法
string $char = $text->substring(0, 1);
```

## 总结

折言语言的字符串操作提供了：

- **基本操作**: 长度、查找、截取
- **大小写转换**: toUpperCase(), toLowerCase()
- **字符串替换**: replace()
- **字符串分割**: split()
- **字符串检查**: startsWith(), endsWith()
- **字符串清理**: trim()

合理使用字符串方法可以：

- 提高代码可读性
- 简化字符串处理逻辑
- 减少手动字符串操作错误
- 提升程序性能
