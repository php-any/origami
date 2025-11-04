# Origami-lang

Origami-lang is an innovative hybrid scripting language that deeply integrates PHP's rapid development capabilities with Go's efficient concurrency model. It also incorporates some Java and TypeScript conventions.

> [中文文档](README_CN.md) | [English Documentation](README.md)

## ⚠️ Current Status

The codebase has not been optimized yet, and performance is not optimized.
Please use it as a tool, do not use it in production environments.

## 🚀 Core Features

### 🎯 Go Reflection Integration

- **Easy Registration**: Register Go functions to the script domain with one line `vm.RegisterFunction("add", func(a, b int) int { return a + b })`
- **Class Reflection**: Automatically convert Go structs to script classes `vm.RegisterReflectClass("User", &User{})`
- **Zero Configuration**: No manual wrapper code needed, automatic type conversion
- **Constructors**: Support named parameters `$user = new User(Name: "Alice")`
- **Method Calls**: Directly call public methods of Go structs `$user->SetName("Bob")`

### Syntax Fusion

- **PHP Compatibility**: Supports most PHP syntax
- **Go Concurrency**: `spawn` keyword to launch coroutines
- **Type System**: Supports type declarations `int $i = 0` and nullable types `?string`

### Special Syntax

- **HTML Embedding**: Supports direct HTML code blocks
- **String Interpolation**: `"Hello {$name}"` and `"@{function()}"` syntax
- **Duck Typing**: `like` keyword for structural matching
- **Chinese Programming**: Supports Chinese keywords `函数`, `输出`, etc.
- **Postfix Parameters**: Supports `function($param: type)` syntax
- **Async Execution**: `spawn` keyword to launch async coroutines
- **Generic Classes**: Supports `class DB<T>` generic syntax

### Array Methods

- **Chained Calls**: `$array->map()->filter()->reduce()`
- **Functional Programming**: `map()`, `filter()`, `reduce()`, `flatMap()`
- **Search Methods**: `find()`, `findIndex()`, `includes()`

### Object-Oriented

- **Class Inheritance**: Supports single inheritance and interface implementation
- **Type Checking**: `instanceof` and `like` operators
- **Parent Access**: `parent::` syntax

## 📝 Examples

### Go Reflection Integration

```go
// Define Go struct
type Calculator struct {
    Name string
}

func (c *Calculator) Add(a, b int) int {
    return a + b
}

func (c *Calculator) GetName() string {
    return c.Name
}

// Register to script domain
vm.RegisterReflectClass("Calculator", &Calculator{})
```

```php
// Use in script
$calc = new Calculator(Name: "MyCalc");
echo $calc->GetName();     // Output: MyCalc
echo $calc->Add(5, 3);     // Output: 8
```

### Function Registration

```go
// Register Go functions
vm.RegisterFunction("add", func(a, b int) int { return a + b })
vm.RegisterFunction("isEven", func(n int) bool { return n%2 == 0 })
```

```php
// Call in script
$result = add(5, 3);     // Returns 8
$even = isEven(4);       // Returns true
```

### Basic Syntax

```php
int $count = 0;
string $name = "World";
echo "Hello {$name}";

function greet(string $name): string {
    return "Hello " . $name;
}
```

### Postfix Parameter Syntax

```php
function div($obj) {
    return "<div>" + $obj->body + "</div>";
}

function span($obj) {
    return "<span>" + $obj->body + "</span>";
}

$html = div {
    "body": span {
        "body": "Content",
    }
}
```

### Generic Classes

```php
class Users {
    public $name = "";
}

class DB<T> {
    public $where = {};

    public function where($key, $value) {
        $this->where[$key] = $value;
        return $this;
    }

    public function get() {
        return [new T()];
    }
}

$list = DB<Users>()->where("name", "John")->get();
```

### Async Coroutines

```php
function fetchData($url: string): string {
    // Simulate network request
    sleep(1);
    return "Data from " . $url;
}

// Launch async coroutine
spawn fetchData("https://api.example.com");

echo "Main thread continues...\n";
```

### HTML Embedding

```php
$content = <div class="container">
    <h1>{$title}</h1>
    <p>This is embedded HTML</p>
</div>;
```

### Array Operations

```php
$numbers = [1, 2, 3, 4, 5];
$doubled = $numbers->map(($n) => $n * 2);
$evens = $numbers->filter(($n) => $n % 2 == 0);
```

### Chinese Programming

```php
函数 用户(名称) {
  输出 名称;
}
用户("张三");
```

## 🚀 Quick Start

```bash
git clone https://github.com/php-any/origami.git
cd origami
go build -o origami .
./origami script.zy
```

## 📚 Documentation

- [Documentation](https://github.com/php-any/origami/tree/main/docs)
- [Test Cases](https://github.com/php-any/origami/tree/main/tests)

## 💬 Discussion Group

![Origami Discussion Group QR Code](https://github.com/php-any/origami/blob/main/qrcode_1753692981069.jpg)

## 📄 License

MIT License
