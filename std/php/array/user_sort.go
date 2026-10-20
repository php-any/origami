package array

import (
	"github.com/php-any/origami/data"
)

// Comparison failures propagate back to PHP. Once a callback fails, no further
// callbacks run. Like PHP, usort renumbers keys even when a comparator throws.
// The algorithm follows zend_sort so inconsistent comparators (Filament's
// one-argument usort) keep the same order as PHP.
func userSort(ctx data.Context, array *data.ArrayValue, callback data.Value, preserve bool) data.Control {
	var failure data.Control
	cmp := func(a, b *data.ZVal) int {
		if failure != nil {
			return 0
		}
		value, control := invokeCallback(ctx, callback, []data.Value{a.ReadValue(), b.ReadValue()})
		if returned, ok := control.(data.ReturnControl); ok {
			value, _ = returned.ReturnValue().(data.Value)
			control = nil
		}
		if control != nil {
			failure = control
			return 0
		}
		return compareSign(value)
	}
	edit := func(slots []*data.ZVal) {
		zendSort(slots, cmp)
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

func compareSign(v data.Value) int {
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	if fv, ok := v.(*data.FloatValue); ok {
		switch {
		case fv.Value > 0:
			return 1
		case fv.Value < 0:
			return -1
		default:
			return 0
		}
	}
	if as, ok := v.(data.AsInt); ok {
		if i, err := as.AsInt(); err == nil {
			return i
		}
	}
	return 0
}

func swapSlot(slots []*data.ZVal, i, j int) {
	slots[i], slots[j] = slots[j], slots[i]
}

func zendSort(slots []*data.ZVal, cmp func(a, b *data.ZVal) int) {
	for {
		n := len(slots)
		if n <= 16 {
			zendInsertSort(slots, cmp)
			return
		}
		offset := n >> 1
		if n>>10 != 0 {
			delta := offset >> 1
			zendSort5(slots, 0, delta, offset, offset+delta, n-1, cmp)
		} else {
			zendSort3(slots, 0, offset, n-1, cmp)
		}
		swapSlot(slots, 1, offset)
		pivot := 1
		i := pivot + 1
		j := n - 1
		for {
			for cmp(slots[pivot], slots[i]) > 0 {
				i++
				if i == j {
					goto done
				}
			}
			j--
			if j == i {
				goto done
			}
			for cmp(slots[j], slots[pivot]) > 0 {
				j--
				if j == i {
					goto done
				}
			}
			swapSlot(slots, i, j)
			i++
			if i == j {
				goto done
			}
		}
	done:
		swapSlot(slots, pivot, i-1)
		left := i - 1
		if left < n-i {
			zendSort(slots[:left], cmp)
			slots = slots[i:]
			continue
		}
		zendSort(slots[i:], cmp)
		slots = slots[:left]
	}
}

func zendSort2(slots []*data.ZVal, a, b int, cmp func(a, b *data.ZVal) int) {
	if cmp(slots[a], slots[b]) > 0 {
		swapSlot(slots, a, b)
	}
}

func zendSort3(slots []*data.ZVal, a, b, c int, cmp func(a, b *data.ZVal) int) {
	if !(cmp(slots[a], slots[b]) > 0) {
		if !(cmp(slots[b], slots[c]) > 0) {
			return
		}
		swapSlot(slots, b, c)
		if cmp(slots[a], slots[b]) > 0 {
			swapSlot(slots, a, b)
		}
		return
	}
	if !(cmp(slots[c], slots[b]) > 0) {
		swapSlot(slots, a, c)
		return
	}
	swapSlot(slots, a, b)
	if cmp(slots[b], slots[c]) > 0 {
		swapSlot(slots, b, c)
	}
}

func zendSort4(slots []*data.ZVal, a, b, c, d int, cmp func(a, b *data.ZVal) int) {
	zendSort3(slots, a, b, c, cmp)
	if cmp(slots[c], slots[d]) > 0 {
		swapSlot(slots, c, d)
		if cmp(slots[b], slots[c]) > 0 {
			swapSlot(slots, b, c)
			if cmp(slots[a], slots[b]) > 0 {
				swapSlot(slots, a, b)
			}
		}
	}
}

func zendSort5(slots []*data.ZVal, a, b, c, d, e int, cmp func(a, b *data.ZVal) int) {
	zendSort4(slots, a, b, c, d, cmp)
	if cmp(slots[d], slots[e]) > 0 {
		swapSlot(slots, d, e)
		if cmp(slots[c], slots[d]) > 0 {
			swapSlot(slots, c, d)
			if cmp(slots[b], slots[c]) > 0 {
				swapSlot(slots, b, c)
				if cmp(slots[a], slots[b]) > 0 {
					swapSlot(slots, a, b)
				}
			}
		}
	}
}

func zendInsertSort(slots []*data.ZVal, cmp func(a, b *data.ZVal) int) {
	n := len(slots)
	switch n {
	case 0, 1:
		return
	case 2:
		zendSort2(slots, 0, 1, cmp)
		return
	case 3:
		zendSort3(slots, 0, 1, 2, cmp)
		return
	case 4:
		zendSort4(slots, 0, 1, 2, 3, cmp)
		return
	case 5:
		zendSort5(slots, 0, 1, 2, 3, 4, cmp)
		return
	}
	limit := 6
	if limit > n {
		limit = n
	}
	for i := 1; i < limit; i++ {
		j := i - 1
		if !(cmp(slots[j], slots[i]) > 0) {
			continue
		}
		for j != 0 {
			j--
			if !(cmp(slots[j], slots[i]) > 0) {
				j++
				break
			}
		}
		for k := i; k > j; k-- {
			swapSlot(slots, k, k-1)
		}
	}
	for i := 6; i < n; i++ {
		j := i - 1
		if !(cmp(slots[j], slots[i]) > 0) {
			continue
		}
		for {
			j -= 2
			if !(cmp(slots[j], slots[i]) > 0) {
				j++
				if !(cmp(slots[j], slots[i]) > 0) {
					j++
				}
				break
			}
			if j == 0 {
				break
			}
			if j == 1 {
				j--
				if cmp(slots[i], slots[j]) > 0 {
					j++
				}
				break
			}
		}
		for k := i; k > j; k-- {
			swapSlot(slots, k, k-1)
		}
	}
}
