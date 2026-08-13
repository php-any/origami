package tools

import (
	"fmt"
	"reflect"
	"strings"
)

// StructAnalyzer 结构体分析器
type StructAnalyzer struct {
	typeConverters map[string]*TypeConverter
}

// NewStructAnalyzer 创建新的结构体分析器
func NewStructAnalyzer() *StructAnalyzer {
	analyzer := &StructAnalyzer{
		typeConverters: make(map[string]*TypeConverter),
	}
	analyzer.initTypeConverters()
	return analyzer
}

// AnalyzeStruct 分析结构体
func (a *StructAnalyzer) AnalyzeStruct(instance interface{}) (*StructInfo, error) {
	typ := reflect.TypeOf(instance)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("实例必须是结构体类型")
	}

	structInfo := &StructInfo{
		Name:       typ.Name(),
		Methods:    a.ExtractMethods(reflect.TypeOf(instance)),
		Properties: a.ExtractProperties(typ, instance),
	}

	return structInfo, nil
}

// ExtractMethods 提取结构体的所有方法
func (a *StructAnalyzer) ExtractMethods(typ reflect.Type) []MethodInfo {
	var methods []MethodInfo

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)

		// 跳过私有方法（Go中首字母小写的方法）
		if !a.isPublicMethod(method.Name) {
			continue
		}

		methodInfo := MethodInfo{
			Name:       method.Name,
			Params:     a.ExtractParams(method),
			ReturnType: a.extractReturnType(method.Type),
			Modifier:   "public",
			IsStatic:   false,
		}

		methods = append(methods, methodInfo)
	}

	return methods
}

// ExtractProperties 提取结构体的所有字段作为属性
func (a *StructAnalyzer) ExtractProperties(typ reflect.Type, instance interface{}) []PropertyInfo {
	var properties []PropertyInfo
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// 跳过私有字段（Go中首字母小写的字段）
		if !a.isPublicField(field.Name) {
			continue
		}

		// 获取字段值
		fieldValue := val.Field(i)
		propertyValue := a.getFieldValue(fieldValue)

		propertyInfo := PropertyInfo{
			Name:  field.Name,
			Type:  a.getPropertyType(field.Type),
			Value: propertyValue,
		}

		properties = append(properties, propertyInfo)
	}

	return properties
}

// ExtractParams 提取方法的参数信息
func (a *StructAnalyzer) ExtractParams(method reflect.Method) []ParamInfo {
	var params []ParamInfo
	funcType := method.Type

	// 跳过第一个参数（接收者）
	for i := 1; i < funcType.NumIn(); i++ {
		paramType := funcType.In(i)

		paramInfo := ParamInfo{
			Name:     a.generateParamName(i - 1),
			Type:     a.getTypeString(paramType),
			Index:    i - 1,
			Required: true,
		}
		params = append(params, paramInfo)
	}

	return params
}

// GetTypeConverter 获取类型转换器
func (a *StructAnalyzer) GetTypeConverter(goType string) *TypeConverter {
	return a.typeConverters[goType]
}

// isPublicMethod 判断是否为公有方法
func (a *StructAnalyzer) isPublicMethod(name string) bool {
	return len(name) > 0 && strings.ToUpper(name[:1]) == name[:1]
}

// isPublicField 判断是否为公有字段
func (a *StructAnalyzer) isPublicField(name string) bool {
	return len(name) > 0 && strings.ToUpper(name[:1]) == name[:1]
}

// getPropertyType 获取属性类型
func (a *StructAnalyzer) getPropertyType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "String"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "Int"
	case reflect.Bool:
		return "Bool"
	case reflect.Float32, reflect.Float64:
		return "Float"
	default:
		return "String"
	}
}

// getFieldValue 获取字段值的字符串表示
func (a *StructAnalyzer) getFieldValue(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		// 处理字符串中的特殊字符
		str := val.String()
		str = strings.ReplaceAll(str, "\n", "\\n")
		str = strings.ReplaceAll(str, "\r", "\\r")
		str = strings.ReplaceAll(str, "\t", "\\t")
		str = strings.ReplaceAll(str, "\"", "\\\"")
		return fmt.Sprintf(`"%s"`, str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Bool:
		return fmt.Sprintf("%t", val.Bool())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", val.Float())
	default:
		return `""`
	}
}

// generateParamName 生成参数名
func (a *StructAnalyzer) generateParamName(index int) string {
	return fmt.Sprintf("param%d", index)
}

// extractReturnType 提取返回类型
func (a *StructAnalyzer) extractReturnType(funcType reflect.Type) string {
	numOut := funcType.NumOut()
	if numOut == 0 {
		return ""
	}
	if numOut == 1 {
		return a.getTypeString(funcType.Out(0))
	}

	var returnTypes []string
	for i := 0; i < numOut; i++ {
		returnTypes = append(returnTypes, a.getTypeString(funcType.Out(i)))
	}
	return "(" + strings.Join(returnTypes, ", ") + ")"
}

// getTypeString 获取类型字符串
func (a *StructAnalyzer) getTypeString(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return "*" + a.getTypeString(t.Elem())
	case reflect.Slice:
		return "[]" + a.getTypeString(t.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), a.getTypeString(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", a.getTypeString(t.Key()), a.getTypeString(t.Elem()))
	case reflect.Interface:
		if t.Name() == "" {
			return "interface{}"
		}
		return t.Name()
	case reflect.Struct:
		if t.PkgPath() != "" {
			return t.PkgPath() + "." + t.Name()
		}
		return t.Name()
	default:
		return t.Name()
	}
}

// initTypeConverters 初始化类型转换器
func (a *StructAnalyzer) initTypeConverters() {
	// 基本类型转换器
	a.typeConverters["string"] = &TypeConverter{
		GoType:    "string",
		DataType:  "*data.StringValue",
		Converter: ".AsString()",
	}

	a.typeConverters["int"] = &TypeConverter{
		GoType:    "int",
		DataType:  "*data.IntValue",
		Converter: ".AsInt()",
	}

	a.typeConverters["bool"] = &TypeConverter{
		GoType:    "bool",
		DataType:  "*data.BoolValue",
		Converter: ".AsBool()",
	}

	a.typeConverters["float64"] = &TypeConverter{
		GoType:    "float64",
		DataType:  "*data.FloatValue",
		Converter: ".AsFloat()",
	}

	// 接口类型转换器
	a.typeConverters["data.Context"] = &TypeConverter{
		GoType:   "data.Context",
		DataType: "data.Context",
	}

	// 指针类型转换器
	a.typeConverters["*data.StringValue"] = &TypeConverter{
		GoType:   "*data.StringValue",
		DataType: "*data.StringValue",
	}

	a.typeConverters["*data.FuncValue"] = &TypeConverter{
		GoType:   "*data.FuncValue",
		DataType: "*data.FuncValue",
	}

	a.typeConverters["*data.IntValue"] = &TypeConverter{
		GoType:   "*data.IntValue",
		DataType: "*data.IntValue",
	}

	a.typeConverters["*data.BoolValue"] = &TypeConverter{
		GoType:   "*data.BoolValue",
		DataType: "*data.BoolValue",
	}

	a.typeConverters["*data.FloatValue"] = &TypeConverter{
		GoType:   "*data.FloatValue",
		DataType: "*data.FloatValue",
	}
}
