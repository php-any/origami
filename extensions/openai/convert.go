package openai

import "github.com/php-any/origami/data"

// goValueToDataValue 将 Go 原生类型递归转换为 data.Value
func goValueToDataValue(v any) data.Value {
	switch val := v.(type) {
	case nil:
		return data.NewNullValue()
	case string:
		return data.NewStringValue(val)
	case int:
		return data.NewIntValue(val)
	case int64:
		return data.NewIntValue(int(val))
	case float64:
		return data.NewFloatValue(val)
	case float32:
		return data.NewFloatValue(float64(val))
	case bool:
		return data.NewBoolValue(val)
	case []any:
		items := make([]data.Value, len(val))
		for i, item := range val {
			items[i] = goValueToDataValue(item)
		}
		return data.NewArrayValue(items)
	case map[string]any:
		obj := data.NewObjectValue()
		for k, item := range val {
			obj.SetProperty(k, goValueToDataValue(item))
		}
		return obj
	default:
		return data.NewAnyValue(val)
	}
}
