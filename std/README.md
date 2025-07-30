 
---

# Origami 脚本域类与函数封装规范

## 一、函数的封装

### 1. 实现 data.FuncStmt 接口

每个函数都需要实现 `data.FuncStmt` 接口，通常包括如下方法：

- `Call(ctx data.Context) (data.GetValue, data.Control)`：函数的实际执行逻辑。
- `GetName() string`：函数名（在脚本中调用的名字）。
- `GetParams() []data.GetValue`：参数定义（用于参数校验和自动补全）。
- `GetVariables() []data.Variable`：变量定义（类型提示等）。

### 2. 工厂函数

每个函数建议提供一个工厂方法，如：

```go
func NewIsDirFunction() data.FuncStmt {
    return &IsDirFunction{}
}
```

### 3. 注册到 VM

在 `std/load.go` 里通过 `vm.AddFunc(NewXXXFunction())` 注册到虚拟机。

### 4. 示例

```go
type IsDirFunction struct{}

func (f *IsDirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
    // 业务逻辑
}
func (f *IsDirFunction) GetName() string { return "is_dir" }
func (f *IsDirFunction) GetParams() []data.GetValue { ... }
func (f *IsDirFunction) GetVariables() []data.Variable { ... }
```

---

## 二、类的封装

### 1. 实现 data.ClassStmt 接口

每个类都需要实现 `data.ClassStmt` 接口，常见方法有：

- `GetName() string`：类名（脚本中调用的名字）。
- `GetMethod(name string) (data.Method, bool)`：根据方法名获取方法。
- `GetMethods() []data.Method`：获取所有方法。
- `GetProperty(name string) (data.Property, bool)`：获取属性。
- `GetProperties() map[string]data.Property`：获取所有属性。
- `GetConstruct() data.Method`：构造方法（可选）。

### 2. 工厂函数

每个类建议提供一个工厂方法，如：

```go
func NewLogClass() data.ClassStmt {
    source := NewLog()
    return &LogClass{
        debug:  &LogDebugMethod{source},
        // ...
    }
}
```

### 3. 方法包装

每个方法都需要实现 `data.Method` 接口，通常以 `XXXMethod` 结尾，并持有原始结构体指针。

### 4. 注册到 VM

在 `std/load.go` 里通过 `vm.AddClass(NewXXXClass())` 注册到虚拟机。

### 5. 示例

```go
type LogClass struct {
    node.Node
    debug  data.Method
    // ...
}
func (s *LogClass) GetName() string { return "Log" }
func (s *LogClass) GetMethod(name string) (data.Method, bool) { ... }
func (s *LogClass) GetMethods() []data.Method { ... }
```

---

## 三、自动生成工具

### 1. 使用 tools 包自动生成

`std/tools` 目录下提供了自动生成 Go 类包装器的工具。核心流程如下：

- 通过 `tools.Generate(instance, config)` 自动分析结构体并生成包装代码。
- 支持自定义包名、类名、结构体名、输出目录等。

#### 示例

```go
err := tools.Generate(newDateTime(), &tools.GeneratorConfig{
    PackageName: "system",
    ClassName:   "System\\DateTime",
    StructName:  "DateTime",
    OutputDir:   ".", // 生成到当前目录
})
```

### 2. 支持的类型

- 基本类型参数 (int, string, bool, float)
- 指针类型参数
- 接口类型参数
- 切片和数组类型
- 自定义结构体类型

### 3. 生成内容

- 类定义文件：`{struct_name}_class.go`
- 每个方法一个包装器文件

--- 