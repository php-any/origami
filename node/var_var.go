package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// VarVar 表示 PHP 的变量变量：$$var
// 解析时会捕获当前作用域内的变量列表，运行时根据名称在这些变量中查找。
type VarVar struct {
	*Node `pp:"-"`

	// NameExpr 用于计算变量名的表达式，一般是一个普通变量（如 $field）
	NameExpr data.GetValue

	// Vars 是解析时快照的当前作用域变量列表
	Vars []data.Variable
}

// NewVarVar 创建一个变量变量节点
func NewVarVar(from data.From, nameExpr data.GetValue, vars []data.Variable) *VarVar {
	return &VarVar{
		Node:     NewNode(from),
		NameExpr: nameExpr,
		Vars:     vars,
	}
}

// GetValue 实现 data.GetValue：按当前上下文计算变量名，再在捕获的变量列表中查找对应变量的值。
func (v *VarVar) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 先计算名称表达式的值
	nameVal, acl := v.NameExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	val, ok := nameVal.(data.Value)
	if !ok {
		return nil, data.NewErrorThrow(v.from, errors.New("变量变量名称必须是可转换为字符串的值"))
	}
	name := val.AsString()
	if name == "" {
		// 名称为空时，返回 null，避免产生不可预期的变量名
		return data.NewNullValue(), nil
	}

	// 在捕获的变量列表中查找同名变量
	for _, vari := range v.Vars {
		if vari == nil {
			continue
		}
		if vari.GetName() == name {
			got, ctl := ctx.GetVariableValue(vari)
			if ctl != nil {
				return nil, ctl
			}
			return got, nil
		}
	}

	// 与 PHP 行为对齐：未找到变量时返回 null
	return data.NewNullValue(), nil
}
