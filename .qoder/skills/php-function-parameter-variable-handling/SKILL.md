---
name: php-function-parameter-variable-handling
description: PHP函数参数和变量列表处理规范。涵盖单个参数、多个参数、引用参数、AST参数等不同类型参数的正确处理方式。用于函数实现和补充时的标准化处理。
---

# PHP函数参数和变量列表处理规范

## 核心原则

在实现PHP函数时，必须正确处理参数列表(GetParams)和变量列表(GetVariables)，确保两者保持一致且完整。

### 重要：错误处理规范

**绝对不能忽略错误返回值！**
```go
// ❌ 错误：忽略错误
param, _ := ctx.GetIndexValue(0)

// ✅ 正确：处理错误
param, exists := ctx.GetIndexValue(0)
if !exists {
    return nil, data.NewErrorThrow(nil, errors.New("缺少必需参数"))
}
```
每个ctx.GetIndexValue()调用都必须检查exists返回值，确保参数存在。

### 构造函数错误处理特别重要
```go
// ✅ 构造函数中的正确错误处理
func (m *MyClassConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 获取必需参数
    iterator, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少必需的迭代器参数"))
    }
    
    // 设置实例状态
    m.instance.field = iterator
    return nil, nil
}
```

### 可选参数的处理
```go
// ✅ 可选参数的安全处理
func (m *MyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 必需参数
    requiredParam, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少必需参数"))
    }
    
    // 可选参数
    optionalParam, exists := ctx.GetIndexValue(1)
    if !exists {
        optionalParam = data.NewStringValue("default_value") // 使用默认值
    }
    
    // 处理逻辑...
    return result, nil
}
```

## 参数类型分类

### 1. 单个参数函数
```go
// 示例：strlen函数
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

### 2. 多个参数函数
```go
// 示例：substr函数
func (f *SubstrFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "str", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "start", 1, nil, data.NewBaseType("int")),
        node.NewParameter(nil, "length", 2, nil, data.NewBaseType("int")),
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

### 3. 引用参数函数（ParameterReference）
```go
// 示例：array_shift函数（引用参数）
func (f *ArrayShiftFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameterReference(nil, "array", 0, data.NewBaseType("array")), // 引用参数
    }
}

func (f *ArrayShiftFunction) GetVariables() []data.Variable {
    return []data.Variable{
        // 注意：引用参数在变量列表中使用普通Variable声明
        node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
    }
}

func (f *ArrayShiftFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    arrayVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少数组参数"))
    }
    
    // 处理引用参数逻辑...
    return result, nil
}
```

### 4. 引用参数函数 - array_pop（修改原数组）
```go
// 示例：array_pop函数（引用参数）
func (f *ArrayPopFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameterReference(nil, "array", 0, data.Mixed{}), // 引用参数
    }
}

func (f *ArrayPopFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "array", 0, data.Mixed{}),
    }
}
```

### 5. AST参数函数（原始AST节点）
```go
// 示例：isset函数（需要原始AST节点）
func (f *IssetFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameterRawAST(nil, "var", 0, data.Mixed{}), // AST参数
    }
}

func (f *IssetFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "var", 0, data.Mixed{}),
    }
}
```

### 6. 可变参数函数（多个值参数）
```go
// 示例：printf函数（可变参数）
func (f *PrintfFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "format", 0, nil, data.NewBaseType("string")),
        node.NewParameters(nil, "values", 1, nil, nil), // 可变参数
    }
}

func (f *PrintfFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
        node.NewVariable(nil, "values", 1, nil), // 可变参数变量
    }
}

func (f *PrintfFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
        node.NewVariable(nil, "values", 1, nil), // 可变参数变量
    }
}
```

## 关键实现规范

### 1. 参数和变量必须一一对应
```go
// ❌ 错误：参数和变量数量不匹配
func (f *WrongFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "param1", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "param2", 1, nil, data.NewBaseType("int")),
    }
}

func (f *WrongFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "param1", 0, data.NewBaseType("string")),
        // 缺少param2的变量声明！
    }
}

// ✅ 正确：参数和变量完全匹配
func (f *CorrectFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "param1", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "param2", 1, nil, data.NewBaseType("int")),
    }
}

func (f *CorrectFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "param1", 0, data.NewBaseType("string")),
        node.NewVariable(nil, "param2", 1, data.NewBaseType("int")),
    }
}
```

### 2. 索引必须连续且正确
```go
// ❌ 错误：索引不连续
func (f *WrongIndexFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "first", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "third", 2, nil, data.NewBaseType("int")), // 跳过了索引1
    }
}

// ✅ 正确：索引连续
func (f *CorrectIndexFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "first", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "second", 1, nil, data.NewBaseType("int")),
        node.NewParameter(nil, "third", 2, nil, data.NewBaseType("bool")),
    }
}
```

### 3. 类型声明要准确
```go
// 根据实际函数需求选择合适的类型
func (f *TypedFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "str_param", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "int_param", 1, nil, data.NewBaseType("int")),
        node.NewParameter(nil, "bool_param", 2, nil, data.NewBaseType("bool")),
        node.NewParameter(nil, "float_param", 3, nil, data.NewBaseType("float")),
        node.NewParameter(nil, "array_param", 4, nil, data.NewBaseType("array")),
        node.NewParameter(nil, "object_param", 5, nil, data.NewBaseType("object")),
    }
}
```

## 特殊情况处理

### 1. 可选参数
```go
// 使用默认值处理可选参数
func (f *OptionalParamFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "required", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "optional", 1, data.NewStringValue("default"), data.NewBaseType("string")),
    }
}
```

### 2. 混合参数类型
```go
// 普通参数 + 可变参数的组合
func (f *MixedParamFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "format", 0, nil, data.NewBaseType("string")),
        node.NewParameters(nil, "args", 1, nil, nil), // 可变参数
    }
}
```

### 3. 引用传递参数
```go
// 需要在文档中明确标注引用语义
func (f *ReferenceFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array")),
        // 该参数会被函数修改（引用传递）
    }
}
```

## 验证清单

实现函数时检查：
- [ ] GetParams和GetVariables返回的参数数量是否一致
- [ ] 参数索引是否从0开始连续递增
- [ ] 参数名称在params和variables中是否完全一致
- [ ] 参数类型声明是否准确反映函数需求
- [ ] 可选参数是否正确设置了默认值
- [ ] 可变参数是否使用了Parameters而不是Parameter
- [ ] 引用参数是否有适当的文档说明
- [ ] Call方法中是否正确处理了所有ctx.GetIndexValue()的错误返回
- [ ] 是否避免使用_忽略错误返回值
- [ ] 构造函数是否正确处理必需参数的缺失情况
- [ ] 错误消息是否清晰描述了问题原因
- [ ] 是否正确导入了errors包用于错误处理