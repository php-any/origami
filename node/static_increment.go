package node

import "github.com/php-any/origami/data"

func incrementStaticProperty(ctx data.Context, expression data.GetValue, delta int, postfix bool) (data.GetValue, data.Control, bool) {
	switch expression.(type) {
	case *CallSelfProperty, *CallStaticProperty, *CallStaticPropertyLater, *CallStaticKeywordProperty:
	default:
		return nil, nil, false
	}
	slot, ctl := expression.(interface {
		GetZVal(data.Context) (*data.ZVal, data.Control)
	}).GetZVal(ctx)
	if ctl != nil {
		return nil, ctl, true
	}
	old := slot.ReadValue()
	var next data.Value
	switch value := old.(type) {
	case *data.IntValue:
		next = data.NewIntValue(value.Value + delta)
	case *data.FloatValue:
		next = data.NewFloatValue(value.Value + float64(delta))
	case *data.NullValue:
		if delta < 0 {
			next = old
		} else {
			next = data.NewIntValue(1)
		}
	default:
		return nil, nil, false
	}
	next, ctl = slot.PrepareWrite(next, ctx)
	if ctl != nil {
		return nil, ctl, true
	}
	data.CowAssign(slot, next)
	if postfix {
		return old, nil, true
	}
	return next, nil, true
}
