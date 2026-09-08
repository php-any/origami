package node

import (
	"github.com/php-any/origami/data"
)

type BinaryLe struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryLe(from data.From, left, right data.GetValue) data.GetValue {
	le := &BinaryLe{Node: NewNode(from), Left: left, Right: right}
	// 解析阶段模式识别：$var <= IntLiteral 是 for 循环条件最常见形式，
	// 发出 VarIntLe 节点，实现 BoolTest 接口，允许 ForStatement 绕过
	// BoolValue 分配直接获取 bool 值。
	if ve, ok := left.(*VariableExpression); ok {
		if lit, ok := right.(*IntLiteral); ok {
			if iv, ok := lit.V.(*data.IntValue); ok {
				return &VarIntLe{Node: NewNode(from), VarIdx: ve.Index, Lit: iv.Value, Le: le}
			}
		}
	}
	return le
}

func (b *BinaryLe) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	if cmp, ok := phpCompareValues(ctx, lv, rv); ok {
		return data.NewBoolValue(cmp <= 0), nil
	}
	return data.NewBoolValue(false), nil
}
