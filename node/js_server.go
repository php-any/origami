package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// JsServerExpression 表示 $.SERVER() 表达式
// 它接收一个参数（变量），并将其转换为 JavaScript 格式输出
type JsServerExpression struct {
	*Node `pp:"-"`
	Args  []data.GetValue // 函数参数列表
}

// NewJsServerExpression 创建一个新的 JS_SERVER 表达式节点
func NewJsServerExpression(from data.From, args []data.GetValue) *JsServerExpression {
	return &JsServerExpression{
		Node: NewNode(from),
		Args: args,
	}
}

// GetValue 获取表达式的值
// 在运行时，这将被转换为 JavaScript 值格式的字符串
func (n *JsServerExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// $.SERVER 应该只有一个参数（变量）
	if len(n.Args) != 1 {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("$.SERVER() 需要一个参数"))
	}

	// 获取参数的值，参考 CallExpression 的处理方式
	param := n.Args[0]
	tempV, acl := param.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	// 检查返回的是 data.Value，参考 CallExpression 的处理
	var varValue data.Value
	if val, ok := tempV.(data.Value); ok {
		varValue = val
	} else {
		// 如果不是 data.Value，使用 null 值
		varValue = data.NewNullValue()
	}

	// 检查变量值的类型，决定返回原始 JavaScript 值还是字符串值
	switch varValue.(type) {
	case *data.IntValue, *data.FloatValue, *data.BoolValue, *data.NullValue, *data.ObjectValue, *data.ArrayValue:
		// 数字、布尔值、null、对象、数组都返回原始 JavaScript 值（不带引号）
		jsValue := convertToJavaScriptValue(varValue)
		return NewJsRawValue(jsValue), nil
	default:
		// 其他类型（主要是字符串）返回字符串值（带引号）
		jsValue := convertToJavaScriptValue(varValue)
		// 检查是否是对象或数组格式（可能是从其他类型转换来的）
		if isJavaScriptObjectOrArray(jsValue) {
			return NewJsRawValue(jsValue), nil
		}
		return data.NewStringValue(jsValue), nil
	}
}

// isJavaScriptObjectOrArray 检查字符串是否是 JavaScript 对象或数组格式
func isJavaScriptObjectOrArray(s string) bool {
	if len(s) == 0 {
		return false
	}
	// 去除前后空白
	trimmed := s
	for len(trimmed) > 0 && (trimmed[0] == ' ' || trimmed[0] == '\t' || trimmed[0] == '\n' || trimmed[0] == '\r') {
		trimmed = trimmed[1:]
	}
	for len(trimmed) > 0 {
		last := len(trimmed) - 1
		if trimmed[last] == ' ' || trimmed[last] == '\t' || trimmed[last] == '\n' || trimmed[last] == '\r' {
			trimmed = trimmed[:last]
		} else {
			break
		}
	}

	// 检查是否以 { 或 [ 开头
	return len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[')
}

// convertToJavaScriptValue 将 Go 值转换为 JavaScript 格式的字符串
func convertToJavaScriptValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case bool:
		if v {
			return "true"
		}
		return "false"
	case int, int8, int16, int32, int64:
		return formatIntValue(v)
	case uint, uint8, uint16, uint32, uint64:
		return formatUintValue(v)
	case float32, float64:
		return formatFloatValue(v)
	case string:
		return formatStringValue(v)
	case []interface{}:
		return formatArrayValue(v)
	case map[string]interface{}:
		return formatObjectValue(v)
	case *data.IntValue:
		// 处理 IntValue，直接输出数字（不带引号）
		return formatDataIntValue(v)
	case *data.FloatValue:
		// 处理 FloatValue，直接输出数字（不带引号）
		return formatDataFloatValue(v)
	case *data.BoolValue:
		// 处理 BoolValue，直接输出布尔值（不带引号）
		return formatDataBoolValue(v)
	case *data.NullValue:
		// 处理 NullValue，输出 null（不带引号）
		return "null"
	case *data.ObjectValue:
		// 处理 ObjectValue，转换为 JavaScript 对象格式
		return formatDataObjectValue(v)
	case *data.ArrayValue:
		// 处理 ArrayValue，转换为 JavaScript 数组格式
		return formatDataArrayValue(v)
	default:
		// 对于其他类型，尝试转换为字符串
		if strValue, ok := value.(data.AsString); ok {
			return formatStringValue(strValue.AsString())
		}
		return formatStringValue("")
	}
}

// formatIntValue 格式化整数值
func formatIntValue(v interface{}) string {
	if intValue, ok := v.(data.AsInt); ok {
		if intVal, err := intValue.AsInt(); err == nil {
			return fmt.Sprintf("%d", intVal)
		}
	}
	return "0"
}

// formatUintValue 格式化无符号整数值
func formatUintValue(v interface{}) string {
	if intValue, ok := v.(data.AsInt); ok {
		if intVal, err := intValue.AsInt(); err == nil {
			return fmt.Sprintf("%d", intVal)
		}
	}
	return "0"
}

// formatFloatValue 格式化浮点数值
func formatFloatValue(v interface{}) string {
	if floatValue, ok := v.(data.AsFloat); ok {
		if floatVal, err := floatValue.AsFloat(); err == nil {
			return fmt.Sprintf("%f", floatVal)
		}
	}
	return "0.0"
}

// formatStringValue 格式化字符串值
func formatStringValue(v string) string {
	// 转义特殊字符
	result := "\""
	for _, r := range v {
		switch r {
		case '"':
			result += "\\\""
		case '\\':
			result += "\\\\"
		case '\n':
			result += "\\n"
		case '\r':
			result += "\\r"
		case '\t':
			result += "\\t"
		default:
			result += string(r)
		}
	}
	result += "\""
	return result
}

// formatArrayValue 格式化数组值
func formatArrayValue(v []interface{}) string {
	if len(v) == 0 {
		return "[]"
	}

	result := "["
	for i, item := range v {
		if i > 0 {
			result += ", "
		}
		result += convertToJavaScriptValue(item)
	}
	result += "]"
	return result
}

// formatObjectValue 格式化对象值
func formatObjectValue(v map[string]interface{}) string {
	if len(v) == 0 {
		return "{}"
	}

	result := "{"
	first := true
	for key, value := range v {
		if !first {
			result += ", "
		}
		result += formatStringValue(key) + ": " + convertToJavaScriptValue(value)
		first = false
	}
	result += "}"
	return result
}

// formatDataObjectValue 格式化 data.ObjectValue 为 JavaScript 对象格式
func formatDataObjectValue(obj *data.ObjectValue) string {
	properties := obj.GetProperties()
	if len(properties) == 0 {
		return "{}"
	}

	result := "{"
	first := true
	for key, value := range properties {
		if !first {
			result += ", "
		}
		// 获取值的实际值
		if val, ok := value.GetValue(nil); ok == nil {
			if valValue, ok := val.(data.Value); ok {
				result += formatStringValue(key) + ": " + convertToJavaScriptValue(valValue)
			} else {
				result += formatStringValue(key) + ": " + convertToJavaScriptValue(value)
			}
		} else {
			result += formatStringValue(key) + ": " + convertToJavaScriptValue(value)
		}
		first = false
	}
	result += "}"
	return result
}

// formatDataArrayValue 格式化 data.ArrayValue 为 JavaScript 数组格式
func formatDataArrayValue(arr *data.ArrayValue) string {
	values := arr.Value
	if len(values) == 0 {
		return "[]"
	}

	result := "["
	for i, value := range values {
		if i > 0 {
			result += ", "
		}
		// 获取值的实际值
		if val, ok := value.GetValue(nil); ok == nil {
			if valValue, ok := val.(data.Value); ok {
				result += convertToJavaScriptValue(valValue)
			} else {
				result += convertToJavaScriptValue(value)
			}
		} else {
			result += convertToJavaScriptValue(value)
		}
	}
	result += "]"
	return result
}

// formatDataIntValue 格式化 data.IntValue 为 JavaScript 数字格式
func formatDataIntValue(v *data.IntValue) string {
	return fmt.Sprintf("%d", v.Value)
}

// formatDataFloatValue 格式化 data.FloatValue 为 JavaScript 数字格式
func formatDataFloatValue(v *data.FloatValue) string {
	return fmt.Sprintf("%g", v.Value)
}

// formatDataBoolValue 格式化 data.BoolValue 为 JavaScript 布尔值格式
func formatDataBoolValue(v *data.BoolValue) string {
	if v.Value {
		return "true"
	}
	return "false"
}

// JsRawValue 表示一个原始的 JavaScript 值（不会被加上引号）
type JsRawValue struct {
	Value string
}

// NewJsRawValue 创建一个新的 JsRawValue
func NewJsRawValue(value string) *JsRawValue {
	return &JsRawValue{Value: value}
}

// GetValue 获取值
func (j *JsRawValue) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return j, nil
}

// AsString 返回原始的 JavaScript 代码（不带引号）
func (j *JsRawValue) AsString() string {
	return j.Value
}

// AsInt 尝试将值转换为整数
func (j *JsRawValue) AsInt() (int, error) {
	// 对于原始 JavaScript 值，尝试解析
	// 这里简化处理，实际可能需要更复杂的解析
	return 0, fmt.Errorf("无法将 JavaScript 原始值转换为整数")
}

// AsFloat 尝试将值转换为浮点数
func (j *JsRawValue) AsFloat() (float64, error) {
	return 0, fmt.Errorf("无法将 JavaScript 原始值转换为浮点数")
}

// AsBool 尝试将值转换为布尔值
func (j *JsRawValue) AsBool() (bool, error) {
	// 检查是否是 null、false、true
	trimmed := j.Value
	for len(trimmed) > 0 && (trimmed[0] == ' ' || trimmed[0] == '\t' || trimmed[0] == '\n' || trimmed[0] == '\r') {
		trimmed = trimmed[1:]
	}
	if trimmed == "null" || trimmed == "false" {
		return false, nil
	}
	if trimmed == "true" {
		return true, nil
	}
	return true, nil // 非空对象/数组视为 true
}

// GetFrom 获取位置信息（JsRawValue 没有位置信息）
func (j *JsRawValue) GetFrom() data.From {
	return nil
}

// Marshal 序列化
func (j *JsRawValue) Marshal(serializer data.Serializer) ([]byte, error) {
	// 将 JsRawValue 转换为 StringValue 进行序列化
	strValue := &data.StringValue{Value: j.Value}
	return serializer.MarshalString(strValue)
}

// Unmarshal 反序列化
func (j *JsRawValue) Unmarshal(dataBytes []byte, serializer data.Serializer) error {
	// 使用 StringValue 进行反序列化
	strValue := &data.StringValue{}
	if err := serializer.UnmarshalString(dataBytes, strValue); err != nil {
		return err
	}
	j.Value = strValue.Value
	return nil
}

// ToGoValue 转换为 Go 值
func (j *JsRawValue) ToGoValue(_ data.Serializer) (any, error) {
	return j.Value, nil
}

// GetMethod 获取方法（JsRawValue 没有方法）
func (j *JsRawValue) GetMethod(name string) (data.Method, bool) {
	return nil, false
}

// GetProperty 获取属性（JsRawValue 没有属性）
func (j *JsRawValue) GetProperty(name string) (data.Value, bool) {
	return nil, false
}
