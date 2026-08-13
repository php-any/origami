package tools

import (
	"reflect"
)

// MethodInfo 方法信息
type MethodInfo struct {
	Name       string      // 方法名
	Params     []ParamInfo // 参数列表
	ReturnType string      // 返回类型
	Modifier   string      // 访问修饰符 (public/private/protected)
	IsStatic   bool        // 是否为静态方法
}

// ParamInfo 参数信息
type ParamInfo struct {
	Name     string // 参数名
	Type     string // 参数类型
	Index    int    // 参数索引
	Required bool   // 是否必需
}

// PropertyInfo 属性信息
type PropertyInfo struct {
	Name  string // 属性名
	Type  string // 属性类型 (String/Int/Bool/Float)
	Value string // 属性值
}

// StructInfo 结构体信息
type StructInfo struct {
	Name       string         // 结构体名
	Methods    []MethodInfo   // 方法列表
	Properties []PropertyInfo // 属性列表
}

// GeneratorConfig 生成器配置
type GeneratorConfig struct {
	PackageName string         // 包名
	ClassName   string         // 类名 (如 "Net\\Http\\Server")
	StructName  string         // 结构体名 (如 "Server")
	OutputDir   string         // 输出目录
	Namespace   string         // 命名空间
	Properties  []PropertyInfo // 属性列表
}

// GeneratedFile 生成的文件信息
type GeneratedFile struct {
	FileName string // 文件名
	Content  string // 文件内容
	FilePath string // 文件路径
	FileType string // 文件类型 (class/method)
}

// TypeConverter 类型转换器
type TypeConverter struct {
	GoType    string // Go 类型
	DataType  string // data 包中的类型
	Converter string // 转换函数
}

// MethodWrapper 方法包装器信息
type MethodWrapper struct {
	StructName   string      // 结构体名
	MethodName   string      // 方法名
	WrapperName  string      // 包装器结构体名
	Params       []ParamInfo // 参数列表
	ReturnType   string      // 返回类型
	SourceCall   string      // 原始方法调用
	ParamChecks  []string    // 参数检查代码
	TypeConverts []string    // 类型转换代码
}

// ClassWrapper 类包装器信息
type ClassWrapper struct {
	StructName    string       // 结构体名
	ClassName     string       // 类名
	PackageName   string       // 包名
	Methods       []MethodInfo // 方法列表
	MethodFields  []string     // 方法字段列表
	MethodCases   []string     // switch case 列表
	MethodReturns []string     // 返回方法列表
}

// Analyzer 结构体分析器接口
type Analyzer interface {
	AnalyzeStruct(instance interface{}) (*StructInfo, error)
	ExtractMethods(value reflect.Value) []MethodInfo
	ExtractParams(method reflect.Method) []ParamInfo
	GetTypeConverter(goType string) *TypeConverter
}

// Generator 代码生成器接口
type Generator interface {
	Generate(instance interface{}, config *GeneratorConfig) ([]GeneratedFile, error)
	GenerateClass(structInfo *StructInfo, config *GeneratorConfig) (*GeneratedFile, error)
	GenerateMethods(structInfo *StructInfo, config *GeneratorConfig) ([]GeneratedFile, error)
	GenerateMethodWrapper(method MethodInfo, config *GeneratorConfig) (*GeneratedFile, error)
}
