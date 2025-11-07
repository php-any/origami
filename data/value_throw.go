package data

import (
	"fmt"
	"strings"
)

// ThrowControl 表示异常抛出控制流
type ThrowControl interface {
	Control
	// IsThrow 是否为异常抛出
	IsThrow() bool
	// GetError 获取异常信息
	GetError() *Error
}

// StackFrame 表示调用栈帧
type StackFrame struct {
	From       From
	ClassName  string
	MethodName string
}

// ThrowValue 表示异常抛出控制流
type ThrowValue struct {
	object           *ClassValue
	getMessage       Method
	getTraceAsString Method

	Name string

	Error *Error
	// 调用栈信息
	StackFrames []StackFrame
}

// AddStackWithInfo 添加调用栈信息，包含类名和方法名
func (t *ThrowValue) AddStackWithInfo(f From, className, methodName string) {
	// 添加到详细调用栈信息
	frame := StackFrame{
		From:       f,
		ClassName:  className,
		MethodName: methodName,
	}
	t.StackFrames = append(t.StackFrames, frame)
}

func (t *ThrowValue) GetFrom() From {
	return t.Error.From
}

func (t *ThrowValue) GetName() string {
	if t.object == nil {
		return t.Name
	}

	return t.object.Class.GetName()
}

func (t *ThrowValue) GetExtend() *string {
	return nil
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
	case "getTraceAsString":
		return t.getTraceAsString, true
	}
	return nil, false
}

func (t *ThrowValue) GetMethods() []Method {
	return []Method{
		t.getMessage,
		t.getTraceAsString,
	}
}

func (t *ThrowValue) GetConstruct() Method {
	return nil
}

func NewErrorThrowFromClassValue(from From, object *ClassValue) Control {
	err := ""
	if str, ok := object.Class.(AsString); ok {
		err = str.AsString()
	} else if method, ok := object.GetMethod("getMessage"); ok {
		ret, acl := method.Call(object)
		if acl != nil {
			return acl
		}
		err = ret.(Value).AsString()
	} else {
		err = "运行时无法处理未继承 Exception 的异常类"
	}

	t := &ThrowValue{
		object: object,
		Name:   "Exception",
		Error:  NewError(from, fmt.Sprintf("Throw %s: %s", object.Class.GetName(), err), nil),
	}
	t.getMessage = &ThrowValueGetMessageMethod{
		source: t,
	}
	t.getTraceAsString = &ThrowValueGetTraceAsStringMethod{
		source: t,
	}
	return t
}

func NewErrorThrow(from From, err error) Control {
	t := &ThrowValue{
		Error: NewError(from, err.Error(), err),
		Name:  "Exception",
	}
	t.getMessage = &ThrowValueGetMessageMethod{
		source: t,
	}
	t.getTraceAsString = &ThrowValueGetTraceAsStringMethod{
		source: t,
	}
	return t
}

// TryErrorThrow 可能不需要抛出的错误
func TryErrorThrow(from From, err error) Control {
	t := &ThrowValue{
		Error: NewError(from, err.Error(), err),
		Name:  "Exception",
	}
	t.getMessage = &ThrowValueGetMessageMethod{
		source: t,
	}
	t.getTraceAsString = &ThrowValueGetTraceAsStringMethod{
		source: t,
	}
	return t
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

type ThrowValueGetTraceAsStringMethod struct {
	source *ThrowValue
}

func (t *ThrowValueGetTraceAsStringMethod) Call(ctx Context) (GetValue, Control) {
	// 构建堆栈跟踪信息
	var trace strings.Builder
	trace.WriteString("Stack trace:\n")

	// 添加当前异常信息
	if t.source.Error != nil && t.source.Error.From != nil {
		start, end := t.source.Error.From.GetPosition()
		sl, sp := t.source.Error.From.GetStartPosition()
		trace.WriteString(fmt.Sprintf("  at %s:%d:%d (position %d-%d)\n",
			t.source.Error.From.GetSource(), sl+1, sp+1, start, end))
	}

	// 添加调用栈信息
	for _, frame := range t.source.StackFrames {
		if frame.From != nil {
			start, end := frame.From.GetPosition()
			sl, sp := frame.From.GetStartPosition()
			trace.WriteString(fmt.Sprintf("  at %s.%s() in %s:%d:%d (position %d-%d)\n",
				frame.ClassName, frame.MethodName, frame.From.GetSource(), sl+1, sp+1, start, end))
		}
	}

	return NewStringValue(trace.String()), nil
}

func (t *ThrowValueGetTraceAsStringMethod) GetName() string {
	return "getTraceAsString"
}

func (t *ThrowValueGetTraceAsStringMethod) GetModifier() Modifier {
	return ModifierPublic
}

func (t *ThrowValueGetTraceAsStringMethod) GetIsStatic() bool {
	return false
}

func (t *ThrowValueGetTraceAsStringMethod) GetParams() []GetValue {
	return []GetValue{}
}

func (t *ThrowValueGetTraceAsStringMethod) GetVariables() []Variable {
	return []Variable{}
}

// GetReturnType 返回方法返回类型
func (t *ThrowValueGetTraceAsStringMethod) GetReturnType() Types {
	return NewBaseType("string")
}
