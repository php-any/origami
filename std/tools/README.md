# 通用结构体包装器生成工具

## 功能说明

这个工具用于根据任意 Go 结构体自动生成符合 `data.Method` 或  `data.ClassStmt` 接口的包装器代码，使得 Go 结构体的方法可以在 PHP 风格的运行时中被调用。

## 核心功能

### 1. 结构体分析
- 通过反射分析任意 Go 结构体
- 提取所有方法的签名信息
- 分析参数类型和返回值类型
- 识别方法的访问修饰符

### 2. 类定义生成
- 生成类结构体，包含所有方法的字段
- 实现 `data.ClassStmt` 接口
- 提供方法查找和实例化功能
- 支持自定义类名和命名空间

### 3. 方法包装器生成
- 为每个方法生成对应的包装器结构体
- 实现 `data.Method` 接口
- 处理参数验证和类型转换
- 调用原始结构体的方法
- 生成错误处理逻辑

### 4. 参数处理
- 根据原始方法的参数类型生成对应的参数验证
- 生成参数类型转换逻辑
- 支持基本类型和复杂类型的转换
- 生成参数缺失时的错误处理

## 工具组件

- `types.go`: 定义数据结构类型
- `analyzer.go`: 结构体分析方法
- `generator.go`: 代码生成器
- `templates.go`: 代码模板定义
- `config.go`: 配置管理
- `main.go`: 主程序入口

## 使用示例

```go
// 分析任意结构体并生成代码
generator := NewWrapperGenerator()
config := &GeneratorConfig{
    PackageName: "http",
    ClassName: "Net\\Http\\Server",
    StructName: "Server",
    OutputDir: "std/net/http",
}
files, err := generator.Generate(&Server{}, config)
```

## 支持的结构体类型

- 基本类型参数 (int, string, bool, float)
- 指针类型参数 (*data.StringValue, *data.FuncValue)
- 接口类型参数 (data.Context)
- 切片和数组类型
- 自定义结构体类型

## 生成的文件命名规则

- 类定义文件: `{struct_name}_class.go`
- 方法包装器文件: `{struct_name}_{method_name}_method.go` 

# 复制代码用于生成代码
```azure
	err := tools.Generate(newDateTime(), &tools.GeneratorConfig{
		PackageName: "system",
		ClassName:   "System\\\\DateTime",
		StructName:  "DateTime",
		OutputDir:   ".", // 生成到当前目录
	})
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}
```