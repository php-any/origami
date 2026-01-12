package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// FunctionStatement 表示函数定义语句
type FunctionStatement struct {
	data.FuncStmt
	*Node  `pp:"-"`
	Name   string          // 函数名
	Params []data.GetValue // 参数列表
	Body   []data.GetValue // 函数体
	vars   []data.Variable // 符号表
	Ret    data.Types      // 返回值类型
}

// NewFunctionStatement 创建一个新的函数定义语句
func NewFunctionStatement(from data.From, name string, params []data.GetValue, body []data.GetValue, vars []data.Variable, ret data.Types) *FunctionStatement {
	return &FunctionStatement{
		Node:   NewNode(from),
		Name:   name,
		Params: params,
		Body:   body,
		vars:   vars,
		Ret:    ret,
	}
}

// GetName 返回函数名
func (f *FunctionStatement) GetName() string {
	return f.Name
}

// GetBody 返回函数体
func (f *FunctionStatement) GetBody() []data.GetValue {
	return f.Body
}

// GetValue 获取函数定义语句的值
func (f *FunctionStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if acl := ctx.GetVM().AddFunc(f); acl != nil {
		return nil, acl
	}
	return nil, nil
}

func (f *FunctionStatement) GetParams() []data.GetValue {
	return f.Params
}

func (f *FunctionStatement) GetVariables() []data.Variable {
	return f.vars
}

// GetReturnType 返回函数返回类型
func (f *FunctionStatement) GetReturnType() data.Types {
	return f.Ret
}

func (f *FunctionStatement) Call(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var ctl data.Control
	for bodyIndex, statement := range f.Body {
		v, ctl = statement.GetValue(ctx)
		if ctl != nil {
			switch rv := ctl.(type) {
			case data.ReturnControl:
				ret := rv.ReturnValue()
				if f.Ret != nil {
					if f.Ret.Is(ret) {
						return ret, nil
					} else {
						return nil, data.NewErrorThrow(f.from, errors.New("函数返回值类型错误"))
					}
				}
				return ret, nil
			case data.YieldControl:
				generator := rv.CreateStackState(ctx, f, f.Body, bodyIndex)
				// 将生成器包装成类值，支持 $data->valid() 等调用
				generatorClass := NewGeneratorClass(generator)
				return generatorClass.GetValue(ctx)
			case data.YieldValueControl:
				generator := NewFuncYieldStackState(ctx, f, f.Body, bodyIndex+1, rv.GetYieldKey(), rv.GetYieldValue())
				// 将生成器包装成类值，支持 $data->valid() 等调用
				generatorClass := NewGeneratorClass(generator)
				return generatorClass.GetValue(ctx)
			case data.AddStack:
				if from, ok := statement.(GetFrom); ok {
					rv.AddStackWithInfo(from.GetFrom(), "function body", TryGetCallClassName(statement))
				}
				rv.AddStackWithInfo(f.from, "function", f.Name)
			}
			return nil, ctl
		}
	}

	return v, nil
}

// Parameter 表示函数参数
type Parameter struct {
	*Node        `pp:"-"`
	Name         string // 变量名
	Index        int    // 变量在作用域中的索引
	Type         data.Types
	DefaultValue data.GetValue // 默认值
}

func (p *Parameter) GetDefaultValue() data.GetValue {
	return p.DefaultValue
}

func (p *Parameter) GetIndex() int {
	return p.Index
}

func (p *Parameter) GetType() data.Types {
	return p.Type
}

func (p *Parameter) SetValue(ctx data.Context, value data.Value) data.Control {
	if p.Type == nil {
		return ctx.SetVariableValue(p, value)
	}
	if p.Type.Is(value) {
		return ctx.SetVariableValue(p, value)
	}
	if p.Type.Is(value) {
		return ctx.SetVariableValue(p, value)
	}
	if _, ok := value.(*data.FuncValue); ok {
		if p.Name == "closure" {
			return ctx.SetVariableValue(p, value)
		}
	}
	return data.NewErrorThrow(p.from, errors.New("变量类型和赋值类型不一致, 变量类型("+p.Type.String()+"), 赋值("+value.AsString()+")"))
}

// NewParameter 创建一个新的参数
func NewParameter(from data.From, name string, index int, defaultValue data.GetValue, ty data.Types) data.GetValue {
	return &Parameter{
		Node:         NewNode(from),
		Name:         name,
		Index:        index,
		Type:         ty,
		DefaultValue: defaultValue,
	}
}

// GetName 返回参数名
func (p *Parameter) GetName() string {
	return p.Name
}

// PromotedParameter 表示属性提升的参数（构造函数参数属性提升）
type PromotedParameter struct {
	*Parameter
	PropertyName string // 对应的属性名（与参数名相同）
}

// NewPromotedParameter 创建一个新的属性提升参数
func NewPromotedParameter(from data.From, name string, index int, defaultValue data.GetValue, ty data.Types) data.GetValue {
	return &PromotedParameter{
		Parameter: &Parameter{
			Node:         NewNode(from),
			Name:         name,
			Index:        index,
			Type:         ty,
			DefaultValue: defaultValue,
		},
		PropertyName: name,
	}
}

func (p *PromotedParameter) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return p.Parameter.GetValue(ctx)
}

func (p *Parameter) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	val, acl := ctx.GetVariableValue(p)
	if acl != nil {
		return nil, acl
	}

	if _, ok := val.(data.AsNull); ok {
		if p.DefaultValue != nil {
			val, acl := p.DefaultValue.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}

			p.SetValue(ctx, val.(data.Value))
		}
	}

	return val, nil
}

// NewParameters 接收多个参数值
func NewParameters(from data.From, name string, index int, defaultValue data.GetValue, ty data.Types) data.GetValue {
	return &Parameters{
		Parameter: &Parameter{
			Node:         NewNode(from),
			Name:         name,
			Index:        index,
			Type:         ty,
			DefaultValue: defaultValue,
		},
	}
}

// Parameters 多值参数
type Parameters struct {
	*Parameter
}

func (p *Parameters) SetValue(ctx data.Context, value data.Value) data.Control {
	//TODO implement me
	panic("implement me")
}

func (p *Parameters) GetDefaultValue() data.GetValue {
	return p.DefaultValue
}

func (p *Parameters) GetName() string {
	return p.Name
}

func (p *Parameters) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	v, acl := ctx.GetVariableValue(p)
	if acl != nil {
		return nil, acl
	}

	if _, ok := v.(*data.ArrayValue); !ok {
		nv := data.NewArrayValue([]data.Value{v})
		ctx.SetVariableValue(p, nv)
		return nv, nil
	}

	return v, nil
}

func (p *Parameters) GetVariables() []data.Variable {
	return nil
}

type ParameterReference struct {
	*Parameter
}

func NewParameterReference(from data.From, name string, index int, ty data.Types) data.GetValue {
	return &ParameterReference{
		Parameter: &Parameter{
			Node:  NewNode(from),
			Name:  name,
			Index: index,
			Type:  ty,
		},
	}
}

func (p *ParameterReference) SetValue(ctx data.Context, value data.Value) data.Control {
	if p.Type != nil {
		if !p.Type.Is(value) {
			return data.NewErrorThrow(p.from, errors.New("变量类型和赋值类型不一致, 变量类型("+p.Type.String()+"), 赋值("+value.AsString()+")"))
		}
	}
	if v, ok := value.(*data.ZValValue); ok {
		ctx.SetIndexZVal(p.Index, v.ZVal)
	} else {
		return ctx.SetVariableValue(p, value)
	}

	return nil
}

// NewParametersReference 接收多个参数值
func NewParametersReference(from data.From, name string, index int, defaultValue data.GetValue, ty data.Types) data.GetValue {
	return &ParametersReference{
		Parameter: &Parameter{
			Node:         NewNode(from),
			Name:         name,
			Index:        index,
			Type:         ty,
			DefaultValue: defaultValue,
		},
	}
}

// ParametersReference 多值参数
type ParametersReference struct {
	*Parameter
}

func (p *ParametersReference) GetDefaultValue() data.GetValue {
	return p.DefaultValue
}

func (p *ParametersReference) GetName() string {
	return p.Name
}

func (p *ParametersReference) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	v, acl := ctx.GetVariableValue(p)
	if acl != nil {
		return nil, acl
	}

	if _, ok := v.(*data.ArrayValue); !ok {
		nv := data.NewArrayValue([]data.Value{v})
		ctx.SetVariableValue(p, nv)
		return nv, nil
	}

	return v, nil
}

func (p *ParametersReference) SetValue(ctx data.Context, value data.Value) data.Control {
	if p.Type == nil {
		return ctx.SetVariableValue(p, value)
	}
	if p.Type.Is(value) {
		return ctx.SetVariableValue(p, value)
	}
	return data.NewErrorThrow(p.from, errors.New("变量类型和赋值类型不一致, 变量类型("+p.Type.String()+"), 赋值("+value.AsString()+")"))
}

// CallerContextParameter 特殊参数类型：用于标记函数需要在调用者的 Context 中执行。
// 主要用于实现类似 func_get_args 这类需要直接访问上级调用入参的函数。
type CallerContextParameter struct {
	*Node `pp:"-"`
}

// NewCallerContextParameter 创建一个新的 CallerContextParameter。
// from 仅用于错误栈信息，可以为 nil。
func NewCallerContextParameter(from data.From) data.GetValue {
	return &CallerContextParameter{
		Node: NewNode(from),
	}
}

// GetValue 对函数调用参数绑定过程不产生实际值，这里直接返回自身即可。
func (p *CallerContextParameter) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return p, nil
}

type CallFunctionLater struct {
	Ctx  data.Context
	Name string
	Fun  data.FuncStmt
}

func (c *CallFunctionLater) Call(ctx data.Context) (data.GetValue, data.Control) {
	if c.Fun == nil {
		c.Fun, _ = c.Ctx.GetVM().GetFunc(c.Name)
	}
	return c.Fun.Call(ctx)
}

func (c *CallFunctionLater) GetName() string {
	if c.Fun == nil {
		c.Fun, _ = c.Ctx.GetVM().GetFunc(c.Name)
	}
	return c.Fun.GetName()
}

func (c *CallFunctionLater) GetParams() []data.GetValue {
	if c.Fun == nil {
		c.Fun, _ = c.Ctx.GetVM().GetFunc(c.Name)
	}
	return c.Fun.GetParams()
}

func (c *CallFunctionLater) GetVariables() []data.Variable {
	if c.Fun == nil {
		c.Fun, _ = c.Ctx.GetVM().GetFunc(c.Name)
	}
	return c.Fun.GetVariables()
}
