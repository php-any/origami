# PHP函数参数处理详细参考

## Node包参数构建函数

### node.NewParameter
创建单个固定参数：
```go
func NewParameter(from data.From, name string, index int, defaultValue data.GetValue, typ data.Types) *ParameterNode
```

**参数说明：**
- `from`: AST来源位置信息
- `name`: 参数名称
- `index`: 参数索引（从0开始）
- `defaultValue`: 默认值（可为nil表示必填参数）
- `typ`: 参数类型

### node.NewParameterReference
创建引用参数：
```go
func NewParameterReference(from data.From, name string, index int, ty data.Types) data.Parameter
```

**使用场景：**
- array_shift、array_pop等需要修改原数组的函数
- sort、ksort等需要就地排序的函数
- 任何需要通过引用修改传入变量的函数

### node.NewParameters
创建可变参数（多个值参数）：
```go
func NewParameters(from data.From, name string, index int, defaultValue data.GetValue, typ data.Types) *ParametersNode
```

**使用场景：**
- printf、sprintf等格式化函数
- array_merge等可变参数函数
- 任意接受不定数量参数的函数

### node.NewParameterRawAST
创建AST参数（原始AST节点）：
```go
func NewParameterRawAST(from data.From, name string, index int, ty data.Types) data.GetValue
```

**使用场景：**
- isset、empty等需要检查变量存在的函数
- eval等需要访问原始代码结构的函数
- 任意需要处理AST节点而非求值结果的函数

### node.NewVariable
创建变量声明：
```go
func NewVariable(from data.From, name string, index int, typ data.Types) *VariableNode
```

## 完整实现示例

### 1. 简单函数（单参数）
```go
type StrlenFunction struct{}

func (f *StrlenFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "str", 0, nil, data.NewBaseType("string")),
    }
}

func (f *StrlenFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "str", 0, data.NewBaseType("string")),
    }
}
```

### 2. 复杂函数（多参数+可选参数）
```go
type SubstrFunction struct{}

func (f *SubstrFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "str", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "start", 1, nil, data.NewBaseType("int")),
        node.NewParameter(nil, "length", 2, nil, data.NewBaseType("int")), // 可选
    }
}

func (f *SubstrFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "str", 0, data.NewBaseType("string")),
        node.NewVariable(nil, "start", 1, data.NewBaseType("int")),
        node.NewVariable(nil, "length", 2, data.NewBaseType("int")),
    }
}
```

### 3. 可变参数函数
```go
type PrintfFunction struct{}

func (f *PrintfFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "format", 0, nil, data.NewBaseType("string")),
        node.NewParameters(nil, "values", 1, nil, nil), // 可变参数
    }
}

func (f *PrintfFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
        node.NewVariable(nil, "values", 1, nil), // 注意类型为nil
    }
}
```

### 4. 引用参数函数
```go
type ArrayPopFunction struct{}

func (f *ArrayPopFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array")),
    }
}

func (f *ArrayPopFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
    }
}

// 在Call方法中处理引用传参
func (f *ArrayPopFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    arrayVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少数组参数"))
    }
    
    if array, ok := arrayVal.(*data.ArrayValue); ok {
        if len(array.Items) > 0 {
            // 引用修改：直接修改原数组
            lastItem := array.Items[len(array.Items)-1]
            array.Items = array.Items[:len(array.Items)-1]
            return lastItem, nil
        }
    }
    return data.NewNullValue(), nil
}
```

## 常见类型映射

| PHP类型 | Go中data.Types表示 |
|---------|-------------------|
| string | data.NewBaseType("string") |
| int | data.NewBaseType("int") |
| bool | data.NewBaseType("bool") |
| float | data.NewBaseType("float") |
| array | data.NewBaseType("array") |
| object | data.NewBaseType("object") |
| mixed | nil 或 data.Mixed{} |
| callable | data.NewBaseType("callable") |

## 错误模式和解决方案

### 1. 参数数量不匹配
```go
// ❌ 错误
func (f *WrongFunction) GetParams() []data.GetValue {
    return []data.GetValue{ /* 3个参数 */ }
}
func (f *WrongFunction) GetVariables() []data.Variable {
    return []data.Variable{ /* 只有2个变量 */ }
}

// ✅ 正确
func (f *CorrectFunction) GetParams() []data.GetValue {
    return []data.GetValue{ /* 3个参数 */ }
}
func (f *CorrectFunction) GetVariables() []data.Variable {
    return []data.Variable{ /* 3个变量，与参数一一对应 */ }
}
```

### 2. 索引混乱
```go
// ❌ 错误
func (f *WrongIndexFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "first", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "second", 2, nil, data.NewBaseType("int")), // 索引跳跃
    }
}

// ✅ 正确
func (f *CorrectIndexFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "first", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "second", 1, nil, data.NewBaseType("int")),
    }
}
```

### 3. 可变参数处理错误
```go
// ❌ 错误：使用Parameter处理可变参数
func (f *WrongVariadicFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "format", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "args", 1, nil, nil), // 错误！
    }
}

// ✅ 正确：使用Parameters处理可变参数
func (f *CorrectVariadicFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "format", 0, nil, data.NewBaseType("string")),
        node.NewParameters(nil, "args", 1, nil, nil),
    }
}
```

## 调试技巧

1. **参数验证**：在Call方法开始处打印所有参数
2. **类型检查**：使用类型断言验证参数类型
3. **索引检查**：确保ctx.GetIndexValue的索引与声明一致
4. **默认值处理**：正确处理可选参数的默认值逻辑
5. **错误处理**：**绝不能使用_忽略ctx.GetIndexValue的第二个返回值**

### 关键错误处理原则

```go
// ❌ 绝对禁止的写法
value, _ := ctx.GetIndexValue(0)

// ✅ 正确的写法
value, exists := ctx.GetIndexValue(0)
if !exists {
    return nil, data.NewErrorThrow(nil, errors.New("参数不存在"))
}
```

每个ctx.GetIndexValue()调用都必须检查exists返回值！