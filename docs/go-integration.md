# Go 集成指南

本文档详细介绍如何将 Go 函数和结构体集成到折言脚本中，扩展语言的功能。

## 概述

折言语言提供了强大的 Go 集成能力，允许您：

- 将 Go 函数注册为脚本函数
- 将 Go 结构体注册为脚本类
- 在脚本中调用 Go 代码
- 实现高性能的扩展功能

## 集成方式

### 1. 函数集成

#### 基本函数注册

```go
package myext

import (
    "fmt"
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

// 定义函数结构体
type MyFunction struct{}

// 实现 FuncStmt 接口
func (f *MyFunction) GetName() string {
    return "myFunction"
}

func (f *MyFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "message", 0, nil, nil),
    }
}

func (f *MyFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "message", 0, nil),
    }
}

// 实现函数调用逻辑
func (f *MyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 获取参数
    params := f.GetParams()
    if len(params) > 0 {
        arg, _ := params[0].GetValue(ctx)
        if msg, ok := arg.(data.AsString); ok {
            fmt.Println("Go function called with:", msg.AsString())
        }
    }

    return nil, nil
}

// 创建函数实例
func NewMyFunction() data.FuncStmt {
    return &MyFunction{}
}
```

#### 带返回值的函数

```go
type MathFunction struct{}

func (f *MathFunction) GetName() string {
    return "add"
}

func (f *MathFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "a", 0, nil, nil),
        node.NewParameter(nil, "b", 0, nil, nil),
    }
}

func (f *MathFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "a", 0, nil),
        node.NewVariable(nil, "b", 0, nil),
    }
}

func (f *MathFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 获取参数
    params := f.GetParams()
    if len(params) >= 2 {
        a, _ := params[0].GetValue(ctx)
        b, _ := params[1].GetValue(ctx)

        // 类型转换和计算
        if aVal, ok := a.(data.AsInt); ok {
            if bVal, ok := b.(data.AsInt); ok {
                result := aVal.AsInt() + bVal.AsInt()
                return data.NewIntValue(result), nil
            }
        }
    }

    return data.NewIntValue(0), nil
}
```

### 2. 类集成

#### 基本类注册

```go
package myext

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

// 定义 Go 结构体
type MyStruct struct {
    Name string
    Age  int
}

// 定义类结构体
type MyClass struct {
    node.Node
    instance *MyStruct
}

// 实现 ClassStmt 接口
func (c *MyClass) GetName() string {
    return "MyClass"
}

func (c *MyClass) GetExtend() *string {
    return nil
}

func (c *MyClass) GetImplements() []string {
    return nil
}

func (c *MyClass) GetProperty(name string) (data.Property, bool) {
    switch name {
    case "name":
        return &MyProperty{name: "name", value: c.instance.Name}, true
    case "age":
        return &MyProperty{name: "age", value: c.instance.Age}, true
    }
    return nil, false
}

func (c *MyClass) GetProperties() map[string]data.Property {
    return map[string]data.Property{
        "name": &MyProperty{name: "name", value: c.instance.Name},
        "age":  &MyProperty{name: "age", value: c.instance.Age},
    }
}

func (c *MyClass) GetMethod(name string) (data.Method, bool) {
    switch name {
    case "greet":
        return &MyMethod{instance: c.instance}, true
    }
    return nil, false
}

func (c *MyClass) GetMethods() []data.Method {
    return []data.Method{
        &MyMethod{instance: c.instance},
    }
}

func (c *MyClass) GetConstruct() data.Method {
    return &MyConstructor{}
}

// 属性实现
type MyProperty struct {
    name  string
    value interface{}
}

func (p *MyProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
    switch v := p.value.(type) {
    case string:
        return data.NewStringValue(v), nil
    case int:
        return data.NewIntValue(v), nil
    }
    return nil, nil
}

// 方法实现
type MyMethod struct {
    instance *MyStruct
}

func (m *MyMethod) GetName() string {
    return "greet"
}

func (m *MyMethod) Call(ctx data.Context, args ...data.GetValue) (data.GetValue, data.Control) {
    message := fmt.Sprintf("Hello, I'm %s, %d years old", m.instance.Name, m.instance.Age)
    return data.NewStringValue(message), nil
}

// 构造函数
type MyConstructor struct{}

func (c *MyConstructor) GetName() string {
    return "__construct"
}

func (c *MyConstructor) Call(ctx data.Context, args ...data.GetValue) (data.GetValue, data.Control) {
    // 创建实例
    instance := &MyStruct{}

    // 设置参数
    if len(args) >= 1 {
        if name, ok := args[0].(data.AsString); ok {
            instance.Name = name.AsString()
        }
    }

    if len(args) >= 2 {
        if age, ok := args[1].(data.AsInt); ok {
            instance.Age = age.AsInt()
        }
    }

    // 返回类实例
    return &MyClass{instance: instance}, nil
}

// 创建类实例
func NewMyClass() data.ClassStmt {
    return &MyClass{}
}
```

### 3. 注册到虚拟机

#### 在 main.go 中注册

```go
package main

import (
    "github.com/php-any/origami/runtime"
    "github.com/php-any/origami/std"
    "myext" // 你的扩展包
)

func main() {
    // 创建虚拟机
    vm := runtime.NewVM(parser)

    // 加载标准库
    std.Load(vm)

    // 注册自定义函数
    vm.AddFunc(myext.NewMyFunction())
    vm.AddFunc(myext.NewMathFunction())

    // 注册自定义类
    vm.AddClass(myext.NewMyClass())

    // 运行脚本
    vm.LoadAndRun("script.zy")
}
```

## 高级集成示例

### 1. HTTP 客户端集成

```go
package httpext

import (
    "net/http"
    "io/ioutil"
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type HTTPClient struct {
    client *http.Client
}

type HTTPGetFunction struct {
    client *HTTPClient
}

func (f *HTTPGetFunction) GetName() string {
    return "httpGet"
}

func (f *HTTPGetFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "url", 0, nil, nil),
    }
}

func (f *HTTPGetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    params := f.GetParams()
    if len(params) > 0 {
        arg, _ := params[0].GetValue(ctx)
        if url, ok := arg.(data.AsString); ok {
            resp, err := f.client.client.Get(url.AsString())
            if err != nil {
                return data.NewStringValue(""), nil
            }
            defer resp.Body.Close()

            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                return data.NewStringValue(""), nil
            }

            return data.NewStringValue(string(body)), nil
        }
    }
    return data.NewStringValue(""), nil
}
```

### 2. 数据库集成

```go
package dbext

import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type DatabaseClass struct {
    node.Node
    db *sql.DB
}

type QueryFunction struct {
    db *sql.DB
}

func (f *QueryFunction) GetName() string {
    return "query"
}

func (f *QueryFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "sql", 0, nil, nil),
    }
}

func (f *QueryFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    params := f.GetParams()
    if len(params) > 0 {
        arg, _ := params[0].GetValue(ctx)
        if sqlStr, ok := arg.(data.AsString); ok {
            rows, err := f.db.Query(sqlStr.AsString())
            if err != nil {
                return data.NewStringValue(""), nil
            }
            defer rows.Close()

            // 处理查询结果
            var results []string
            for rows.Next() {
                var name string
                rows.Scan(&name)
                results = append(results, name)
            }

            return data.NewStringValue(strings.Join(results, ",")), nil
        }
    }
    return data.NewStringValue(""), nil
}
```

### 3. 文件系统集成

```go
package fsext

import (
    "os"
    "io/ioutil"
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type ReadFileFunction struct{}

func (f *ReadFileFunction) GetName() string {
    return "readFile"
}

func (f *ReadFileFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "path", 0, nil, nil),
    }
}

func (f *ReadFileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    params := f.GetParams()
    if len(params) > 0 {
        arg, _ := params[0].GetValue(ctx)
        if path, ok := arg.(data.AsString); ok {
            content, err := ioutil.ReadFile(path.AsString())
            if err != nil {
                return data.NewStringValue(""), nil
            }
            return data.NewStringValue(string(content)), nil
        }
    }
    return data.NewStringValue(""), nil
}

type WriteFileFunction struct{}

func (f *WriteFileFunction) GetName() string {
    return "writeFile"
}

func (f *WriteFileFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "path", 0, nil, nil),
        node.NewParameter(nil, "content", 0, nil, nil),
    }
}

func (f *WriteFileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    params := f.GetParams()
    if len(params) >= 2 {
        pathArg, _ := params[0].GetValue(ctx)
        contentArg, _ := params[1].GetValue(ctx)

        if path, ok := pathArg.(data.AsString); ok {
            if content, ok := contentArg.(data.AsString); ok {
                err := ioutil.WriteFile(path.AsString(), []byte(content.AsString()), 0644)
                if err != nil {
                    return data.NewBoolValue(false), nil
                }
                return data.NewBoolValue(true), nil
            }
        }
    }
    return data.NewBoolValue(false), nil
}
```

## 最佳实践

### 1. 错误处理

```go
type SafeFunction struct{}

func (f *SafeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    defer func() {
        if r := recover(); r != nil {
            // 记录错误日志
            fmt.Printf("Function panic: %v\n", r)
        }
    }()

    // 函数逻辑
    return nil, nil
}
```

### 2. 类型安全

```go
func safeGetString(value data.GetValue, ctx data.Context) (string, bool) {
    if value == nil {
        return "", false
    }

    val, _ := value.GetValue(ctx)
    if str, ok := val.(data.AsString); ok {
        return str.AsString(), true
    }
    return "", false
}

func safeGetInt(value data.GetValue, ctx data.Context) (int, bool) {
    if value == nil {
        return 0, false
    }

    val, _ := value.GetValue(ctx)
    if num, ok := val.(data.AsInt); ok {
        return num.AsInt(), true
    }
    return 0, false
}
```

### 3. 性能优化

```go
// 缓存常用值
type CachedFunction struct {
    cache map[string]interface{}
}

func (f *CachedFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 使用缓存避免重复计算
    key := "cache_key"
    if cached, exists := f.cache[key]; exists {
        return data.NewStringValue(cached.(string)), nil
    }

    // 计算结果并缓存
    result := "computed_value"
    f.cache[key] = result
    return data.NewStringValue(result), nil
}
```

## 调试技巧

### 1. 日志记录

```go
type DebugFunction struct{}

func (f *DebugFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    fmt.Printf("Debug: Function called with context %+v\n", ctx)

    // 获取所有参数
    params := f.GetParams()
    for i, param := range params {
        val, _ := param.GetValue(ctx)
        fmt.Printf("Param %d: %+v\n", i, val)
    }

    return nil, nil
}
```

### 2. 参数验证

```go
func validateParams(params []data.GetValue, ctx data.Context) error {
    if len(params) == 0 {
        return fmt.Errorf("expected at least one parameter")
    }

    for i, param := range params {
        if param == nil {
            return fmt.Errorf("parameter %d is nil", i)
        }
    }

    return nil
}
```

## 常见问题

### 1. 类型转换失败

**问题**: 参数类型转换失败
**解决方案**: 使用类型断言和默认值

```go
func (f *MyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    params := f.GetParams()
    if len(params) > 0 {
        arg, _ := params[0].GetValue(ctx)

        // 安全的类型转换
        var result string
        if str, ok := arg.(data.AsString); ok {
            result = str.AsString()
        } else if num, ok := arg.(data.AsInt); ok {
            result = fmt.Sprintf("%d", num.AsInt())
        } else {
            result = "unknown"
        }

        return data.NewStringValue(result), nil
    }
    return data.NewStringValue(""), nil
}
```

### 2. 内存泄漏

**问题**: 长时间运行的程序内存泄漏
**解决方案**: 及时清理资源

```go
type ResourceFunction struct {
    resources map[string]interface{}
}

func (f *ResourceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 使用资源
    resource := "some_resource"
    f.resources[resource] = "value"

    // 确保在函数结束时清理
    defer func() {
        delete(f.resources, resource)
    }()

    return nil, nil
}
```

## 总结

通过 Go 集成，您可以：

1. **扩展功能**: 将 Go 的强大功能引入折言脚本
2. **提高性能**: 关键路径使用 Go 代码优化性能
3. **复用代码**: 利用现有的 Go 库和工具
4. **类型安全**: 利用 Go 的类型系统保证代码质量

记住以下要点：

- 始终进行类型检查和错误处理
- 使用适当的接口实现
- 注意内存管理和资源清理
- 编写清晰的文档和测试

通过遵循这些最佳实践，您可以创建强大、可靠且易于维护的 Go 扩展。
