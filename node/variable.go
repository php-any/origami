package node

import (
	"errors"
	"github.com/php-any/origami/data"
)

// VariableExpression 表示变量表达式
type VariableExpression struct {
	*Node `pp:"-"`
	Name  string // 变量名
	Index int    // 变量在作用域中的索引
	Type  data.Types
}

// NewVariableWithFirst 解释器创建变量前, 需要先识别定义时的信息 p.scopeManager.LookupVariable(name)
func NewVariableWithFirst(from data.From, first data.Variable) data.Variable {
	if _, ok := first.(*VariableReference); ok {
		return NewVariableReference(from, first.GetName(), first.GetIndex(), first.GetType())
	}
	return &VariableExpression{
		Node:  NewNode(from),
		Name:  first.GetName(),
		Index: first.GetIndex(),
		Type:  first.GetType(),
	}
}
func NewVariable(from data.From, name string, index int, ty data.Types) *VariableExpression {
	if name[0:1] == "$" {
		name = name[1:]
	}
	return &VariableExpression{
		Node:  NewNode(from),
		Name:  name,
		Index: index,
		Type:  ty,
	}
}

// GetValue 获取变量表达式的值
func (v *VariableExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return ctx.GetVariableValue(v)
}

func (v *VariableExpression) GetIndex() int {
	return v.Index
}
func (v *VariableExpression) GetName() string {
	return v.Name
}
func (v *VariableExpression) GetType() data.Types {
	return v.Type
}

func (v *VariableExpression) SetValue(ctx data.Context, value data.Value) data.Control {
	if v.Type == nil {
		return ctx.SetVariableValue(v, value)
	}
	if v.Type.Is(value) {
		return ctx.SetVariableValue(v, value)
	}
	return data.NewErrorThrow(v.from, errors.New("变量类型和赋值类型不一致, 变量类型("+v.Type.String()+"), 赋值("+value.AsString()+")"))
}

// VariableList 支持多变量解包赋值

type VariableList struct {
	Vars []*VariableExpression
}

func NewVariableList(vars []*VariableExpression) *VariableList {
	return &VariableList{Vars: vars}
}

func (vl *VariableList) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 返回所有变量的值组成的数组
	var values []data.Value
	for _, v := range vl.Vars {
		val, ctl := v.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if vv, ok := val.(data.Value); ok {
			values = append(values, vv)
		} else {
			values = append(values, data.NewNullValue())
		}
	}
	return data.NewArrayValue(values), nil
}

func (vl *VariableList) GetIndex() int {
	if len(vl.Vars) > 0 {
		return vl.Vars[0].GetIndex()
	}
	return 0
}

func (vl *VariableList) GetName() string {
	// 用逗号拼接所有变量名
	names := ""
	for i, v := range vl.Vars {
		if i > 0 {
			names += ","
		}
		names += v.GetName()
	}
	return names
}

func (vl *VariableList) GetType() data.Types {
	// 多变量类型一般不做类型约束，返回 nil
	return nil
}

func (vl *VariableList) SetValue(ctx data.Context, value data.Value) data.Control {
	// value 应为 ArrayValue，依次赋值
	arr, ok := value.(*data.ArrayValue)
	if !ok {
		// 单值赋给第一个变量
		if len(vl.Vars) > 0 {
			return vl.Vars[0].SetValue(ctx, value)
		}
		return nil
	}
	for i, v := range vl.Vars {
		var val data.Value = data.NewNullValue()
		if i < len(arr.Value) {
			val = arr.Value[i]
		}
		ctl := v.SetValue(ctx, val)
		if ctl != nil {
			return ctl
		}
	}
	return nil
}

type VariableReference struct {
	*Node `pp:"-"`
	Name  string // 变量名
	Index int    // 变量在作用域中的索引
	Type  data.Types

	ctx data.Context
}

// NewVariableReference 创建一个新的变量引用
func NewVariableReference(from data.From, name string, index int, ty data.Types) *VariableReference {
	if name[0:1] == "$" {
		name = name[1:]
	}
	return &VariableReference{
		Node:  NewNode(from),
		Name:  name,
		Index: index,
		Type:  ty,
	}
}

// GetValue 获取变量表达式的值
func (v *VariableReference) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return v.ctx.GetVariableValue(v)
}

func (v *VariableReference) GetIndex() int {
	return v.Index
}
func (v *VariableReference) GetName() string {
	return v.Name
}
func (v *VariableReference) GetType() data.Types {
	return v.Type
}

func (v *VariableReference) SetValue(ctx data.Context, value data.Value) data.Control {
	if v.Type != nil {
		if !v.Type.Is(value) {
			return data.NewErrorThrow(v.from, errors.New("变量类型和赋值类型不一致, 变量类型("+v.Type.String()+"), 赋值("+value.AsString()+")"))
		}
	}
	temp, ok := ctx.GetVariableValue(v)
	if ok != nil {
		return ok
	}
	if ref, ok := temp.(*data.ReferenceValue); ok {
		ctx.SetVariableValue(v, value)
		return ref.Ctx.SetVariableValue(v, value)
	}

	return nil
}
