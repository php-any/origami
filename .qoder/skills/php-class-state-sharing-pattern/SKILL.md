---
name: php-class-state-sharing-pattern
description: 类声明和函数状态共享的正确模式。确保类状态数据存储在类中，方法通过instance字段访问实例状态。用于避免RecursiveDirectoryIterator等类实现中的常见错误。
---

# 类状态共享正确模式

## 核心原则

类的状态数据应该直接存储在类结构体中，方法通过持有instance引用访问和修改状态。

## 正确的类设计模式

### 1. 类结构体定义
```go
type RecursiveDirectoryIteratorClass struct {
    node.Node
    files []string  // 状态数据直接存储在类中
    pos   int       
    path  string    
}
```

### 2. 方法结构体设计
```go
// 方法结构体必须持有instance字段
type RecursiveDirectoryIteratorConstruct struct {
    instance *RecursiveDirectoryIteratorClass  // 这是正确的设计！
}

func (m *RecursiveDirectoryIteratorConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 通过 m.instance 访问和修改类状态
    m.instance.path = path
    m.instance.pos = 0
    m.instance.files = []string{}
    // ...
}
```

### 3. 方法获取机制
```go
func (r *RecursiveDirectoryIteratorClass) GetMethod(name string) (data.Method, bool) {
    switch name {
    case "__construct":
        // 传递当前实例引用给方法 - 这是关键！
        return &RecursiveDirectoryIteratorConstruct{instance: r}, true
    case "rewind":
        return &RecursiveDirectoryIteratorRewind{instance: r}, true
    // ...
    }
}
```

## GetValue 方法实现规范

### 基本要求
- 必须返回 `data.NewClassValue()` 格式的值
- 确保每次调用返回新的实例（对于需要状态隔离的类）

### 上下文选择指导
```go
// ✅ 推荐模式：创建独立实例副本
func (r *MyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
    clone := &MyClass{
        // 复制所有状态字段
        stateField: r.stateField,
        // ...
    }
    return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}

// ⚠️ 特殊情况：简单数据类可直接返回
func (r *SimpleClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
    return data.NewClassValue(r, ctx), nil
}
```

### 状态克隆规范
```go
// 正确的状态克隆实现
func (r *RecursiveIteratorIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
    // 创建新的实例，每个实例有自己的状态
    clone := &RecursiveIteratorIteratorClass{
        innerIterator: r.innerIterator,
        mode:          r.mode,
        currentKey:    r.currentKey,
        currentValue:  r.currentValue,
        valid:         r.valid,
    }
    return data.NewClassValue(clone, ctx.CreateBaseContext()), nil
}
```

## 常见错误模式（要避免）

### ❌ 错误1：方法不持有实例引用
```go
// 错误：方法无法访问类状态
type RecursiveDirectoryIteratorConstruct struct{}

func (m *RecursiveDirectoryIteratorConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 无法访问类的状态数据！
}
```

### ❌ 错误2：状态数据分散存储
```go
// 错误：状态数据不应该存储在方法中
type RecursiveDirectoryIteratorConstruct struct {
    files []string  // 错误：状态应该在类中
    pos   int       
    path  string    
}
```

### ❌ 错误3：GetValue返回错误格式
```go
// 错误：应该返回包装后的值
func (r *RecursiveDirectoryIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
    clone := &RecursiveDirectoryIteratorClass{...}
    return clone, nil  // 错误：缺少包装
}
```

### ❌ 错误4：忽略错误处理
```go
// 错误：不处理ctx.GetIndexValue的返回值
func (m *MyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
    param, _ := ctx.GetIndexValue(0)  // 错误：忽略错误
    // ...
}
```

## 状态共享原理

1. **状态内聚**：所有状态数据存储在类实例中
2. **方法引用**：每个方法结构体持有指向所属实例的引用
3. **统一访问**：方法通过`m.instance`统一访问状态
4. **状态一致性**：确保所有方法操作同一份状态数据

## 应用场景

当实现以下类型的类时应用此模式：
- 迭代器类（Iterator, RecursiveIterator）
- 数据处理类
- 状态管理类
- 任何需要维护内部状态的类

## 验证清单

实现类时检查：
- [ ] 状态数据是否存储在类结构体中
- [ ] 方法结构体是否持有`instance`字段
- [ ] `GetMethod`是否正确传递实例引用
- [ ] `GetValue`是否返回`data.NewClassValue`格式
- [ ] 是否正确使用`ctx.CreateBaseContext()`创建独立上下文
- [ ] 是否正确克隆所有状态字段
- [ ] `Call`方法中是否正确处理`ctx.GetIndexValue()`的错误返回
- [ ] 方法是否通过`m.instance`访问状态