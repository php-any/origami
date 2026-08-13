package http

import (
	"encoding/json"
	"io"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

// RequestBodyMethod 获取请求体数据
type RequestBodyMethod struct {
	source *httpsrc.Request
}

func (h *RequestBodyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewStringValue(""), nil
	}

	// 获取原始请求体
	if h.source.Body != nil {
		body, err := io.ReadAll(h.source.Body)
		if err == nil {
			bodyStr := string(body)
			// 尝试解析为 JSON
			if json.Valid(body) {
				var result interface{}
				if err := json.Unmarshal(body, &result); err == nil {
					return convertToDataValue(result), nil
				}
			}
			return data.NewStringValue(bodyStr), nil
		}
	}

	// 如果没有请求体，返回空字符串
	return data.NewStringValue(""), nil
}

// convertToDataValue 将 Go 值转换为 data.Value
func convertToDataValue(v interface{}) data.Value {
	switch val := v.(type) {
	case map[string]interface{}:
		// JSON 对象转换为带命名键的数组
		list := make([]*data.ZVal, 0, len(val))
		for k, v := range val {
			list = append(list, data.NewNamedZVal(k, convertToDataValue(v)))
		}
		return &data.ArrayValue{List: list}
	case []interface{}:
		// JSON 数组转换为普通数组
		values := make([]data.Value, len(val))
		for i, v := range val {
			values[i] = convertToDataValue(v)
		}
		return data.NewArrayValue(values)
	case string:
		return data.NewStringValue(val)
	case float64:
		return data.NewFloatValue(val)
	case bool:
		return data.NewBoolValue(val)
	case nil:
		return data.NewNullValue()
	default:
		return data.NewStringValue("")
	}
}

func (h *RequestBodyMethod) GetName() string               { return "body" }
func (h *RequestBodyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (h *RequestBodyMethod) GetIsStatic() bool             { return false }
func (h *RequestBodyMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (h *RequestBodyMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (h *RequestBodyMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
