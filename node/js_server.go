package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// JsServerExpression 表示 $.SERVER() 表达式
type JsServerExpression struct {
	from    data.From
	varName string
}

// NewJsServerExpression 创建一个新的 JS_SERVER 表达式节点
func NewJsServerExpression(from data.From, varName string) *JsServerExpression {
	return &JsServerExpression{
		from:    from,
		varName: varName,
	}
}

// GetFrom 获取位置信息
func (n *JsServerExpression) GetFrom() data.From {
	return n.from
}

// GetValue 获取表达式的值
// 在运行时，这将被转换为 JavaScript 值格式的字符串
func (n *JsServerExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建一个变量引用
	variable := NewVariable(n.from, n.varName, 0, nil)

	// 获取变量的值
	varValue, ctl := ctx.GetVariableValue(variable)
	if ctl != nil {
		return nil, ctl
	}

	// 将变量值转换为 JavaScript 格式的字符串
	jsValue := convertToJavaScriptValue(varValue)
	return data.NewStringValue(jsValue), nil
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

// GetVarName 获取变量名
func (n *JsServerExpression) GetVarName() string {
	return n.varName
}
