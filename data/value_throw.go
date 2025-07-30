package data

import "fmt"

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
	error      *Error
	extend     string
	getMessage Method
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

func (t *ThrowValue) GetProperties() map[string]Property {
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
	return &ThrowValue{
		error: NewError(from, err.Error(), err),
	}
}

// TryErrorThrow 可能不需要抛出的错误
func TryErrorThrow(from From, err error) Control {
	return &ThrowValue{
		error: NewError(from, err.Error(), err),
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
	return t.error
}

// AsString 获取字符串表示
func (t *ThrowValue) AsString() string {
	return fmt.Sprintf("throw %v", t.error.Error())
}
