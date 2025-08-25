# 函数参数后置调用语法

折言支持一种特殊的函数调用语法，允许将参数放在函数名后面的花括号或者方括号中，这种语法类似于某些现代编程语言中的块参数语法。

## 花括号形式

### 基本语法

```php
function div($obj) {
    return "<div>" + $obj->body + "</div>";
}

// 传统函数调用
$html = div($config);

// 参数后置调用语法
$html = div {
    "body": "内容"
}
```

### 语法解析

当解析器遇到 `div { ... }` 这种语法时：

1. `div` 被识别为函数名
2. 花括号 `{ ... }` 内的内容被解析为参数
3. 花括号内的内容被解析为对象/字典结构
4. 该对象作为参数传递给 `div` 函数

### 嵌套调用

参数后置调用支持嵌套结构：

```php
function span($obj) {
    return "<span>" + $obj->body + "</span>";
}

$html = div {
    "body": span {
        "body": "内容"
    }
}
```

在这个例子中：
- `span { "body": "内容" }` 创建一个对象，包含 `body` 属性
- 这个对象作为参数传递给 `span` 函数
- `span` 函数的返回值作为 `div` 函数的 `body` 参数

### 复杂参数结构

```php
function createComponent($config) {
    return "<div class='" + $config->class + "'>" +
           "<h1>" + $config->title + "</h1>" +
           "<p>" + $config->content + "</p>" +
           "</div>";
}

$component = createComponent {
    "class": "card",
    "title": "欢迎",
    "content": "这是一个组件"
}
```

### 函数属性调用

```php
function createConfig($options) {
    return {
        "debug": $options->debug ?? false,
        "timeout": $options->timeout ?? 5000
    };
}

$config = createConfig {
    "debug": true,
    "timeout": 10000
}

// 调用配置中的函数
$config->onSuccess("数据加载完成");
```

### 内容解析

#### 对象/字典语法

花括号内的内容被解析为键值对结构：

```php
$obj = {
    "key1": "value1",
    "key2": "value2",
    "number": 123,
    "boolean": true
}
```

#### 支持的数据类型

花括号内支持各种数据类型：

```php
$config = {
    "string": "文本",
    "number": 42,
    "boolean": true,
    "array": [1, 2, 3],
    "function": function($x) { return $x * 2; },
    "nested": {
        "inner": "嵌套值"
    }
}
```

### 实际应用示例

#### HTML 组件构建

```php
namespace components;

function createCard($config) {
    return "<div class='card'>" +
           "<h3>" + $config->title + "</h3>" +
           "<p>" + $config->content + "</p>" +
           "</div>";
}

$card = createCard {
    "title": "欢迎",
    "content": "这是一个卡片组件"
}
```

#### 配置对象创建

```php
function createAppConfig($options) {
    return {
        "debug": $options->debug ?? false,
        "timeout": $options->timeout ?? 5000,
        "retries": $options->retries ?? 3,
        "api": $options->api ?? "https://api.example.com"
    };
}

$appConfig = createAppConfig {
    "debug": true,
    "timeout": 10000,
    "api": "https://myapi.com"
}
```

#### 事件处理器

```php
function createEventHandler($handlers) {
    return {
        "onClick": $handlers->onClick ?? function() {},
        "onSubmit": $handlers->onSubmit ?? function() {},
        "onError": $handlers->onError ?? function($error) {
            echo "错误: " + $error;
        }
    };
}

$handlers = createEventHandler {
    "onClick": function($event) {
        echo "按钮被点击了";
    },
    "onSubmit": function($form) {
        echo "表单提交: " + $form->data;
    }
}
```

#### 数据库查询构建

```php
function createQuery($params) {
    return "SELECT " + $params->fields + 
           " FROM " + $params->table +
           " WHERE " + $params->where;
}

$query = createQuery {
    "fields": "id, name, email",
    "table": "users",
    "where": "status = 'active'"
}
```

### 语法优势

1. **可读性**: 参数结构更清晰，易于理解
2. **简洁性**: 比传统函数调用更简洁
3. **结构化**: 参数以结构化方式组织
4. **嵌套支持**: 支持深层嵌套的参数结构
5. **类型安全**: 支持类型检查和类型提示

### 与传统语法的对比

#### 传统函数调用

```php
// 传统方式
$html = div(array(
    "body" => "内容",
    "class" => "container"
));

// 参数后置调用
$html = div {
    "body": "内容",
    "class": "container"
}
```

#### 复杂参数对比

```php
// 传统方式 - 难以阅读
$component = createComponent(array(
    "class" => "card",
    "title" => "欢迎",
    "content" => "这是一个组件",
    "handlers" => array(
        "onClick" => function($event) { /* ... */ },
        "onSubmit" => function($form) { /* ... */ }
    )
));

// 参数后置调用 - 更清晰
$component = createComponent {
    "class": "card",
    "title": "欢迎", 
    "content": "这是一个组件",
    "handlers": {
        "onClick": function($event) { /* ... */ },
        "onSubmit": function($form) { /* ... */ }
    }
}
```

### 注意事项

1. **函数定义**: 函数必须能够接受对象/字典类型的参数
2. **属性访问**: 在函数内部使用 `->` 操作符访问参数属性
3. **类型检查**: 支持静态类型检查，提供更好的开发体验
4. **作用域**: 花括号内的变量遵循正常的作用域规则
5. **错误处理**: 如果函数不接受对象参数，会抛出类型错误

### 与其他语言的对比

#### JavaScript 对比

```javascript
// JavaScript - 对象字面量作为参数
const html = div({
    body: "内容",
    class: "container"
});

// 折言语法
$html = div {
    "body": "内容",
    "class": "container"
}
```

#### Python 对比

```python
# Python - 字典作为参数
html = div({
    "body": "内容",
    "class": "container"
})

# 折言语法
$html = div {
    "body": "内容",
    "class": "container"
}
```

这种**函数参数后置调用**语法让折言在保持 PHP 兼容性的同时，提供了更现代和直观的函数调用体验，特别适合需要传递复杂参数结构的场景。 

## 方括号形式

### 基本语法

```php
function hello(array $args): string {
    return "Hello ". $args->join("")
}

// 传统函数调用
$str = hello(["w", "o", "r", "l", "d"])

// 参数后置调用语法
$str = hello ["w", "o", "r", "l", "d"]
```

### 语法解析

当解析器遇到`hello [...]`这种语法时：

1. `hello`被识别为函数名
2. 方括号`[...]`内的内容被解析为参数
3. 方括号内的内容被解析为数组结构
4. 该数组作为参数传递给`hello`函数
