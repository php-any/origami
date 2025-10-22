package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
)

// DatabaseSerializer 实现 data.Serializer 接口，用于数据库相关的序列化
type DatabaseSerializer struct{}

// NewDatabaseSerializer 创建数据库序列化器
func NewDatabaseSerializer() *DatabaseSerializer {
	return &DatabaseSerializer{}
}

// 实现 data.Serializer 接口的所有方法
func (ds *DatabaseSerializer) MarshalInt(v *data.IntValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (ds *DatabaseSerializer) UnmarshalInt(data []byte, v *data.IntValue) error {
	return json.Unmarshal(data, &v.Value)
}

func (ds *DatabaseSerializer) MarshalString(v *data.StringValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (ds *DatabaseSerializer) UnmarshalString(data []byte, v *data.StringValue) error {
	return json.Unmarshal(data, &v.Value)
}

func (ds *DatabaseSerializer) MarshalNull(v *data.NullValue) ([]byte, error) {
	return json.Marshal(nil)
}

func (ds *DatabaseSerializer) UnmarshalNull(data []byte, v *data.NullValue) error {
	return nil
}

func (ds *DatabaseSerializer) MarshalArray(v *data.ArrayValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (ds *DatabaseSerializer) UnmarshalArray(data []byte, v *data.ArrayValue) error {
	return json.Unmarshal(data, &v.Value)
}

func (ds *DatabaseSerializer) MarshalObject(v *data.ObjectValue) ([]byte, error) {
	props := v.GetProperties()
	encoded := make(map[string]json.RawMessage, len(props))
	for k, val := range props {
		if vs, ok := val.(data.ValueSerializer); ok {
			b, err := vs.Marshal(ds)
			if err != nil {
				return nil, err
			}
			encoded[k] = b
		} else {
			b, _ := json.Marshal(val.AsString())
			encoded[k] = b
		}
	}
	return json.Marshal(encoded)
}

func (ds *DatabaseSerializer) UnmarshalObject(data []byte, v *data.ObjectValue) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, raw := range m {
		val, err := ds.unmarshalValue(raw)
		if err != nil {
			return err
		}
		v.SetProperty(k, val)
	}
	return nil
}

func (ds *DatabaseSerializer) MarshalBool(v *data.BoolValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (ds *DatabaseSerializer) UnmarshalBool(data []byte, v *data.BoolValue) error {
	return json.Unmarshal(data, &v.Value)
}

func (ds *DatabaseSerializer) MarshalFloat(v *data.FloatValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (ds *DatabaseSerializer) UnmarshalFloat(data []byte, v *data.FloatValue) error {
	return json.Unmarshal(data, &v.Value)
}

func (ds *DatabaseSerializer) MarshalAny(v *data.AnyValue) ([]byte, error) {
	return json.Marshal(v.Value)
}

func (ds *DatabaseSerializer) UnmarshalAny(data []byte, v *data.AnyValue) error {
	var tmp any
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Value = tmp
	return nil
}

func (ds *DatabaseSerializer) MarshalClass(v *data.ClassValue) ([]byte, error) {
	props := v.GetProperties()
	encoded := make(map[string]json.RawMessage, len(props))
	for k, val := range props {
		if vs, ok := val.(data.ValueSerializer); ok {
			b, err := vs.Marshal(ds)
			if err != nil {
				return nil, err
			}
			encoded[k] = b
		} else {
			b, _ := json.Marshal(val.AsString())
			encoded[k] = b
		}
	}
	return json.Marshal(encoded)
}

func (ds *DatabaseSerializer) UnmarshalClass(data []byte, v *data.ClassValue) error {
	var props map[string]json.RawMessage
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	for k, raw := range props {
		val, err := ds.unmarshalValue(raw)
		if err != nil {
			return err
		}
		v.SetProperty(k, val)
	}
	return nil
}

// unmarshalValue 根据 JSON 数据推断类型并反序列化
func (ds *DatabaseSerializer) unmarshalValue(raw []byte) (data.Value, error) {
	// 尝试按顺序解析为不同类型
	var i int
	if err := json.Unmarshal(raw, &i); err == nil {
		return data.NewIntValue(i), nil
	}

	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return data.NewFloatValue(f), nil
	}

	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return data.NewBoolValue(b), nil
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return data.NewStringValue(s), nil
	}

	// 检查是否为 null
	if string(raw) == "null" {
		return data.NewNullValue(), nil
	}

	// 检查是否为数组
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err == nil {
		av := &data.ArrayValue{}
		if err := ds.UnmarshalArray(raw, av); err != nil {
			return nil, err
		}
		return av, nil
	}

	// 检查是否为对象
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err == nil {
		ov := data.NewObjectValue()
		if err := ds.UnmarshalObject(raw, ov); err != nil {
			return nil, err
		}
		return ov, nil
	}

	// 如果都不匹配，作为任意值处理
	var anyVal any
	if err := json.Unmarshal(raw, &anyVal); err != nil {
		return nil, fmt.Errorf("无法解析 JSON 数据: %v", err)
	}
	return data.NewAnyValue(anyVal), nil
}

// 数据库扫描相关的辅助方法

// ScanRowToInstance 扫描数据库行到实例
func (ds *DatabaseSerializer) ScanRowToInstance(instance *data.ClassValue, row *sql.Row) error {
	// 获取实例的类定义
	classStmt := instance.Class
	if classStmt == nil {
		return fmt.Errorf("实例没有类定义")
	}

	// 获取类的属性定义
	properties := classStmt.GetPropertyList()
	if len(properties) == 0 {
		return fmt.Errorf("类没有属性定义")
	}

	// 创建扫描目标
	values := make([]interface{}, len(properties))
	valuePtrs := make([]interface{}, len(properties))

	// 为每个属性创建扫描目标
	propertyNames := make([]string, 0, len(properties))
	for _, property := range properties {
		propertyNames = append(propertyNames, property.GetName())
	}

	// 初始化扫描目标
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 扫描数据
	err := row.Scan(valuePtrs...)
	if err != nil {
		return fmt.Errorf("扫描数据库行失败: %w", err)
	}

	// 将扫描结果映射到实例属性
	for i, val := range values {
		if i >= len(propertyNames) {
			break
		}

		propertyName := propertyNames[i]
		value := ds.convertToValue(val)
		instance.SetProperty(propertyName, value)
	}

	return nil
}

// convertToValue 将数据库值转换为脚本值
func (ds *DatabaseSerializer) convertToValue(val interface{}) data.Value {
	if val == nil {
		return data.NewNullValue()
	}

	switch v := val.(type) {
	case int:
		return data.NewIntValue(v)
	case int8:
		return data.NewIntValue(int(v))
	case int16:
		return data.NewIntValue(int(v))
	case int32:
		return data.NewIntValue(int(v))
	case int64:
		// 检查是否超出 int 范围
		if v > int64(^uint(0)>>1) || v < int64(-1<<63) {
			// 如果超出 int 范围，转换为字符串
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case uint:
		if v > uint(^uint(0)>>1) {
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case uint8:
		return data.NewIntValue(int(v))
	case uint16:
		return data.NewIntValue(int(v))
	case uint32:
		return data.NewIntValue(int(v))
	case uint64:
		if v > uint64(^uint(0)>>1) {
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case float32:
		return data.NewFloatValue(float64(v))
	case float64:
		return data.NewFloatValue(v)
	case string:
		return data.NewStringValue(v)
	case []byte:
		return data.NewStringValue(string(v))
	case bool:
		return data.NewBoolValue(v)
	default:
		// 对于其他类型，转换为字符串
		return data.NewStringValue(fmt.Sprintf("%v", v))
	}
}
