package node

import (
	"errors"
	"github.com/php-any/origami/data"
	"math"
)

func assignArrayIndex(ctx data.Context, array *data.ArrayValue, index data.GetValue, appendIndex bool, value data.Value, from data.From) data.Control {
	key, ok := index.(data.Value)
	if appendIndex {
		next := array.NextAppendIntKey()
		if next == math.MaxInt {
			if slot, _ := array.FindSlotByIntKey(next); slot != nil {
				return data.NewErrorThrowByName(from, errors.New("Cannot add element to the array as the next element is already occupied"), "Error")
			}
		}
		key, ok = data.NewIntValue(next), true
	}
	if !ok {
		return data.NewTypeError(from, errors.New("Illegal offset type"))
	}
	accepted, ctl := array.AssignKey(ctx, key, value)
	if ctl != nil {
		return ctl
	}
	if !accepted {
		return data.NewTypeError(from, errors.New("Illegal offset type"))
	}
	return nil
}
