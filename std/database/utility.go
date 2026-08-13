package database

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/serializer/json"
)

var serializer = json.NewJsonSerializer()

// ConvertValueToGoType 将 data.Value 转换为 Go 原生类型
// 这是一个独立的工具函数，用于替换各个数据库方法中重复的 convertValueToGoType 方法
func ConvertValueToGoType(val data.Value) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case data.ValueSerializer:
		b, err := v.ToGoValue(serializer)
		if err == nil {
			return b
		}
		if vv, ok := v.(data.Value); ok {
			return vv.AsString()
		}
		return nil
	default:
		// 对于其他类型，尝试转换为字符串
		return v.AsString()
	}
}
