package data

import (
	"fmt"
)

func NewClassValue(class ClassStmt, ctx Context) *ClassValue {
	// 开始初始化属性
	return &ClassValue{
		ObjectValue: NewObjectValue(),
		Class:       class,
		Context:     ctx,
	}
}

type AsClass interface {
	AsObject
}

type ClassValue struct {
	Context
	*ObjectValue
	Class ClassStmt
}

func (c *ClassValue) GetValue(ctx Context) (GetValue, Control) {
	return c, nil
}

func (c *ClassValue) AsString() string {
	result := ""
	c.property.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(Value)
		result += fmt.Sprintf("%s: %s\n", k, v.AsString())
		return true
	})

	if len(result) > 2 {
		result = result[:len(result)-1] // 移除最后一个换行符
	}

	// 构建输出字符串
	return fmt.Sprintf("%s {\n"+
		"\t%s\n"+
		"}",
		c.Class.GetName(), result,
	)
}

func (c *ClassValue) GetProperty(name string) (Property, bool) {
	if v, ok := c.Class.GetProperty(name); ok {
		return v, true
	}

	vm := c.GetVM()
	// 执行父级
	last := c.Class
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		next, ok := vm.GetClass(*ext)
		if !ok {
			return nil, false
		}

		property, ok := next.GetProperty(name)
		if ok {
			if property.GetModifier() == ModifierPrivate {
				return nil, false
			}
			return property, true
		}
		last = next
	}

	return nil, false
}

func (c *ClassValue) GetMethod(name string) (Method, bool) {
	if fn, ok := c.Class.GetMethod(name); ok {
		return fn, true
	}

	vm := c.GetVM()
	// 执行父级
	last := c.Class
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		next, ok := vm.GetClass(*ext)
		if !ok {
			return nil, false
		}

		fn, ok := next.GetMethod(name)
		if ok {
			if fn.GetModifier() == ModifierPrivate {
				return nil, false
			}
			return fn, true
		}
		last = next
	}

	return nil, false
}

func (c *ClassValue) GetProperties() map[string]Value {
	return map[string]Value{}
}

func (c *ClassValue) CreateContext(vars []Variable) Context {
	ctx := c.Context.CreateContext(vars)
	return &ClassMethodContext{
		&ClassValue{
			ObjectValue: c.ObjectValue,
			Class:       c.Class,
			Context:     ctx,
		},
	}
}

func (c *ClassValue) SetVariableValue(variable Variable, value Value) Control {
	c.SetProperty(variable.GetName(), value)
	return nil
}

func (c *ClassValue) GetVariableValue(variable Variable) (Value, Control) {
	return c.ObjectValue.GetVariableValue(variable)
}

type ClassMethodContext struct {
	*ClassValue
}

func (c *ClassMethodContext) SetVariableValue(variable Variable, value Value) Control {
	c.Context.SetVariableValue(variable, value)
	return nil
}

func (c *ClassMethodContext) GetVariableValue(variable Variable) (Value, Control) {
	if _, ok := variable.(Property); ok {
		return c.ObjectValue.GetVariableValue(variable)
	}
	return c.Context.GetVariableValue(variable)
}

func (c *ClassMethodContext) GetIndexValue(index int) (Value, bool) {
	return c.Context.GetIndexValue(index)
}
