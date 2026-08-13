package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

type PostfixDecr struct {
	*Node `pp:"-"`
	Left  data.GetValue
}

func NewPostfixDecr(from data.From, left data.GetValue) data.GetValue {
	if ve, ok := left.(*VariableExpression); ok {
		fb := &PostfixDecr{Node: NewNode(from), Left: left}
		return &VarPostDecr{Node: NewNode(from), VarIdx: ve.Index, Var: ve, Fallback: fb}
	}
	return &PostfixDecr{
		Node: NewNode(from),
		Left: left,
	}
}

func (p *PostfixDecr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if cop, ok := p.Left.(*CallObjectProperty); ok {
		return p.decrObjectProperty(ctx, cop)
	}

	if ie, ok := p.Left.(*IndexExpression); ok {
		if className, overloaded := arrayAccessOverloadedClass(ctx, ie); overloaded {
			lv, ctl := ie.GetValue(ctx)
			emitIndirectModificationNotice(p.GetFrom(), className)
			return lv, ctl
		}
	}

	// 先获取左操作数的值
	lv, lCtl := p.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	var originalValue data.GetValue

	switch v := lv.(type) {
	case *data.IntValue:
		i, err := v.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(p.from, err)
		}
		originalValue = data.NewIntValue(i)
		newValue := data.NewIntValue(i - 1)
		if variable, ok := p.Left.(data.Variable); ok {
			if ctl := variable.SetValue(ctx, newValue); ctl != nil {
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
		newValue := data.NewFloatValue(f - 1.0)
		if variable, ok := p.Left.(data.Variable); ok {
			if ctl := variable.SetValue(ctx, newValue); ctl != nil {
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
		newValue := data.NewIntValue(i - 1)
		if variable, ok := p.Left.(data.Variable); ok {
			if ctl := variable.SetValue(ctx, newValue); ctl != nil {
				return nil, ctl
			}
		}
		return originalValue, nil
	}

	if asInt, ok := lv.(data.AsInt); ok {
		i, err := asInt.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(p.from, err)
		}
		originalValue = data.NewIntValue(i)
		newValue := data.NewIntValue(i - 1)
		if variable, ok := p.Left.(data.Variable); ok {
			if ctl := variable.SetValue(ctx, newValue); ctl != nil {
				return nil, ctl
			}
		} else if cop, ok := p.Left.(*CallObjectProperty); ok {
			if ctl := cop.SetValue(ctx, newValue); ctl != nil {
				return nil, ctl
			}
		}
		return originalValue, nil
	}

	return nil, data.NewErrorThrow(p.from, errors.New("不支持的类型自减"))
}

func (p *PostfixDecr) decrObjectProperty(ctx data.Context, cop *CallObjectProperty) (data.GetValue, data.Control) {
	lv, lCtl := cop.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}
	i := 0
	switch v := lv.(type) {
	case *data.NullValue:
		i = 0
	case *data.IntValue:
		var err error
		i, err = v.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(p.from, err)
		}
	case *data.FloatValue:
		f, err := v.AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(p.from, err)
		}
		i = int(f)
	default:
		if asInt, ok := lv.(data.AsInt); ok {
			var err error
			i, err = asInt.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(p.from, err)
			}
		}
	}
	originalValue := data.NewIntValue(i)
	if _, ok := lv.(*data.NullValue); ok {
		originalValue = data.NewNullValue()
	}
	newValue := data.NewIntValue(i - 1)
	if ctl := cop.SetValue(ctx, newValue); ctl != nil {
		return nil, ctl
	}
	return originalValue, nil
}
