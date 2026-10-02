package array

import (
	"sort"

	"github.com/php-any/origami/data"
)

// Comparison failures propagate back to PHP. Once a callback fails, no further
// callbacks run. Like PHP, usort renumbers keys even when a comparator throws.
func userSort(ctx data.Context, array *data.ArrayValue, callback data.Value, preserve bool) data.Control {
	var failure data.Control
	edit := func(slots []*data.ZVal) {
		sort.SliceStable(slots, func(i, j int) bool {
			if failure != nil {
				return false
			}
			value, control := invokeCallback(ctx, callback, []data.Value{slots[i].Value, slots[j].Value})
			if returned, ok := control.(data.ReturnControl); ok {
				value, _ = returned.ReturnValue().(data.Value)
				control = nil
			}
			if control != nil {
				failure = control
				return false
			}
			return compareLess(value)
		})
		if !preserve {
			for _, slot := range slots {
				slot.Name = ""
				slot.EmptyStrKey = false
			}
			array.ReplaceAll(slots)
		}
	}
	if preserve {
		array.EditPreservingKeys(edit)
	} else {
		array.EditSlots(edit)
	}
	return failure
}
