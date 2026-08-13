# 脚本反射模块

本文档介绍折言语言中的脚本反射功能，允许在脚本中查询和获取类、方法、属性的信息。

## 概述

折言语言提供了 `Reflect` 类，用于在脚本中获取类和对象的结构信息。这个模块可以帮助您：

- 查询脚本中定义的类信息
- 获取类的方法和属性列表
- 查看方法和属性的详细信息
- 动态分析脚本结构

## 使用方法

### 1. 基本用法

```php
<?php
// 定义一个测试类
class MyClass {
    public string $name = "default";
    public int $age = 25;
    public float $score = 95.5;
    public bool $active = true;

    public function getName() {
        return $this->name;
    }

    public function setName($name) {
        $this->name = $name;
    }
}

// 创建反射对象
$reflect = new Reflect();

// 获取类信息
$classInfo = $reflect->getClassInfo("MyClass");
echo "类信息: " . $classInfo . "\n";

// 列出所有方法
$methods = $reflect->listMethods("MyClass");
echo "方法列表: " . $methods . "\n";

// 列出所有属性
$properties = $reflect->listProperties("MyClass");
echo "属性列表: " . $properties . "\n";
?>
```

### 2. 获取类信息

```php
<?php
class TestClass {
    public string $title = "Hello";
    private int $count = 0;

    public function getTitle() {
        return $this->title;
    }

    private function increment() {
        $this->count++;
    }
}

$reflect = new Reflect();

// 获取完整的类信息
$info = $reflect->getClassInfo("TestClass");
echo $info;
// 输出: {"name":"TestClass","methodCount":2,"propertyCount":2,"methods":["getTitle","increment"],"properties":["title","count"]}
?>
```

### 3. 获取方法信息

```php
<?php
class Calculator {
    public function add($a, $b) {
        return $a + $b;
    }

    public function multiply($x, $y) {
        return $x * $y;
    }
}

$reflect = new Reflect();

// 获取特定方法的详细信息
$methodInfo = $reflect->getMethodInfo("Calculator", "add");
echo $methodInfo;
// 输出: {"name":"add","modifier":"public","isStatic":false,"paramCount":2}
?>
```

### 4. 获取属性信息

```php
<?php
class User {
    public string $name = "John";
    public int $age = 30;
    private string $email = "john@example.com";
    public static string $version = "1.0";
}

$reflect = new Reflect();

// 获取特定属性的详细信息
$propertyInfo = $reflect->getPropertyInfo("User", "name");
echo $propertyInfo;
// 输出: {"name":"name","type":"string","modifier":"public","isStatic":false,"defaultValue":"John"}

// 获取所有属性信息
$properties = $reflect->listProperties("User");
echo $properties;
// 输出: [name,age,email,version]
?>
```

## API 参考

### Reflect 类方法

#### getClassInfo(className: string): string

获取指定类的完整信息。

**参数：**

- `className`: 类名

**返回值：**
JSON 格式的字符串，包含：

- `name`: 类名
- `methodCount`: 方法数量
- `propertyCount`: 属性数量
- `methods`: 方法名列表
- `properties`: 属性名列表

#### getMethodInfo(className: string, methodName: string): string

获取指定方法的详细信息。

**参数：**

- `className`: 类名
- `methodName`: 方法名

**返回值：**
JSON 格式的字符串，包含：

- `name`: 方法名
- `modifier`: 访问修饰符（"public", "private", "protected"）
- `isStatic`: 是否为静态方法
- `paramCount`: 参数数量

#### getPropertyInfo(className: string, propertyName: string): string

获取指定属性的详细信息。

**参数：**

- `className`: 类名
- `propertyName`: 属性名

**返回值：**
JSON 格式的字符串，包含：

- `name`: 属性名
- `type`: 属性类型（"string", "int", "float", "bool", "null"）
- `modifier`: 访问修饰符
- `isStatic`: 是否为静态属性
- `defaultValue`: 默认值

#### listClasses(): string

列出所有可用的类。

**返回值：**
逗号分隔的类名列表字符串。

#### listMethods(className: string): string

列出指定类的所有方法。

**参数：**

- `className`: 类名

**返回值：**
逗号分隔的方法名列表字符串。

#### listProperties(className: string): string

列出指定类的所有属性。

**参数：**

- `className`: 类名

**返回值：**
逗号分隔的属性名列表字符串。

## 类型推断

反射模块会自动推断属性的类型：

- **字符串类型**: `string $name = "value"`
- **整数类型**: `int $count = 42`
- **浮点类型**: `float $price = 99.99`
- **布尔类型**: `bool $active = true`
- **空值类型**: `$nullable = null`

## 访问修饰符

支持以下访问修饰符：

- `public`: 公有成员
- `private`: 私有成员
- `protected`: 保护成员

## 注意事项

1. **JSON 格式**: 所有返回的信息都是 JSON 格式的字符串，需要在脚本中解析使用。

2. **类型推断**: 属性类型基于默认值进行推断，如果没有默认值则显示为 "unknown"。

3. **静态成员**: 支持检测静态方法和静态属性。

4. **错误处理**: 如果类或成员不存在，会返回相应的错误信息。

## 示例：动态分析

```php
<?php
class Database {
    public string $host = "localhost";
    public int $port = 3306;
    public bool $connected = false;

    public function connect() {
        $this->connected = true;
    }

    public function disconnect() {
        $this->connected = false;
    }
}

$reflect = new Reflect();

// 动态分析类结构
$classInfo = $reflect->getClassInfo("Database");
$methods = $reflect->listMethods("Database");
$properties = $reflect->listProperties("Database");

echo "数据库类信息:\n";
echo "类信息: " . $classInfo . "\n";
echo "方法: " . $methods . "\n";
echo "属性: " . $properties . "\n";

// 分析特定属性
$hostInfo = $reflect->getPropertyInfo("Database", "host");
echo "host 属性: " . $hostInfo . "\n";
?>
```

这个反射模块为脚本提供了强大的内省能力，让您可以在运行时分析和了解脚本中定义的结构。
