# PHP类状态共享模式实例

## 完整示例：SimpleCounter类

```go
package example

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

// SimpleCounterClass 计数器类
type SimpleCounterClass struct {
    node.Node
    count int  // 状态数据存储在类中
}

// NewSimpleCounterClass 创建计数器类
func NewSimpleCounterClass() *SimpleCounterClass {
    return &SimpleCounterClass{
        count: 0,
    }
}

func (c *SimpleCounterClass) GetName() string { return "SimpleCounter" }
func (c *SimpleCounterClass) GetExtend() *string { return nil }
func (c *SimpleCounterClass) GetImplements() []string { return nil }

// GetValue 实现值克隆
func (c *SimpleCounterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
    clone := &SimpleCounterClass{
        count: c.count,  // 复制状态
    }
    return data.NewClassValue(clone, ctx), nil
}

// GetMethod 返回方法（关键：传递实例引用）
func (c *SimpleCounterClass) GetMethod(name string) (data.Method, bool) {
    switch name {
    case "__construct":
        return &SimpleCounterConstruct{instance: c}, true
    case "increment":
        return &SimpleCounterIncrement{instance: c}, true
    case "decrement":
        return &SimpleCounterDecrement{instance: c}, true
    case "getValue":
        return &SimpleCounterGetValue{instance: c}, true
    case "reset":
        return &SimpleCounterReset{instance: c}, true
    }
    return nil, false
}

// 构造函数方法
type SimpleCounterConstruct struct {
    instance *SimpleCounterClass
}

func (m *SimpleCounterConstruct) GetName() string { return "__construct" }
func (m *SimpleCounterConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SimpleCounterConstruct) GetIsStatic() bool { return false }
func (m *SimpleCounterConstruct) GetReturnType() data.Types { return nil }
func (m *SimpleCounterConstruct) GetParams() []data.GetValue { return nil }
func (m *SimpleCounterConstruct) GetVariables() []data.Variable { return nil }

func (m *SimpleCounterConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 初始化状态
    m.instance.count = 0
    return nil, nil
}

// 增加方法
type SimpleCounterIncrement struct {
    instance *SimpleCounterClass
}

func (m *SimpleCounterIncrement) GetName() string { return "increment" }
func (m *SimpleCounterIncrement) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SimpleCounterIncrement) GetIsStatic() bool { return false }
func (m *SimpleCounterIncrement) GetReturnType() data.Types { return data.Int{} }
func (m *SimpleCounterIncrement) GetParams() []data.GetValue { return nil }
func (m *SimpleCounterIncrement) GetVariables() []data.Variable { return nil }

func (m *SimpleCounterIncrement) Call(ctx data.Context) (data.GetValue, data.Control) {
    m.instance.count++  // 通过instance访问状态
    return data.NewIntValue(m.instance.count), nil
}

// 减少方法
type SimpleCounterDecrement struct {
    instance *SimpleCounterClass
}

func (m *SimpleCounterDecrement) Call(ctx data.Context) (data.GetValue, data.Control) {
    m.instance.count--  // 修改状态
    return data.NewIntValue(m.instance.count), nil
}

// 获取值方法
type SimpleCounterGetValue struct {
    instance *SimpleCounterClass
}

func (m *SimpleCounterGetValue) Call(ctx data.Context) (data.GetValue, data.Control) {
    return data.NewIntValue(m.instance.count), nil  // 读取状态
}

// 重置方法
type SimpleCounterReset struct {
    instance *SimpleCounterClass
}

func (m *SimpleCounterReset) Call(ctx data.Context) (data.GetValue, data.Control) {
    m.instance.count = 0  // 重置状态
    return nil, nil
}
```

## 使用示例（PHP代码）

```php
<?php
// 创建计数器实例
$counter1 = new SimpleCounter();
$counter2 = new SimpleCounter();

// 各自独立计数
echo $counter1->increment(); // 输出: 1
echo $counter1->increment(); // 输出: 2
echo $counter2->increment(); // 输出: 1

// 状态隔离验证
echo $counter1->getValue(); // 输出: 2
echo $counter2->getValue(); // 输出: 1

// 重置操作
$counter1->reset();
echo $counter1->getValue(); // 输出: 0
```

## 关键要点说明

1. **状态存储**：`count`字段直接存储在`SimpleCounterClass`中
2. **方法引用**：每个方法结构体都有`instance`字段指向所属实例
3. **状态操作**：所有方法都通过`m.instance.count`访问和修改状态
4. **实例独立**：`counter1`和`counter2`有完全独立的状态
5. **正确复制**：`GetValue`方法正确复制状态并返回包装值

## 对比错误实现

```go
// ❌ 错误实现
type WrongCounterIncrement struct {
    count *int  // 错误：状态不应该在方法中
}

func (m *WrongCounterIncrement) Call(ctx data.Context) (data.GetValue, data.Control) {
    *m.count++  // 错误：方法间状态不共享
    return data.NewIntValue(*m.count), nil
}

// ✅ 正确实现
type CorrectCounterIncrement struct {
    instance *SimpleCounterClass  // 正确：引用类实例
}

func (m *CorrectCounterIncrement) Call(ctx data.Context) (data.GetValue, data.Control) {
    m.instance.count++  // 正确：通过实例访问状态
    return data.NewIntValue(m.instance.count), nil
}
```

这个例子展示了完整的类设计模式，可以帮助理解状态共享的正确实现方式。