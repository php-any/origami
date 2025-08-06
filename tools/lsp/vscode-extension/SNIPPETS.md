# Origami 语言代码片段

本扩展为 Origami 语言提供了丰富的代码片段，帮助您快速编写常用的语法结构。

## 使用方法

在 VSCode 中打开 `.cjp` 或 `.origami` 文件，输入以下前缀并按 `Tab` 键或选择建议来插入代码片段：

## 控制流语句

### 条件语句
- `if` - 基本 if 语句
- `ifelse` - if-else 语句  
- `ifelseif` - if-elseif-else 语句

### 循环语句
- `for` - for 循环
- `foreach` - foreach 循环
- `foreachkey` - 带键的 foreach 循环
- `while` - while 循环
- `dowhile` - do-while 循环

### 其他控制流
- `switch` - switch 语句
- `try` - try-catch 块
- `tryfinally` - try-catch-finally 块

## 函数和类

### 函数
- `function` - 函数声明
- `return` - 返回语句

### 类和方法
- `class` - 类声明
- `pubmethod` - 公共方法
- `privmethod` - 私有方法
- `protmethod` - 受保护方法

### 命名空间
- `namespace` - 命名空间声明
- `use` - use 语句

## 变量和数据

### 变量
- `var` - 变量声明
- `const` - 常量声明
- `array` - 数组声明
- `assocarray` - 关联数组声明

### 对象
- `new` - 创建新实例
- `instanceof` - instanceof 检查

## 输出和调试

- `echo` - echo 语句
- `dump` - var_dump 语句

## 异常处理

- `throw` - 抛出异常
- `break` - break 语句
- `continue` - continue 语句

## 示例

输入 `if` 并按 Tab 键：
```php
if (condition) {
    // code
}
```

输入 `foreach` 并按 Tab 键：
```php
foreach ($array as $item) {
    // code
}
```

输入 `class` 并按 Tab 键：
```php
class ClassName {
    public function __construct(parameters) {
        // constructor code
    }

    public function methodName(parameters) {
        // method code
        return value;
    }
}
```

## 提示

- 使用 `Ctrl+Space`（Windows/Linux）或 `Cmd+Space`（Mac）来手动触发代码补全
- 代码片段中的占位符可以使用 `Tab` 键在它们之间跳转
- 某些代码片段包含多个占位符，您可以依次填写它们