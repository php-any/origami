package json

import (
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
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
			serializer := NewJsonSerializer()
			b, err := vs.Marshal(serializer)
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

func (j *JsonSerializer) UnmarshalArray(data []byte, v *data.ArrayValue) error {
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	values := make([]data.Value, 0, len(items))
	for _, item := range items {
		// 尝试推断类型并反序列化
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
			serializer := NewJsonSerializer()
			b, err := vs.Marshal(serializer)
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
			serializer := NewJsonSerializer()
			b, err := vs.Marshal(serializer)
			if err != nil {
				return nil, err
			}
			encoded[k] = b
		} else {
			b, _ := json.Marshal(val.AsString())
			encoded[k] = b
		}
	}

	payload := map[string]any{
		"name":       v.Class.GetName(),
		"properties": encoded,
	}
	return json.Marshal(payload)
}

func (j *JsonSerializer) UnmarshalClass(data []byte, v *data.ClassValue) error {
	var payload struct {
		Name       string                     `json:"name"`
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	// 只恢复属性，不改变类定义
	for k, raw := range payload.Properties {
		val, err := j.unmarshalValue(raw)
		if err != nil {
			return err
		}
		v.SetProperty(k, val)
	}
	return nil
}

// 辅助方法：根据 JSON 数据推断类型并反序列化
func (j *JsonSerializer) unmarshalValue(data []byte) (data.Value, error) {
	// 尝试按顺序解析为不同类型
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		return data.NewIntValue(i), nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		return data.NewFloatValue(f), nil
	}

	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		return data.NewBoolValue(b), nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return data.NewStringValue(s), nil
	}

	// 检查是否为 null
	if string(data) == "null" {
		return data.NewNullValue(), nil
	}

	// 检查是否为数组
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err == nil {
		av := &data.ArrayValue{}
		if err := j.UnmarshalArray(data, av); err != nil {
			return nil, err
		}
		return av, nil
	}

	// 检查是否为对象
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err == nil {
		ov := data.NewObjectValue()
		if err := j.UnmarshalObject(data, ov); err != nil {
			return nil, err
		}
		return ov, nil
	}

	// 如果都不匹配，作为任意值处理
	var anyVal any
	if err := json.Unmarshal(data, &anyVal); err != nil {
		return nil, fmt.Errorf("无法解析 JSON 数据: %v", err)
	}
	return data.NewAnyValue(anyVal), nil
}
