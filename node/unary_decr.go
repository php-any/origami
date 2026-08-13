package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

type UnaryDecr struct {
	*Node `pp:"-"`
	Right data.GetValue
}

func NewUnaryDecr(from data.From, right data.GetValue) *UnaryDecr {
	return &UnaryDecr{
		Node:  NewNode(from),
		Right: right,
	}
}

func (u *UnaryDecr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 先获取右操作数的值
	rv, rCtl := u.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 根据类型进行自减操作
	switch v := rv.(type) {
	case *data.IntValue:
		i, err := v.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(u.from, err)
		}
		newValue := data.NewIntValue(i - 1)

		// 如果是变量，需要更新变量的值
		if variable, ok := u.Right.(data.Variable); ok {
			ctl := variable.SetValue(ctx, newValue)
			if ctl != nil {
				return nil, ctl
			}
		}

		return newValue, nil
	case *data.FloatValue:
		f, err := v.AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(u.from, err)
		}
		newValue := data.NewFloatValue(f - 1.0)

		// 如果是变量，需要更新变量的值
		if variable, ok := u.Right.(data.Variable); ok {
			ctl := variable.SetValue(ctx, newValue)
			if ctl != nil {
				return nil, ctl
			}
		}

		return newValue, nil
	}

	return nil, data.NewErrorThrow(u.from, errors.New("不支持的类型自减"))
}
