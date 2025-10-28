# HTML 解析器示例

这个示例演示了 Origami HTML 解析器的 if 和 for 特征功能，以及路由系统的使用。

## 功能特性

### 1. if 条件特征

- 支持条件表达式
- 根据条件决定是否渲染 HTML 内容
- 支持嵌套 HTML 结构

### 2. for 循环特征

- 支持数组遍历
- 支持对象遍历
- 支持键值对遍历
- 支持空数组处理

### 3. 路由系统

- 可配置的特征名称
- 动态注册新的属性处理器
- 支持自定义处理器

## 文件说明

- `html_features.zy` - PHP 脚本示例，演示 HTML 特征的使用
- `main.go` - 主程序，用于运行 PHP 脚本
- `go.mod` - Go 模块文件
- `run_example.sh` - 运行脚本

## 运行示例

### 1. 运行 HTML 解析器示例

```bash
go run main.go
```

### 2. 使用运行脚本

```bash
./run_example.sh
```

## 使用示例

### 基本 if 条件

```php
$showDiv = true;
echo "<div if='$showDiv'><h1>条件内容</h1></div>";
```

### 基本 for 循环

```php
$items = ["苹果", "香蕉", "橙子"];
echo "<ul for='item in $items'><li>$item</li></ul>";
```

### 带键值对的 for 循环

```php
$items = ["苹果", "香蕉", "橙子"];
echo "<ul for='index, item in $items'><li>$index: $item</li></ul>";
```

### 对象遍历

```php
$userInfo = ["name" => "张三", "age" => 25];
echo "<dl for='key, value in $userInfo'><dt>$key</dt><dd>$value</dd></dl>";
```

### 嵌套结构

```php
$showDiv = true;
$items = ["苹果", "香蕉"];
echo "<div if='$showDiv'><h2>商品列表</h2><ul for='item in $items'><li>$item</li></ul></div>";
```

## 扩展功能

### 注册自定义处理器

```go
htmlParser.RegisterAttributeHandler("show", func(h *parser.HtmlParser, from *node.TokenFrom, tagName string, attributes map[string]data.GetValue, children []data.GetValue, isSelfClosing bool, attrValue data.GetValue) (data.GetValue, data.Control) {
    // 自定义show逻辑
    return node.NewHtmlNode(from, tagName, attributes, children, isSelfClosing), nil
})
```

### 自定义配置

```go
config := parser.NewHtmlFeatureConfig()
config.IfName = "condition"
config.ForName = "loop"
customParser := parser.NewHtmlParserWithConfig(p, config)
```

## 注意事项

1. HTML 特征属性名称可以通过配置自定义
2. 支持动态注册新的属性处理器
3. 路由系统使得添加新特征变得非常简单
4. 所有特征都支持嵌套 HTML 结构
