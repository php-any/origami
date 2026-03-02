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

// CloneObjectValue 创建一个新的 ObjectValue 副本。
// 语义上用于模拟 PHP 数组的 copy-on-write（特别是关联数组在 Origami 中可能用 ObjectValue 表示）：
// - 仅复制属性存储结构本身（PropertyStore），不深拷贝每个元素的 Value
// - 这样结构性修改（新增/覆盖某个 key）不会影响原对象
// - 元素内部若是对象/数组，仍按其自身语义共享
func CloneObjectValue(src *ObjectValue) *ObjectValue {
	if src == nil {
		return nil
	}

	clone := &ObjectValue{
		Value:    src.Value,
		Context:  src.Context,
		property: NewOrderedMap(),
	}

	// 按插入顺序复制所有键值对
	src.property.Range(func(key string, value Value) bool {
		clone.property.Set(key, value)
		return true
	})

	return clone
}

type ObjectValue struct {
	Value
	Context
	property PropertyStore
	iterator int // 迭代器当前位置索引
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

func (o *ObjectValue) GetProperty(name string) (Value, Control) {
	v, ok := o.property.Get(name)
	if !ok {
		switch name {
		case "length":
			return NewIntValue(o.property.Len()), nil
		}
		return NewNullValue(), nil
	}
	return v, nil
}

func (o *ObjectValue) GetZVal(name string) (*ZVal, Control) {
	v, _ := o.property.GetZVal(name)
	return v, nil
}

func (o *ObjectValue) SetProperty(name string, value Value) Control {
	// 为了更贴近 PHP 数组的 copy-on-write 语义，当属性值是数组时存储一个克隆，
	// 避免多个属性/变量共享同一个 ArrayValue 实例，被 array_shift/array_pop 等原地修改时互相影响。
	if arr, ok := value.(*ArrayValue); ok {
		value = CloneArrayValue(arr)
	}
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

// RangeProperties 按插入顺序遍历所有属性
// 使用此方法可保证遍历顺序与插入顺序一致，避免 Go map 遍历顺序随机的问题
func (o *ObjectValue) RangeProperties(fn func(key string, value Value) bool) {
	o.property.Range(fn)
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

func (o *ObjectValue) ToGoValue(serializer Serializer) (any, error) {
	return serializer.MarshalObject(o)
}

// Iterator 接口实现

// Rewind 将指针重置到第一个元素
func (o *ObjectValue) Rewind(_ Context) (Value, Control) {
	o.iterator = 0
	return nil, nil
}

// Valid 检查当前位置是否有效（是否还有元素）
func (o *ObjectValue) Valid(_ Context) (Value, Control) {
	count := o.property.Len()
	valid := o.iterator >= 0 && o.iterator < count
	return NewBoolValue(valid), nil
}

// Current 返回当前元素
func (o *ObjectValue) Current(_ Context) (Value, Control) {
	_, value, ok := o.property.GetByIndex(o.iterator)
	if !ok {
		return NewNullValue(), nil
	}
	return value, nil
}

// Key 返回当前元素的键
func (o *ObjectValue) Key(_ Context) (Value, Control) {
	key, _, ok := o.property.GetByIndex(o.iterator)
	if !ok {
		return NewNullValue(), nil
	}
	return NewStringValue(key), nil
}

// Next 将指针向前移动到下一个元素
func (o *ObjectValue) Next(_ Context) Control {
	o.iterator++
	return nil
}
