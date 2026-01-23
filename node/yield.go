package node

import "github.com/php-any/origami/data"

// YieldStatement 表示yield语句
type YieldStatement struct {
	*Node `pp:"-"`
	Key   data.GetValue // 可选的键
	Value data.GetValue // 值
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

	// 如果没有指定 key，使用 null 作为默认值（但 GetKey() 应该始终返回非 nil）
	if key == nil {
		key = data.NewNullValue()
	}
	// 返回 yield 控制流（普通 yield 使用 YieldValueControl）
	return nil, data.NewYieldControlWithContext(key, value, ctx)
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
