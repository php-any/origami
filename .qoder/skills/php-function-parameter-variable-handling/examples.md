# PHP函数参数处理实例

## 实际函数实现示例

### 1. 字符串处理函数 - strlen
```go
package stringfuncs

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type StrlenFunction struct{}

func (f *StrlenFunction) GetName() string { return "strlen" }
func (f *StrlenFunction) GetIsStatic() bool { return true }

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

func (f *StrlenFunction) GetReturnType() data.Types {
    return data.NewBaseType("int")
}

func (f *StrlenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    strVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少字符串参数"))
    }
    
    if str, ok := strVal.(data.AsString); ok {
        return data.NewIntValue(len(str.AsString())), nil
    }
    return data.NewIntValue(0), nil
}
```

### 2. 引用参数函数 - array_shift（修改原数组）
```go
package arrayfuncs

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type ArrayShiftFunction struct{}

func (f *ArrayShiftFunction) GetName() string { return "array_shift" }
func (f *ArrayShiftFunction) GetIsStatic() bool { return true }

func (f *ArrayShiftFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameterReference(nil, "array", 0, data.NewBaseType("array")), // 引用参数
    }
}

func (f *ArrayShiftFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
    }
}

func (f *ArrayShiftFunction) GetReturnType() data.Types {
    return data.NewMixedType()
}

func (f *ArrayShiftFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    arrayVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少数组参数"))
    }
    
    if array, ok := arrayVal.(*data.ArrayValue); ok {
        if len(array.List) == 0 {
            return data.NewNullValue(), nil
        }
        
        // 移除并返回第一个元素
        firstElement := array.List[0]
        array.List = array.List[1:]
        
        return firstElement.Value, nil
    }
    
    return data.NewNullValue(), nil
}
```

### 3. 多个参数函数 - substr
```go
package arrayfuncs

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type ArrayPushFunction struct{}

func (f *ArrayPushFunction) GetName() string { return "array_push" }
func (f *ArrayPushFunction) GetIsStatic() bool { return true }

func (f *ArrayPushFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array")),
        node.NewParameters(nil, "values", 1, nil, nil),
    }
}

func (f *ArrayPushFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
        node.NewVariable(nil, "values", 1, nil),
    }
}

func (f *ArrayPushFunction) GetReturnType() data.Types {
    return data.NewBaseType("int")
}

func (f *ArrayPushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    arrayVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少数组参数"))
    }
    
    if array, ok := arrayVal.(*data.ArrayValue); ok {
        // 处理可变参数
        for i := 1; ; i++ {
            val, exists := ctx.GetIndexValue(i)
            if !exists {
                break
            }
            array.Items = append(array.Items, val)
        }
        return data.NewIntValue(len(array.Items)), nil
    }
    
    return data.NewIntValue(0), nil
}
```

### 4. 可变参数函数 - printf（多个值参数）
```go
package stringfuncs

import (
    "fmt"
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type PrintfFunction struct{}

func (f *PrintfFunction) GetName() string { return "printf" }
func (f *PrintfFunction) GetIsStatic() bool { return true }

func (f *PrintfFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "format", 0, nil, data.NewBaseType("string")),
        node.NewParameters(nil, "values", 1, nil, nil), // 可变参数
    }
}

func (f *PrintfFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
        node.NewVariable(nil, "values", 1, nil),
    }
}

func (f *PrintfFunction) GetReturnType() data.Types {
    return data.NewBaseType("int")
}

func (f *PrintfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    formatVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少格式字符串参数"))
    }
    
    if formatStr, ok := formatVal.(data.AsString); ok {
        // 收集可变参数
        var args []interface{}
        for i := 1; ; i++ {
            arg, exists := ctx.GetIndexValue(i)
            if !exists {
                break
            }
            
            // 转换为Go值
            if asString, ok := arg.(data.AsString); ok {
                args = append(args, asString.AsString())
            } else if asInt, ok := arg.(data.AsInt); ok {
                args = append(args, asInt.AsInt())
            } else if asFloat, ok := arg.(data.AsFloat); ok {
                args = append(args, asFloat.AsFloat64())
            } else {
                args = append(args, arg)
            }
        }
        
        // 执行格式化并输出
        result := fmt.Printf(formatStr.AsString(), args...)
        return data.NewIntValue(result), nil
    }
    
    return data.NewIntValue(0), nil
}
        if numVal, ok := val.(data.AsNumber); ok {
            if minNum, ok := minValue.(data.AsNumber); ok {
                if numVal.AsFloat64() < minNum.AsFloat64() {
                    minValue = val
                }
            }
        }
    }
    
    if first {
        return data.NewNullValue(), nil
    }
    
    return minValue, nil
}
```

### 4. 文件处理函数 - file_get_contents（带可选参数）
```go
package filefuncs

import (
    "io/ioutil"
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type FileGetContentsFunction struct{}

func (f *FileGetContentsFunction) GetName() string { return "file_get_contents" }
func (f *FileGetContentsFunction) GetIsStatic() bool { return true }

func (f *FileGetContentsFunction) GetParams() []data.GetValue {
    return []data.GetValue{
        node.NewParameter(nil, "filename", 0, nil, data.NewBaseType("string")),
        node.NewParameter(nil, "use_include_path", 1, data.NewBoolValue(false), data.NewBaseType("bool")),
        node.NewParameter(nil, "context", 2, data.NewNullValue(), data.NewBaseType("resource")),
        node.NewParameter(nil, "offset", 3, data.NewIntValue(0), data.NewBaseType("int")),
        node.NewParameter(nil, "maxlen", 4, data.NewNullValue(), data.NewBaseType("int")),
    }
}

func (f *FileGetContentsFunction) GetVariables() []data.Variable {
    return []data.Variable{
        node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
        node.NewVariable(nil, "use_include_path", 1, data.NewBaseType("bool")),
        node.NewVariable(nil, "context", 2, data.NewBaseType("resource")),
        node.NewVariable(nil, "offset", 3, data.NewBaseType("int")),
        node.NewVariable(nil, "maxlen", 4, data.NewBaseType("int")),
    }
}

func (f *FileGetContentsFunction) GetReturnType() data.Types {
    return data.NewBaseType("string")
}

func (f *FileGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    filenameVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少文件名参数"))
    }
    
    if filenameStr, ok := filenameVal.(data.AsString); ok {
        filename := filenameStr.AsString()
        
        // 读取文件内容
        content, err := ioutil.ReadFile(filename)
        if err != nil {
            return data.NewBoolValue(false), data.NewErrorThrow(nil, err)
        }
        
        return data.NewStringValue(string(content)), nil
    }
    
    return data.NewBoolValue(false), nil
}
```

### 5. AST参数函数 - isset（原始AST节点）
```go
package core

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

type IssetFunction struct{}

func (f *IssetFunction) GetName() string { return "isset" }
func (f *IssetFunction) GetIsStatic() bool { return true }

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

func (f *IssetFunction) GetReturnType() data.Types {
    return data.NewBaseType("bool")
}

func (f *IssetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    varVal, exists := ctx.GetIndexValue(0)
    if !exists {
        return data.NewBoolValue(false), nil
    }
    
    // 检查变量是否已设置且不为null
    if varVal == nil || varVal.IsNull() {
        return data.NewBoolValue(false), nil
    }
    
    return data.NewBoolValue(true), nil
}
```

## 参数处理最佳实践

### 1. 参数验证模式
```go
func (f *ValidationExampleFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 1. 获取参数并处理错误
    param1, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少必需参数"))
    }
    
    param2, exists := ctx.GetIndexValue(1)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少第二个参数"))
    }
    
    // 2. 类型验证
    strParam, ok := param1.(data.AsString)
    if !ok {
        return nil, data.NewErrorThrow(nil, errors.New("第一个参数必须是字符串"))
    }
    
    intParam, ok := param2.(data.AsInt)
    if !ok {
        return nil, data.NewErrorThrow(nil, errors.New("第二个参数必须是整数"))
    }
    
    // 3. 值验证
    if intParam.AsInt() < 0 {
        return nil, data.NewErrorThrow(nil, errors.New("参数值不能为负数"))
    }
    
    // 4. 执行逻辑
    result := process(strParam.AsString(), intParam.AsInt())
    return data.NewStringValue(result), nil
}
```

### 2. 可选参数处理
```go
func (f *OptionalParamFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 检查必需参数
    requiredParam, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少必需参数"))
    }
    
    // 检查可选参数是否存在
    optionalParam, hasOptional := ctx.GetIndexValue(1)
    if !hasOptional {
        // 使用默认值
        optionalParam = data.NewStringValue("default_value")
    }
    
    // 处理逻辑...
}
```

### 3. 可变参数处理
```go
func (f *VariadicFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 检查必需参数
    requiredParam, exists := ctx.GetIndexValue(0)
    if !exists {
        return nil, data.NewErrorThrow(nil, errors.New("缺少必需参数"))
    }
    
    // 收集所有可变参数
    var variadicArgs []data.GetValue
    for i := 1; ; i++ {
        arg, exists := ctx.GetIndexValue(i)
        if !exists {
            break
        }
        variadicArgs = append(variadicArgs, arg)
    }
    
    // 处理逻辑...
}
```

这些示例展示了不同类型的参数处理方式，可以作为实现新函数时的参考模板。