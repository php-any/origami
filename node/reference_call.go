package node

import "github.com/php-any/origami/data"

type referenceCall interface {
	GetReferenceValue(data.Context) (data.GetValue, data.Control)
}

func callValue(value data.GetValue, ctl data.Control) (data.GetValue, data.Control) {
	if ctl == nil {
		if ref, ok := value.(*data.ArraySlotRef); ok && ref.Slot != nil {
			return ref.Slot.ReadValue(), nil
		}
		if ref, ok := value.(*data.ReferenceValue); ok {
			return ref.Val.GetValue(ref.Ctx)
		}
	}
	return value, ctl
}
