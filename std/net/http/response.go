package http

import (
	"encoding/json"
	"github.com/php-any/origami/data"
	"net/http"
)

func NewResponse(w http.ResponseWriter, r *http.Request) *Response {
	return &Response{w, r}
}

type Response struct {
	w http.ResponseWriter
	r *http.Request
}

func (req *Response) Write(text string) (data.GetValue, data.Control) {
	n, err := req.w.Write([]byte(text))
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewIntValue(n), nil
}

func (req *Response) Json(v any) (data.GetValue, data.Control) {
	// 将 data.Value 转换为 Go 原生类型
	nativeValue := req.convertToNative(v)

	// 序列化为 JSON
	jsonBytes, err := json.Marshal(nativeValue)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	// 设置 Content-Type 头
	req.w.Header().Set("Content-Type", "application/json")

	// 写入响应
	n, err := req.w.Write(jsonBytes)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	return data.NewIntValue(n), nil
}

// convertToNative 将 data.Value 类型转换为 Go 原生类型，用于 JSON 序列化
func (req *Response) convertToNative(v any) any {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case *data.StringValue:
		return val.AsString()

	case *data.IntValue:
		intVal, _ := val.AsInt()
		return intVal

	case *data.FloatValue:
		floatVal, _ := val.AsFloat()
		return floatVal

	case *data.BoolValue:
		boolVal, _ := val.AsBool()
		return boolVal

	case *data.NullValue:
		return nil

	case *data.ArrayValue:
		result := make([]any, len(val.Value))
		for i, item := range val.Value {
			result[i] = req.convertToNative(item)
		}
		return result

	case *data.ObjectValue:
		result := make(map[string]any)
		properties := val.GetProperties()
		for key, value := range properties {
			result[key] = req.convertToNative(value)
		}
		return result

	case *data.ClassValue:
		result := make(map[string]any)
		properties := val.GetProperties()
		for key, value := range properties {
			result[key] = req.convertToNative(value)
		}
		return result

	default:
		// 对于其他类型，尝试转换为字符串
		if stringer, ok := val.(interface{ AsString() string }); ok {
			return stringer.AsString()
		}
		// 如果无法转换，返回 nil
		return nil
	}
}
