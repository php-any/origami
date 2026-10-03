package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

func globalsDimensionSlot(ctx data.Context, index *IndexExpression) (*data.ZVal, data.Control) {
	if index.Append {
		return nil, data.NewErrorThrowByName(index.GetFrom(), fmt.Errorf("Cannot append to $GLOBALS"), "Error")
	}
	key, ctl := index.Index.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name, ok := indexKeyString(key)
	if !ok {
		return nil, data.NewTypeError(index.GetFrom(), fmt.Errorf("Illegal offset type"))
	}
	return ctx.GetVM().EnsureGlobalZVal(name), nil
}

func assignGlobalDimension(ctx data.Context, slot *data.ZVal, value data.Value) data.Control {
	if source, ctl := data.ReferenceSlot(value); source != nil || ctl != nil {
		if ctl != nil {
			return ctl
		}
		if source == slot {
			return nil
		}
		name := slot.Name
		slot.ReleaseRefSlot()
		source.AddRefSlot()
		*slot = *data.CopyReferenceBucket(source)
		slot.Name, slot.Defined = name, true
		return nil
	}
	prepared, ctl := slot.PrepareWrite(value, ctx)
	if ctl != nil {
		return ctl
	}
	data.CowAssign(slot, prepared)
	return nil
}
