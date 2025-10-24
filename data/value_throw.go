package data

import (
	"fmt"
)

// ThrowControl 表示异常抛出控制流
type ThrowControl interface {
	Control
	// IsThrow 是否为异常抛出
	IsThrow() bool
	// GetError 获取异常信息
	GetError() *Error
}

// ThrowValue 表示异常抛出控制流
type ThrowValue struct {
	object     *ClassValue
	extend     string
	getMessage Method

	Error *Error
	// 堆栈
	Stack []From
}

func (t *ThrowValue) AddStack(f From) {
	t.Stack = append(t.Stack, f)
}

func (t *ThrowValue) GetFrom() From {
	return t.Error.From
}

func (t *ThrowValue) GetName() string {
	return "Error"
}

func (t *ThrowValue) GetExtend() *string {
	return &t.extend
}

func (t *ThrowValue) GetImplements() []string {
	return nil
}

func (t *ThrowValue) GetProperty(name string) (Property, bool) {
	return nil, false
}

func (t *ThrowValue) GetPropertyList() []Property {
	return nil
}

func (t *ThrowValue) GetMethod(name string) (Method, bool) {
	switch name {
	case "getMessage":
		return t.getMessage, true
	}
	return nil, false
}

func (t *ThrowValue) GetMethods() []Method {
	return nil
}

func (t *ThrowValue) GetConstruct() Method {
	return nil
}

func NewErrorThrow(from From, err error) Control {
	t := &ThrowValue{
		Error: NewError(from, err.Error(), err),
	}
	t.getMessage = &ThrowValueGetMessageMethod{
		source: t,
	}
	return t
}

// TryErrorThrow 可能不需要抛出的错误
func TryErrorThrow(from From, err error) Control {
	return &ThrowValue{
		Error: NewError(from, err.Error(), err),
	}
}

func (t *ThrowValue) SetValue(v Value) {
	panic("TODO ThrowValue.SetValue")
}

func (t *ThrowValue) GetValue(ctx Context) (GetValue, Control) {
	return t, nil
}

// IsThrow 是否为异常抛出
func (t *ThrowValue) IsThrow() bool {
	return true
}

// GetError 获取异常信息
func (t *ThrowValue) GetError() *Error {
	return t.Error
}

// AsString 获取字符串表示
func (t *ThrowValue) AsString() string {
	return fmt.Sprintf("throw %v", t.Error.Error())
}

type ThrowValueGetMessageMethod struct {
	source *ThrowValue
}

func (t *ThrowValueGetMessageMethod) Call(ctx Context) (GetValue, Control) {
	return NewStringValue(t.source.Error.Error()), nil
}

func (t *ThrowValueGetMessageMethod) GetName() string {
	return "getMessage"
}

func (t *ThrowValueGetMessageMethod) GetModifier() Modifier {
	return ModifierPublic
}

func (t *ThrowValueGetMessageMethod) GetIsStatic() bool {
	return false
}

func (t *ThrowValueGetMessageMethod) GetParams() []GetValue {
	return []GetValue{}
}

func (t *ThrowValueGetMessageMethod) GetVariables() []Variable {
	return []Variable{}
}

// GetReturnType 返回方法返回类型
func (t *ThrowValueGetMessageMethod) GetReturnType() Types {
	return NewBaseType("string")
}
