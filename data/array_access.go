package data

import "iter"

// NewArrayValueFromSlots takes ownership of an ordered slot buffer. Slots must
// already carry their PHP keys; use SetKey or AppendValue when building values.
func NewArrayValueFromSlots(slots []*ZVal) *ArrayValue {
	return &ArrayValue{flatArrayStore: FlatArrayStore{entries: slots}}
}

// AppendEntries appends already keyed entries without publishing the backing
// buffer. Callers transfer ownership of their entry metadata.
func (a *ArrayValue) AppendEntries(slots ...*ZVal) {
	a.ensureIndex()
	for _, slot := range slots {
		if slot == nil {
			a.entries = append(a.entries, nil)
			a.invalidateIndex()
			a.ensureIndex()
			continue
		}
		a.appendSlotIncremental(slot)
	}
}

// ReplaceSlot replaces one insertion position, including its key metadata.
func (a *ArrayValue) ReplaceSlot(position int, slot *ZVal) {
	a.entries[position] = slot
	a.invalidateIndex()
}

// SwapSlots is for dense internal containers such as SPL heaps. PHP key
// preserving sorts must use EditPreservingKeys instead.
func (a *ArrayValue) SwapSlots(left, right int) {
	a.entries[left], a.entries[right] = a.entries[right], a.entries[left]
	a.invalidateIndex()
}

// EditPreservingKeys makes implicit integer keys explicit before reordering.
// It preserves the automatic append counter and resets the array pointer.
func (a *ArrayValue) EditPreservingKeys(edit func([]*ZVal)) {
	a.ensureIndex()
	a.normalizeDenseIntKeys()
	a.EditSlots(edit)
	a.iterator = 0
}

// EditReindexing drops all keys, as required by sort/rsort/usort.
func (a *ArrayValue) EditReindexing(edit func([]*ZVal)) {
	a.EditSlots(func(slots []*ZVal) {
		edit(slots)
		for _, slot := range slots {
			if slot != nil {
				slot.Name = ""
				slot.EmptyStrKey = false
			}
		}
	})
	a.appendKeyKnown = false
	a.intKeySeen = false
	a.iterator = 0
}

// Len returns the number of entries, including entries with sparse PHP keys.
func (a *ArrayValue) Len() int { return len(a.entries) }

// View borrows entries for read-only traversal without allocating a snapshot.
func (a *ArrayValue) View() Span[*ZVal] { return NewSpan(a.entries) }

// RemovePositions removes a dense internal container range without copying
// the buffer. PHP unset/shift/pop have separate key and append-counter rules.
func (a *ArrayValue) RemovePositions(start, end int) {
	if start < 0 || end < start || end > a.Len() {
		panic("invalid array position range")
	}
	n := copy(a.entries[start:], a.entries[end:])
	clear(a.entries[start+n:])
	a.ReplaceAll(a.entries[:start+n])
	a.iterator = 0
}

// PrependDense adds a dense value and reindexes the ordered buffer once.
func (a *ArrayValue) PrependDense(value Value) {
	a.entries = append(a.entries, nil)
	copy(a.entries[1:], a.entries[:len(a.entries)-1])
	a.entries[0] = NewZVal(value)
	a.ReindexIntKeys(a.entries)
}

// FilterSlots preserves PHP keys and the automatic append state while compacting
// the existing buffer. The callback must not mutate the array structure.
func (a *ArrayValue) FilterSlots(keep func(int, *ZVal) bool) {
	a.EditPreservingKeys(func(slots []*ZVal) {
		write := 0
		for position, slot := range slots {
			if keep(position, slot) {
				slots[write] = slot
				write++
			}
		}
		clear(slots[write:])
		a.entries = slots[:write]
	})
}

// At reads an entry by insertion position, not by its PHP integer key.
func (a *ArrayValue) At(position int) *ZVal {
	if position < 0 || position >= len(a.entries) {
		return nil
	}
	return a.entries[position]
}

// Range visits slots in insertion order. Structural writes during iteration need
// an explicit Snapshot; Range does not allocate a copy or provide synchronization.
func (a *ArrayValue) Range() iter.Seq2[int, *ZVal] {
	return func(yield func(int, *ZVal) bool) {
		for i, slot := range a.entries {
			if !yield(i, slot) {
				return
			}
		}
	}
}

// Snapshot copies the ordered slot list. ZVals remain shared, so references keep
// their identity. Use CloneArrayValue for PHP value-copy semantics.
func (a *ArrayValue) Snapshot() []*ZVal {
	return a.AppendSlotsTo(nil)
}

// AppendSlotsTo copies slot pointers into a caller-owned buffer without exposing
// the array's mutable slice header. Existing destination capacity is reused.
func (a *ArrayValue) AppendSlotsTo(dst []*ZVal) []*ZVal {
	return append(dst, a.entries...)
}

// EditSlots is the bulk-mutation boundary for operations such as sorting.
// The view is valid only inside edit. Future stores can materialize here once.
func (a *ArrayValue) EditSlots(edit func([]*ZVal)) {
	defer a.invalidateIndex()
	edit(a.entries)
}

// ReplaceAll takes ownership of slots after a bulk replacement or reordering.
// Callers must finish editing keys before calling it, to invalidate the key cache.
func (a *ArrayValue) ReplaceAll(slots []*ZVal) {
	a.entries = slots
	a.appendKeyKnown = false
	a.intKeySeen = false
	a.invalidateIndex()
}

// PopSlot removes the last inserted entry. Unlike unset, PHP array_pop moves
// nNextFreeElement back by one when it removes that exact preceding integer key.
func (a *ArrayValue) PopSlot() *ZVal {
	if a.Len() == 0 {
		return nil
	}
	a.ensureIndex()
	position := a.Len() - 1
	slot := a.entries[position]
	if slot.IsPackedIntSlot() {
		if position == a.nextIntKey-1 {
			a.nextIntKey--
		}
	} else if slot != nil && !slot.EmptyStrKey {
		if key, ok := ParseIntArrayKeyName(slot.Name); ok && key == a.nextIntKey-1 {
			a.nextIntKey--
		}
	}
	a.entries[position] = nil
	a.entries = a.entries[:position]
	a.invalidateIndex()
	a.iterator = 0
	return slot
}

// ReindexIntKeys resets integer keys after array_shift/unshift/splice/sort.
// String keys keep their identity and position. Bulk callers own the slot list.
func (a *ArrayValue) ReindexIntKeys(slots []*ZVal) {
	next := 0
	for position, slot := range slots {
		if slot == nil || slot.EmptyStrKey {
			continue
		}
		_, numeric := ParseIntArrayKeyName(slot.Name)
		if slot.Name == "" || numeric {
			if position == next {
				slot.Name = ""
			} else {
				slot.Name = IntArrayKeyName(next)
			}
			next++
		}
	}
	a.ReplaceAll(slots)
	a.iterator = 0
}

func (a *ArrayValue) ShiftSlot() *ZVal {
	if a.Len() == 0 {
		return nil
	}
	slot := a.entries[0]
	a.ReindexIntKeys(a.entries[1:])
	return slot
}

// SetKey uses PHP scalar array-key conversion. Unsupported keys return false.
// It does not perform copy-on-write separation of the owning variable.
func (a *ArrayValue) SetKey(key Value, value Value) bool {
	switch k := key.(type) {
	case *NullValue:
		a.SetStringKey("", value)
	case *StringValue:
		a.SetStringKey(k.AsString(), value)
	case AsInt:
		i, err := k.AsInt()
		if err != nil {
			return false
		}
		a.SetIntKey(i, value)
	default:
		return false
	}
	return true
}

// slots is a temporary internal view for legacy array methods and bulk
// operations. Writes to this view must end in ReplaceAll; it must not escape data.
func (a *ArrayValue) slots() []*ZVal { return a.entries }
