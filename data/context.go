package data

import "errors"

type Context interface {
	SetNamespace(name string) Context
	GetNamespace() string

	GetVariableValue(variable Variable) (Value, Control)
	GetIndexValue(index int) (Value, bool)
	SetVariableValue(variable Variable, value Value) Control

	CreateContext(vars []Variable) Context

	CreateBaseContext() Context

	GetVM() VM
}

type VM interface {
	AddClass(c ClassStmt) Control
	GetClass(pkg string) (ClassStmt, bool)
	AddInterface(i InterfaceStmt) Control
	GetInterface(pkg string) (InterfaceStmt, bool)
	AddFunc(f FuncStmt) Control
	GetFunc(pkg string) (FuncStmt, bool)
	RegisterFunction(name string, fn interface{}) Control
	RegisterReflectClass(name string, instance interface{}) Control
	CreateContext(vars []Variable) Context
	SetThrowControl(func(acl Control))
	ThrowControl(acl Control)
	LoadAndRun(file string) (GetValue, Control)

	SetClassPathCache(name, path string)
	GetClassPathCache(name string) (string, bool)
}

type ClassStmt interface {
	GetValue(ctx Context) (GetValue, Control)
	GetFrom() From
	GetName() string

	GetExtend() *string                       // 父类名
	GetImplements() []string                  // 实现的接口列表
	GetProperty(name string) (Property, bool) // 属性列表
	GetProperties() map[string]Property       // 属性列表
	GetMethod(name string) (Method, bool)     // 方法列表
	GetMethods() []Method                     // 方法列表

	GetConstruct() Method
}

type Modifier int

const (
	ModifierPublic    Modifier = iota // 公有
	ModifierPrivate                   // 私有
	ModifierProtected                 // 保护
)

func NewModifier(v string) Modifier {
	switch v {
	case "private":
		return ModifierPrivate
	case "protected":
		return ModifierProtected
	default:
		return ModifierPublic
	}
}

type Property interface {
	GetValue(ctx Context) (GetValue, Control)
	GetName() string           // 属性名
	GetModifier() Modifier     // 访问修饰符
	GetIsStatic() bool         // 是否是静态属性
	GetDefaultValue() GetValue // 默认值
}

type Method interface {
	Call(ctx Context) (GetValue, Control)
	GetName() string          // 方法名
	GetModifier() Modifier    // 访问修饰符
	GetIsStatic() bool        // 是否是静态方法
	GetParams() []GetValue    // 参数列表, 用于接收参数; 实现这个接口的参数结构体 node.Parameter 接收单个值, node.Parameters 只有任意多个值
	GetVariables() []Variable // 变量列表, 用于创建符号表
}

type FuncStmt interface {
	Call(ctx Context) (GetValue, Control)
	GetName() string
	// GetParams 获取参数列表
	GetParams() []GetValue // 参数列表
	// GetVariables 获取变量列表, 用于创建符号表
	GetVariables() []Variable // 变量列表
}

type Function struct {
	FuncStmt
}

func (f Function) GetVariables() []Variable {
	return nil
}

func (f Function) GetParams() []GetValue {
	var got []GetValue
	for _, value := range f.GetVariables() {
		got = append(got, value)
	}

	return got
}

type Variable interface {
	GetValue

	GetIndex() int
	GetName() string
	GetType() Types

	SetValue(ctx Context, value Value) Control
}

func NewVariable(name string, index int, ty Types) Variable {
	return &VariableTODO{name: name, index: index, ty: ty}
}

type VariableTODO struct {
	name  string
	index int
	ty    Types
}

func (v VariableTODO) GetValue(_ Context) (GetValue, Control) {
	return nil, nil
}

func (v VariableTODO) GetIndex() int {
	return v.index
}

func (v VariableTODO) GetName() string {
	return v.name
}

func (v VariableTODO) GetType() Types {
	return v.ty
}

func (v VariableTODO) SetValue(ctx Context, value Value) Control {
	if v.ty == nil {
		return ctx.SetVariableValue(v, value)
	}
	if v.ty.Is(value) {
		return ctx.SetVariableValue(v, value)
	}
	return NewErrorThrow(nil, errors.New("变量类型和赋值类型不一致, 变量类型("+v.ty.String()+"), 赋值("+value.AsString()+")"))
}

type Parameter interface {
	GetIndex() int
	GetType() Types
	SetValue(ctx Context, value Value) Control
	GetName() string
	GetValue(ctx Context) (GetValue, Control)
	GetDefaultValue() GetValue
}

type Parameters interface {
	GetIndex() int
	GetType() Types
	SetValue(ctx Context, value Value) Control
	GetName() string
	GetValue(ctx Context) (GetValue, Control)
	GetDefaultValue() GetValue
	GetVariables() []Variable
}

// NewParameter 单个参数接收
func NewParameter(name string, index int) Parameter {
	return &ParameterTODO{Name: name, Index: index}
}

func NewParameterDefault(name string, index int, defaultValue GetValue, ty Types) Parameter {
	return &ParameterTODO{Name: name, Index: index, DefaultValue: defaultValue, Type: ty}
}

// NewParameters 支持多个参数接收
func NewParameters(name string, index int) Parameter {
	return &ParametersTODO{Name: name, Index: index}
}

type ParameterTODO struct {
	Name         string // 变量名
	Index        int    // 变量在作用域中的索引
	Type         Types
	DefaultValue GetValue // 默认值
}

func (p *ParameterTODO) GetName() string {
	return p.Name
}

func (p *ParameterTODO) GetValue(ctx Context) (GetValue, Control) {
	val, acl := ctx.GetVariableValue(p)
	if acl != nil {
		return nil, acl
	}

	if _, ok := val.(AsNull); ok {
		if p.DefaultValue != nil {
			val, acl := p.DefaultValue.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}

			p.SetValue(ctx, val.(Value))
		}
	}

	return val, nil
}

func (p *ParameterTODO) GetDefaultValue() GetValue {
	return p.DefaultValue
}

func (p *ParameterTODO) GetIndex() int {
	return p.Index
}

func (p *ParameterTODO) GetType() Types {
	return p.Type
}

func (p *ParameterTODO) SetValue(ctx Context, value Value) Control {
	if p.Type == nil {
		return ctx.SetVariableValue(p, value)
	}
	if p.Type.Is(value) {
		return ctx.SetVariableValue(p, value)
	}
	return NewErrorThrow(nil, errors.New("变量类型和赋值类型不一致, 变量类型("+p.Type.String()+"), 赋值("+value.AsString()+")"))
}

type ParametersTODO struct {
	Name         string // 变量名
	Index        int    // 变量在作用域中的索引
	Type         Types
	DefaultValue GetValue // 默认值
}

func (p *ParametersTODO) GetName() string {
	return p.Name
}

func (p *ParametersTODO) GetValue(ctx Context) (GetValue, Control) {
	v, acl := ctx.GetVariableValue(p)
	if acl != nil {
		return nil, acl
	}

	if _, ok := v.(*ArrayValue); !ok {
		nv := NewArrayValue([]Value{v})
		ctx.SetVariableValue(p, nv)
		return nv, nil
	}

	return v, nil
}

func (p *ParametersTODO) GetDefaultValue() GetValue {
	return p.DefaultValue
}

func (p *ParametersTODO) GetIndex() int {
	return p.Index
}

func (p *ParametersTODO) GetType() Types {
	return p.Type
}

func (p *ParametersTODO) SetValue(ctx Context, value Value) Control {
	if p.Type == nil {
		return ctx.SetVariableValue(p, value)
	}
	if p.Type.Is(value) {
		return ctx.SetVariableValue(p, value)
	}
	return NewErrorThrow(nil, errors.New("变量类型和赋值类型不一致, 变量类型("+p.Type.String()+"), 赋值("+value.AsString()+")"))
}

func (p *ParametersTODO) GetVariables() []Variable {
	return nil
}
