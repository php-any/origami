package node

import "github.com/php-any/origami/data"

// YieldStatement 表示yield语句
type YieldStatement struct {
	*Node       `pp:"-"`
	Key         data.GetValue // 可选的键
	Value       data.GetValue // 值
	ByReference bool
}

// NewYieldStatement 创建一个新的yield语句
func NewYieldStatement(from *TokenFrom, key data.GetValue, value data.GetValue) *YieldStatement {
	return &YieldStatement{
		Node:  NewNode(from),
		Key:   key,
		Value: value,
	}
}

// GetValue 获取yield语句的值
func (y *YieldStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var key data.Value
	var value data.Value

	// 解析键（如果有）
	if y.Key != nil {
		k, ctl := y.Key.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if k != nil {
			key = k.(data.Value)
		}
	}

	// 解析值
	if y.Value != nil {
		if y.ByReference {
			slot, ctl := y.referenceSlot(ctx)
			if ctl != nil {
				return nil, ctl
			}
			return nil, data.NewYieldControlWithContext(key, &data.ArraySlotRef{Slot: slot}, ctx)
		}
		v, ctl := y.Value.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if v != nil {
			value = v.(data.Value)
		} else {
			value = data.NewNullValue()
		}
	} else {
		value = data.NewNullValue()
	}

	// 返回 yield 控制流（普通 yield 使用 YieldValueControl）
	// key 为 nil 表示未显式指定，Generator 会自动分配自增整数 key
	return nil, data.NewYieldControlWithContext(key, value, ctx)
}

// Reference yields accept temporary values with a Notice. Do not use the
// stricter argument rule for rvalues, or evaluate a failed argument twice.
func (y *YieldStatement) referenceSlot(ctx data.Context) (*data.ZVal, data.Control) {
	switch y.Value.(type) {
	case data.Variable, *IndexExpression, interface {
		GetZVal(data.Context) (*data.ZVal, data.Control)
	}:
		return argumentReferenceSlot(ctx, y.Value)
	}
	var value data.GetValue
	var ctl data.Control
	if call, ok := y.Value.(referenceCall); ok {
		value, ctl = call.GetReferenceValue(ctx)
	} else {
		value, ctl = y.Value.GetValue(ctx)
	}
	if ctl != nil {
		return nil, ctl
	}
	v, ok := value.(data.Value)
	if !ok {
		v = data.NewNullValue()
	}
	if slot, ctl := data.ReferenceSlot(v); slot != nil || ctl != nil {
		return slot, ctl
	}
	if ctl := data.EmitPHPError(ctx, 8, "Only variable references should be yielded by reference", y.GetFrom()); ctl != nil {
		return nil, ctl
	}
	return data.NewZVal(v), nil
}

// YieldFromStatement 表示yield from语句
type YieldFromStatement struct {
	*Node  `pp:"-"`
	Source data.GetValue // 要委托的生成器或可迭代对象
}

// NewYieldFromStatement 创建一个新的yield from语句
func NewYieldFromStatement(from *TokenFrom, source data.GetValue) *YieldFromStatement {
	return &YieldFromStatement{
		Node:   NewNode(from),
		Source: source,
	}
}

// GetValue 获取yield from语句的值
func (y *YieldFromStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 解析源表达式
	source, ctl := y.Source.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	if source == nil {
		return nil, data.NewErrorThrow(y.GetFrom(), data.NewError(y.GetFrom(), "yield from 需要一个可迭代对象", nil))
	}

	// 返回 yield from 控制流：使用单独的 YieldFromControl，实现委托逻辑
	return nil, NewYieldFromControl(ctx, source)
}
