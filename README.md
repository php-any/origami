# 折言(origami-lang)

折言(origami-lang) 是一门创新性的融合型脚本语言，深度结合 PHP 的快速开发基因与 Go 的高效并发模型。同时还有部分go、ts习惯引入。

## ⚠️ 当前状态

当前未对代码分支进行任何优化，性能尚未优化。
请作为一个工具使用，请勿用于生产环境。

## 🚀 核心特征

### 🎯 Go 反射集成

- **便捷注册**: 一键将 Go 函数注册到脚本域 `vm.RegisterFunction("add", func(a, b int) int { return a + b })`
- **类反射**: 自动将 Go 结构体转换为脚本类 `vm.RegisterReflectClass("User", &User{})`
- **零配置**: 无需手动编写包装代码，自动处理类型转换
- **构造函数**: 支持命名参数 `$user = new User(Name: "Alice")`
- **方法调用**: 直接调用 Go 结构体的公开方法 `$user->SetName("Bob")`

### 语法融合

- **PHP 兼容**: 支持大部分 PHP 语法
- **Go 并发**: `spawn` 关键字启动协程
- **类型系统**: 支持类型声明 `int $i = 0` 和可空类型 `?string`

### 特殊语法

- **HTML 内嵌**: 支持直接内嵌 HTML 代码块
- **字符串插值**: `"Hello {$name}"` 和 `"@{function()}"` 语法
- **鸭子类型**: `like` 关键字进行结构匹配
- **中文编程**: 支持中文关键字 `函数`、`输出` 等
- **参数后置**: 支持 `function($param: type)` 语法
- **异步执行**: `spawn` 关键字启动异步协程
- **泛型类**: 支持 `class DB<T>` 泛型语法

### 数组方法

- **链式调用**: `$array->map()->filter()->reduce()`
- **函数式编程**: `map()`, `filter()`, `reduce()`, `flatMap()`
- **查找方法**: `find()`, `findIndex()`, `includes()`

### 面向对象

- **类继承**: 支持单继承和接口实现
- **类型检查**: `instanceof` 和 `like` 操作符
- **父类访问**: `parent::` 语法

## 📝 示例

### Go 反射集成

```go
// 定义 Go 结构体
type Calculator struct {
    Name string
}

func (c *Calculator) Add(a, b int) int {
    return a + b
}

func (c *Calculator) GetName() string {
    return c.Name
}

// 注册到脚本域
vm.RegisterReflectClass("Calculator", &Calculator{})
```

```php
// 在脚本中使用
$calc = new Calculator(Name: "MyCalc");
echo $calc->GetName();     // 输出: MyCalc
echo $calc->Add(5, 3);     // 输出: 8
```

### 函数注册

```go
// 注册 Go 函数
vm.RegisterFunction("add", func(a, b int) int { return a + b })
vm.RegisterFunction("isEven", func(n int) bool { return n%2 == 0 })
```

```php
// 脚本中调用
$result = add(5, 3);     // 返回 8
$even = isEven(4);       // 返回 true
```

### 基础语法

```php
int $count = 0;
string $name = "World";
echo "Hello {$name}";

function greet(string $name): string {
    return "Hello " . $name;
}
```

### 参数后置语法

```php
function div($obj) {
    return "<div>" + $obj->body + "</div>";
}

function span($obj) {
    return "<span>" + $obj->body + "</span>";
}

$html = div {
    "body": span {
        "body": "内容",
    }
}
```

### 泛型类

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

$list = DB<Users>()->where("name", "张三")->get();
```

### 异步协程

```php
function fetchData($url: string): string {
    // 模拟网络请求
    sleep(1);
    return "Data from " . $url;
}

// 启动异步协程
spawn fetchData("https://api.example.com");

echo "Main thread continues...\n";
```

### HTML 内嵌

```php
$content = <div class="container">
    <h1>{$title}</h1>
    <p>This is embedded HTML</p>
</div>;
```

### 数组操作

```php
$numbers = [1, 2, 3, 4, 5];
$doubled = $numbers->map(($n) => $n * 2);
$evens = $numbers->filter(($n) => $n % 2 == 0);
```

### 中文编程

```php
函数 用户(名称) {
  输出 名称;
}
用户("张三");
```

## 🚀 快速开始

```bash
git clone https://github.com/php-any/origami.git
cd origami
go build -o origami .
./origami script.php
```

## 📚 文档

- [文档](https://github.com/php-any/origami/tree/main/docs)
- [测试用例](https://github.com/php-any/origami/tree/main/tests)

## 💬 讨论群

![折言讨论群二维码](https://github.com/php-any/origami/blob/main/qrcode_1753692981069.jpg)

## 📄 许可证

MIT 许可证
