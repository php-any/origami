package data

import (
	"context"
	"fmt"
)

type AsObject interface {
}

func NewObjectValue() *ObjectValue {
	return &ObjectValue{
		property: NewOrderedMap(),
	}
}

type ObjectValue struct {
	Value
	Context
	property *OrderedMap
}

func (o *ObjectValue) GoContext() context.Context {
	return context.Background()
}

func (o *ObjectValue) GetValue(ctx Context) (GetValue, Control) {
	return o, nil
}

func (o *ObjectValue) AsString() string {
	result := ""
	o.property.Range(func(key string, value Value) bool {
		result += fmt.Sprintf("\t%s: %s\n", key, value.AsString())
		return true
	})
	if len(result) > 2 {
		result = result[:len(result)-1] // 移除最后一个换行符
	}

	// 构建输出字符串
	return fmt.Sprintf("Object {\n"+
		"%s\n"+
		"}",
		result,
	)
}

func (o *ObjectValue) AsBool() (bool, error) {
	return true, nil
}

func (o *ObjectValue) GetProperty(name string) (Value, bool) {
	v, ok := o.property.Get(name)
	if !ok {
		return NewNullValue(), false
	}
	return v, ok
}

func (o *ObjectValue) SetProperty(name string, value Value) Control {
	o.property.Set(name, value)
	return nil
}

func (o *ObjectValue) GetProperties() map[string]Value {
	properties := make(map[string]Value)

	// 遍历 OrderedMap 中的所有属性
	o.property.Range(func(key string, value Value) bool {
		properties[key] = value
		return true
	})

	return properties
}

func (o *ObjectValue) SetVariableValue(variable Variable, value Value) Control {
	o.SetProperty(variable.GetName(), value)
	return nil
}

func (o *ObjectValue) GetVariableValue(variable Variable) (Value, Control) {
	v, ok := o.property.Get(variable.GetName())
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (o *ObjectValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalObject(o)
}

func (o *ObjectValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalObject(data, o)
}
