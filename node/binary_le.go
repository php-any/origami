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

	switch l := lv.(type) {
	case *data.IntValue:
		if ri, ok := rv.(data.AsInt); ok {
			li, err := l.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			riVal, err := ri.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewBoolValue(li <= riVal), nil
		}
	case *data.FloatValue:
		if rf, ok := rv.(data.AsFloat); ok {
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rfVal, err := rf.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewBoolValue(lf <= rfVal), nil
		}
	case *data.StringValue:
		if rs, ok := rv.(data.AsString); ok {
			ls := l.AsString()
			rsVal := rs.AsString()
			return data.NewBoolValue(ls <= rsVal), nil
		}
	}

	return data.NewBoolValue(false), nil
}
