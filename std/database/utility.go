package database

import (
	"github.com/php-any/origami/data"
)

// ConvertValueToGoType 将 data.Value 转换为 Go 原生类型
// 这是一个独立的工具函数，用于替换各个数据库方法中重复的 convertValueToGoType 方法
func ConvertValueToGoType(val data.Value) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case *data.IntValue:
		return v.Value
	case *data.StringValue:
		return v.Value
	case *data.BoolValue:
		return v.Value
	case *data.FloatValue:
		return v.Value
	case *data.NullValue:
		return nil
	case *data.ArrayValue:
		// 对于数组，转换为 []interface{}
		result := make([]interface{}, len(v.Value))
		for i, item := range v.Value {
			result[i] = ConvertValueToGoType(item)
		}
		return result
	case *data.ObjectValue:
		// 对于对象，转换为 map[string]interface{}
		result := make(map[string]interface{})
		for k, item := range v.GetProperties() {
			result[k] = ConvertValueToGoType(item)
		}
		return result
	default:
		// 对于其他类型，尝试转换为字符串
		return v.AsString()
	}
}
