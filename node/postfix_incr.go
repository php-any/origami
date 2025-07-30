package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// PostfixIncr 前自增
type PostfixIncr struct {
	*Node `pp:"-"`
	Left  data.GetValue
}

func NewPostfixIncr(from data.From, left data.GetValue) *PostfixIncr {
	return &PostfixIncr{
		Node: NewNode(from),
		Left: left,
	}
}

func (p *PostfixIncr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 先获取左操作数的值
	lv, lCtl := p.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	// 保存原始值用于返回
	var originalValue data.GetValue

	// 根据类型进行自增操作
	switch v := lv.(type) {
	case *data.IntValue:
		i, err := v.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(p.from, err)
		}
		originalValue = data.NewIntValue(i)
		newValue := data.NewIntValue(i + 1)

		// 如果是变量，需要更新变量的值
		if variable, ok := p.Left.(data.Variable); ok {
			ctl := variable.SetValue(ctx, newValue)
			if ctl != nil {
				return nil, ctl
			}
		}

		return originalValue, nil
	case *data.FloatValue:
		f, err := v.AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(p.from, err)
		}
		originalValue = data.NewFloatValue(f)
		newValue := data.NewFloatValue(f + 1.0)

		// 如果是变量，需要更新变量的值
		if variable, ok := p.Left.(data.Variable); ok {
			ctl := variable.SetValue(ctx, newValue)
			if ctl != nil {
				return nil, ctl
			}
		}

		return originalValue, nil
	case *data.NullValue:
		i, err := v.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(p.from, err)
		}
		originalValue = data.NewIntValue(i)
		newValue := data.NewIntValue(i + 1)

		// 如果是变量，需要更新变量的值
		if variable, ok := p.Left.(data.Variable); ok {
			ctl := variable.SetValue(ctx, newValue)
			if ctl != nil {
				return nil, ctl
			}
		}

		return originalValue, nil
	}

	return nil, data.NewErrorThrow(p.from, errors.New("不支持的类型自增"))
}
