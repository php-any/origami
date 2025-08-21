package data

import (
	"fmt"
	"sync"
)

type AsObject interface {
}

func NewObjectValue() *ObjectValue {
	return &ObjectValue{
		property: sync.Map{},
	}
}

type ObjectValue struct {
	Value
	Context
	property sync.Map
}

func (o *ObjectValue) GetValue(ctx Context) (GetValue, Control) {
	return o, nil
}

func (o *ObjectValue) AsString() string {
	result := ""
	o.property.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(Value)
		result += fmt.Sprintf("\t%s: %s\n", k, v.AsString())
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
	notEmpty := false
	o.property.Range(func(key, value any) bool {
		notEmpty = true
		return false
	})
	return notEmpty, nil
}

func (o *ObjectValue) GetProperty(name string) (Value, bool) {
	v, ok := o.property.Load(name)
	if !ok {
		return NewNullValue(), false
	}
	return v.(Value), ok
}

func (o *ObjectValue) SetProperty(name string, value Value) {
	o.property.Store(name, value)
}

func (o *ObjectValue) DeleteProperty(name string) {
	o.property.Delete(name)
}

func (o *ObjectValue) HasProperty(name string) bool {
	_, ok := o.property.Load(name)
	return ok
}

func (o *ObjectValue) GetProperties() map[string]Value {
	properties := make(map[string]Value)

	// 遍历 sync.Map 中的所有属性
	o.property.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(Value)
		properties[k] = v
		return true
	})

	return properties
}

func (o *ObjectValue) SetVariableValue(variable Variable, value Value) Control {
	o.SetProperty(variable.GetName(), value)
	return nil
}

func (o *ObjectValue) GetVariableValue(variable Variable) (Value, Control) {
	v, ok := o.property.Load(variable.GetName())
	if !ok {
		return nil, nil
	}
	return v.(Value), nil
}
