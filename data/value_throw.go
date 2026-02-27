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

type ThrowValue struct {
	Object           *ClassValue
	getMessage       Method
	getTraceAsString Method

	// previous 表示前一个异常，用于异常链；目前 Origami 尚未维护完整链，
	// 该字段暂未在 NewErrorThrow* 中赋值，getPrevious() 将返回 null。
	previous *ThrowValue

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
	if t.Object == nil {
		return t.Name
	}

	return t.Object.Class.GetName()
}

func (t *ThrowValue) GetExtend() *string {
	if t.Object != nil {
		return t.Object.Class.GetExtend()
	}
	return nil
}

func (t *ThrowValue) GetImplements() []string {
	if t.Object != nil {
		return t.Object.Class.GetImplements()
	}
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
	case "getPrevious":
		return &ThrowValueGetPreviousMethod{source: t}, true
	case "getCode":
		return &ThrowValueGetCodeMethod{source: t}, true
	case "getFile":
		return &ThrowValueGetFileMethod{source: t}, true
	case "getLine":
		return &ThrowValueGetLineMethod{source: t}, true
	case "getTrace":
		return &ThrowValueGetTraceMethod{source: t}, true
	}
	return nil, false
}

func (t *ThrowValue) GetMethods() []Method {
	return []Method{
		t.getMessage,
		t.getTraceAsString,
		&ThrowValueGetPreviousMethod{source: t},
		&ThrowValueGetCodeMethod{source: t},
		&ThrowValueGetFileMethod{source: t},
		&ThrowValueGetLineMethod{source: t},
		&ThrowValueGetTraceMethod{source: t},
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
		switch ret.(type) {
		case Value:
			err = ret.(Value).AsString()
		default:
			panic(ret) // 不可能执行到这里的, 如果有报错就是解析器的问题
		}
	} else {
		classStmt := object.Class
		for classStmt != nil {
			if method, ok := classStmt.GetMethod("getMessage"); ok {
				ret, acl := method.Call(object)
				if acl != nil {
					return acl
				}
				switch ret.(type) {
				case Value:
					err = ret.(Value).AsString()
				}
				classStmt = nil
			} else {
				if classStmt.GetExtend() != nil {
					classStmt, _ = object.GetVM().GetOrLoadClass(*classStmt.GetExtend())
				} else {
					classStmt = nil
				}
			}
		}
		if err == "" {
			err = "运行时无法处理未继承 Exception 的异常类"
		}
	}

	t := &ThrowValue{
		Object: object,
		Name:   object.GetName(),
		Error:  NewError(from, err, nil),
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

func NewErrorThrowByName(from From, err error, name string) Control {
	t := &ThrowValue{
		Error: NewError(from, err.Error(), err),
		Name:  name,
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

// ---------------- 额外的 Throwable 方法实现 ----------------

// getPrevious(): ?Throwable
type ThrowValueGetPreviousMethod struct {
	source *ThrowValue
}

func (m *ThrowValueGetPreviousMethod) Call(ctx Context) (GetValue, Control) {
	if m.source == nil || m.source.previous == nil {
		return NewNullValue(), nil
	}
	return m.source.previous, nil
}

func (m *ThrowValueGetPreviousMethod) GetName() string       { return "getPrevious" }
func (m *ThrowValueGetPreviousMethod) GetModifier() Modifier { return ModifierPublic }
func (m *ThrowValueGetPreviousMethod) GetIsStatic() bool     { return false }
func (m *ThrowValueGetPreviousMethod) GetParams() []GetValue { return []GetValue{} }
func (m *ThrowValueGetPreviousMethod) GetVariables() []Variable {
	return []Variable{}
}

func (m *ThrowValueGetPreviousMethod) GetReturnType() Types {
	return NewBaseType("Throwable")
}

// getCode(): int
type ThrowValueGetCodeMethod struct {
	source *ThrowValue
}

func (m *ThrowValueGetCodeMethod) Call(ctx Context) (GetValue, Control) {
	// 目前 Error 中未区分错误码，这里统一返回 0
	return NewIntValue(0), nil
}

func (m *ThrowValueGetCodeMethod) GetName() string       { return "getCode" }
func (m *ThrowValueGetCodeMethod) GetModifier() Modifier { return ModifierPublic }
func (m *ThrowValueGetCodeMethod) GetIsStatic() bool     { return false }
func (m *ThrowValueGetCodeMethod) GetParams() []GetValue { return []GetValue{} }
func (m *ThrowValueGetCodeMethod) GetVariables() []Variable {
	return []Variable{}
}

func (m *ThrowValueGetCodeMethod) GetReturnType() Types {
	return NewBaseType("int")
}

// getFile(): string
type ThrowValueGetFileMethod struct {
	source *ThrowValue
}

func (m *ThrowValueGetFileMethod) Call(ctx Context) (GetValue, Control) {
	if m.source == nil || m.source.Error == nil || m.source.Error.From == nil {
		return NewStringValue(""), nil
	}
	return NewStringValue(m.source.Error.From.GetSource()), nil
}

func (m *ThrowValueGetFileMethod) GetName() string       { return "getFile" }
func (m *ThrowValueGetFileMethod) GetModifier() Modifier { return ModifierPublic }
func (m *ThrowValueGetFileMethod) GetIsStatic() bool     { return false }
func (m *ThrowValueGetFileMethod) GetParams() []GetValue { return []GetValue{} }
func (m *ThrowValueGetFileMethod) GetVariables() []Variable {
	return []Variable{}
}

func (m *ThrowValueGetFileMethod) GetReturnType() Types {
	return NewBaseType("string")
}

// getLine(): int
type ThrowValueGetLineMethod struct {
	source *ThrowValue
}

func (m *ThrowValueGetLineMethod) Call(ctx Context) (GetValue, Control) {
	if m.source == nil || m.source.Error == nil || m.source.Error.From == nil {
		return NewIntValue(0), nil
	}
	sl, _ := m.source.Error.From.GetStartPosition()
	return NewIntValue(sl + 1), nil
}

func (m *ThrowValueGetLineMethod) GetName() string       { return "getLine" }
func (m *ThrowValueGetLineMethod) GetModifier() Modifier { return ModifierPublic }
func (m *ThrowValueGetLineMethod) GetIsStatic() bool     { return false }
func (m *ThrowValueGetLineMethod) GetParams() []GetValue { return []GetValue{} }
func (m *ThrowValueGetLineMethod) GetVariables() []Variable {
	return []Variable{}
}

func (m *ThrowValueGetLineMethod) GetReturnType() Types {
	return NewBaseType("int")
}

// getTrace(): array
type ThrowValueGetTraceMethod struct {
	source *ThrowValue
}

func (m *ThrowValueGetTraceMethod) Call(ctx Context) (GetValue, Control) {
	frames := make([]Value, 0, len(m.source.StackFrames))

	for _, frame := range m.source.StackFrames {
		obj := NewObjectValue()
		if frame.From != nil {
			sl, _ := frame.From.GetStartPosition()
			obj.SetProperty("file", NewStringValue(frame.From.GetSource()))
			obj.SetProperty("line", NewIntValue(sl+1))
		}
		if frame.ClassName != "" {
			obj.SetProperty("class", NewStringValue(frame.ClassName))
		}
		if frame.MethodName != "" {
			obj.SetProperty("function", NewStringValue(frame.MethodName))
			// 简化：统一认为是对象方法
			obj.SetProperty("type", NewStringValue("->"))
		}
		frames = append(frames, obj)
	}

	return NewArrayValue(frames), nil
}

func (m *ThrowValueGetTraceMethod) GetName() string       { return "getTrace" }
func (m *ThrowValueGetTraceMethod) GetModifier() Modifier { return ModifierPublic }
func (m *ThrowValueGetTraceMethod) GetIsStatic() bool     { return false }
func (m *ThrowValueGetTraceMethod) GetParams() []GetValue { return []GetValue{} }
func (m *ThrowValueGetTraceMethod) GetVariables() []Variable {
	return []Variable{}
}

func (m *ThrowValueGetTraceMethod) GetReturnType() Types {
	return NewBaseType("array")
}
