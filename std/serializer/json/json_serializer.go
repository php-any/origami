package json

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// JsonSerializer 实现 data.Serializer 接口
type JsonSerializer struct{}

func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

// Int
func (j *JsonSerializer) MarshalInt(v *data.IntValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (j *JsonSerializer) UnmarshalInt(data []byte, v *data.IntValue) error {
	return json.Unmarshal(data, &v.Value)
}

// String
func (j *JsonSerializer) MarshalString(v *data.StringValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (j *JsonSerializer) UnmarshalString(data []byte, v *data.StringValue) error {
	return json.Unmarshal(data, &v.Value)
}

// Null
func (j *JsonSerializer) MarshalNull(v *data.NullValue) ([]byte, error) {
	return json.Marshal(nil)
}

func (j *JsonSerializer) UnmarshalNull(data []byte, v *data.NullValue) error {
	// NullValue 不需要反序列化任何内容
	return nil
}

// Array
func (j *JsonSerializer) MarshalArray(v *data.ArrayValue) ([]byte, error) {
	// 递归序列化数组中的每个元素
	items := make([]json.RawMessage, 0, len(v.Value))
	for _, elem := range v.Value {
		if vs, ok := elem.(data.ValueSerializer); ok {
			// 为每个元素创建新的序列化器
			b, err := vs.Marshal(j)
			if err != nil {
				return nil, err
			}
			items = append(items, b)
		} else {
			// 如果不支持序列化，转换为字符串
			b, _ := json.Marshal(elem.AsString())
			items = append(items, b)
		}
	}
	return json.Marshal(items)
}

func (j *JsonSerializer) UnmarshalArray(msg []byte, v *data.ArrayValue) error {
	var items []json.RawMessage
	if err := json.Unmarshal(msg, &items); err != nil {
		return err
	}

	values := make([]data.Value, 0, len(items))
	for idx, item := range items {
		// 优先使用已有元素的类型信息
		if idx < len(v.Value) && v.Value[idx] != nil {
			val, err := j.unmarshalWithExpected(item, v.Value[idx])
			if err != nil {
				return err
			}
			values = append(values, val)
			continue
		}

		// 回退到类型猜测
		val, err := j.unmarshalValue(item)
		if err != nil {
			return err
		}
		values = append(values, val)
	}
	v.Value = values
	return nil
}

// Object
func (j *JsonSerializer) MarshalObject(v *data.ObjectValue) ([]byte, error) {
	props := v.GetProperties()

	// 递归序列化对象属性
	encoded := make(map[string]json.RawMessage, len(props))
	for k, val := range props {
		if vs, ok := val.(data.ValueSerializer); ok {
			b, err := vs.Marshal(j)
			if err != nil {
				return nil, err
			}
			encoded[k] = b
		} else {
			// 如果不支持序列化，转换为字符串
			b, _ := json.Marshal(val.AsString())
			encoded[k] = b
		}
	}
	return json.Marshal(encoded)
}

func (j *JsonSerializer) UnmarshalObject(data []byte, v *data.ObjectValue) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, raw := range m {
		// 优先根据已有属性值的具体类型进行反序列化
		if existing, ok := v.GetProperty(k); ok && existing != nil {
			val, err := j.unmarshalWithExpected(raw, existing)
			if err != nil {
				return err
			}
			v.SetProperty(k, val)
			continue
		}

		// 回退到类型猜测
		val, err := j.unmarshalValue(raw)
		if err != nil {
			return err
		}
		v.SetProperty(k, val)
	}
	return nil
}

// Bool
func (j *JsonSerializer) MarshalBool(v *data.BoolValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (j *JsonSerializer) UnmarshalBool(data []byte, v *data.BoolValue) error {
	return json.Unmarshal(data, &v.Value)
}

// Float
func (j *JsonSerializer) MarshalFloat(v *data.FloatValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (j *JsonSerializer) UnmarshalFloat(data []byte, v *data.FloatValue) error {
	return json.Unmarshal(data, &v.Value)
}

// Any
func (j *JsonSerializer) MarshalAny(v *data.AnyValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (j *JsonSerializer) UnmarshalAny(data []byte, v *data.AnyValue) error {
	var tmp any
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Value = tmp
	return nil
}

// Class
func (j *JsonSerializer) MarshalClass(v *data.ClassValue) ([]byte, error) {
	props := v.GetProperties()

	// 序列化类实例的属性
	encoded := make(map[string]json.RawMessage, len(props))
	for k, val := range props {
		if vs, ok := val.(data.ValueSerializer); ok {
			b, err := vs.Marshal(j)
			if err != nil {
				return nil, err
			}
			encoded[k] = b
		} else {
			b, _ := json.Marshal(val.AsString())
			encoded[k] = b
		}
	}

	// 与对象保持一致：仅输出属性映射
	return json.Marshal(encoded)
}

func (j *JsonSerializer) UnmarshalClass(msg []byte, v *data.ClassValue) error {
	// 直接从对象结构恢复属性
	var props map[string]json.RawMessage
	if err := json.Unmarshal(msg, &props); err != nil {
		return err
	}

	for k, raw := range props {
		// 1) 首先尝试根据类定义的属性类型创建相应类型的值
		if prop, ok := v.GetProperty(k); ok && prop != nil {
			// 获取属性的类型信息
			if typedProp, ok := prop.(*node.ClassProperty); ok {
				propType := typedProp.GetType()
				if propType != nil {
					// 根据属性类型创建相应类型的默认值
					expectedValue := j.createDefaultValueForType(propType)
					if expectedValue != nil {
						// 使用属性的类型作为预期类型
						val, err := j.unmarshalWithExpected(raw, expectedValue)
						if err != nil {
							return err
						}
						v.SetProperty(k, val)
						continue
					}
				}
			}
		}

		// 2) 使用实例当前值类型作为备选
		if instVal, ok := v.ObjectValue.GetProperty(k); ok && instVal != nil {
			val, err := j.unmarshalWithExpected(raw, instVal)
			if err != nil {
				return err
			}
			v.SetProperty(k, val)
			continue
		}

		// 3) 回退到类型猜测
		val, err := j.unmarshalValue(raw)
		if err != nil {
			return err
		}
		v.SetProperty(k, val)
	}
	return nil
}

// createDefaultValueForType 根据类型信息创建相应类型的默认值
func (j *JsonSerializer) createDefaultValueForType(ty data.Types) data.Value {
	// 这里可以根据类型信息创建相应类型的默认值
	// 由于目前类型系统还没有完全实现，我们暂时返回nil
	// 在实际应用中，可以根据ty的类型创建相应的默认值
	switch ty.(type) {
	case data.String:
		return data.NewStringValue("")
	case data.Int:
		return data.NewIntValue(0)
	case data.Bool:
		return data.NewBoolValue(false)
	case data.Float:
		return data.NewFloatValue(0.0)
	case data.Arrays:
		return data.NewArrayValue([]data.Value{})
	case data.Object:
		return data.NewObjectValue()
	default:
		// 对于其他类型，暂时返回nil，让后续逻辑处理
		return nil
	}
}

// 辅助方法：根据 JSON 数据推断类型并反序列化
func (j *JsonSerializer) unmarshalValue(raw []byte) (data.Value, error) {
	// 基于首字符进行类型判定，避免多次“猜测”解析
	s := bytes.TrimSpace(raw)
	if len(s) == 0 {
		return data.NewNullValue(), nil
	}
	ch := s[0]
	switch ch {
	case '"': // 字符串
		var str string
		if err := json.Unmarshal(s, &str); err != nil {
			return nil, err
		}
		return data.NewStringValue(str), nil
	case 'n': // null 字面量
		if bytes.Equal(s, []byte("null")) {
			return data.NewNullValue(), nil
		}
	case 't', 'f': // 布尔字面量 true/false
		if bytes.Equal(s, []byte("true")) {
			return data.NewBoolValue(true), nil
		}
		if bytes.Equal(s, []byte("false")) {
			return data.NewBoolValue(false), nil
		}
	case '[': // 数组
		av := &data.ArrayValue{}
		if err := j.UnmarshalArray(s, av); err != nil {
			return nil, err
		}
		return av, nil
	case '{': // 对象
		ov := data.NewObjectValue()
		if err := j.UnmarshalObject(s, ov); err != nil {
			return nil, err
		}
		return ov, nil
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // 数字
		// 按是否包含小数点/指数来区分 int 与 float
		isFloat := bytes.Contains(s, []byte{'.'}) || bytes.ContainsAny(s, "eE")
		if isFloat {
			var f float64
			if err := json.Unmarshal(s, &f); err != nil {
				return nil, err
			}
			return data.NewFloatValue(f), nil
		}
		var i int
		if err := json.Unmarshal(s, &i); err != nil {
			return nil, err
		}
		return data.NewIntValue(i), nil
	}
	// 不认识的首字符，退回 any
	var anyVal any
	if err := json.Unmarshal(s, &anyVal); err != nil {
		return nil, fmt.Errorf("无法解析 JSON 数据: %v", err)
	}
	return data.NewAnyValue(anyVal), nil
}

// unmarshalWithExpected 基于期望的 Value 类型进行精准反序列化
func (j *JsonSerializer) unmarshalWithExpected(raw []byte, expected data.Value) (data.Value, error) {
	switch ev := expected.(type) {
	case data.ValueSerializer:
		err := ev.Unmarshal(raw, j)
		return expected, err
	default:
		// 未知具体类型，回退到猜测
		return j.unmarshalValue(raw)
	}
}
