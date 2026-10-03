package data

import "strconv"

// ArrayOverlayScope owns the aliases of one request's array graph. Parents
// must be immutable while the scope is alive; only the private index snapshot
// and request slots are mutated. Ordinary PHP value copies still use COW.
type ArrayOverlayScope struct {
	arrays    map[*ArrayValue]*ArrayValue
	refs      map[*ZVal]*ZVal
	cells     map[*ReferenceCell]*ReferenceCell
	bindValue func(Value) Value
	bindSlot  func(*ZVal) *ZVal
}

func NewArrayOverlayScope() *ArrayOverlayScope {
	return &ArrayOverlayScope{arrays: make(map[*ArrayValue]*ArrayValue), refs: make(map[*ZVal]*ZVal), cells: make(map[*ReferenceCell]*ReferenceCell)}
}

func (s *ArrayOverlayScope) copySlot(source *ZVal, bindGuard func(*ReferenceGuard) *ReferenceGuard) *ZVal {
	if source == nil {
		return nil
	}
	if target := s.refs[source]; target != nil {
		return target
	}
	target := CopyReferenceBucket(source)
	target.reference = nil
	s.refs[source] = target
	if source.reference != nil {
		if cell := s.cells[source.reference]; cell != nil {
			target.reference = cell
			return target
		}
		cell := &ReferenceCell{count: source.reference.count}
		s.cells[source.reference] = cell
		target.reference = cell
		if s.bindValue != nil {
			cell.value = s.bindValue(source.ReadValue())
		} else if array, ok := source.ReadValue().(*ArrayValue); ok {
			cell.value = s.Array(array)
		} else {
			cell.value = source.ReadValue()
		}
		if bindGuard != nil {
			cell.guard = bindGuard(source.Guard())
		} else {
			cell.guard = source.Guard()
		}
	} else if s.bindValue != nil {
		target.StoreRaw(s.bindValue(source.ReadValue()))
	} else if array, ok := source.ReadValue().(*ArrayValue); ok {
		target.StoreRaw(s.Array(array))
	}
	return target
}

// OverlayArrayStore records local slots and tombstones without copying the
// parent's entry buffer. A writable slot is promoted before it leaves this API.
type OverlayArrayStore struct {
	parent  FlatArrayStore
	scope   *ArrayOverlayScope
	changed map[string]*ZVal
	deleted map[string]bool
	order   []string
	added   map[string]int
	length  int
}

func (s *ArrayOverlayScope) Array(parent *ArrayValue) *ArrayValue {
	if parent == nil {
		return nil
	}
	if previous := s.arrays[parent]; previous != nil {
		return previous
	}
	// Layering a live overlay requires a value snapshot. Worker parents are flat.
	parent.Materialize()
	snapshot := parent.flatArrayStore
	snapshot.ensureIndex()
	o := &OverlayArrayStore{parent: snapshot, scope: s, length: len(snapshot.entries)}
	result := &ArrayValue{flatArrayStore: FlatArrayStore{
		appendKeyKnown: true, intKeySeen: snapshot.intKeySeen,
		nextIntKey: snapshot.nextIntKey, iterator: snapshot.iterator, overlay: o,
	}, IndirectOverloadClass: parent.IndirectOverloadClass}
	s.arrays[parent] = result
	return result
}

func overlayKey(slot *ZVal, position int) string {
	if slot == nil || slot.IsPackedIntSlot() {
		return strconv.Itoa(position)
	}
	return slot.Name
}

func (o *OverlayArrayStore) source(key string) (*ZVal, bool) {
	if o.deleted[key] {
		return nil, false
	}
	return o.parent.LookupZValByStringKey(key)
}

func (o *OverlayArrayStore) promote(key string, slot *ZVal) *ZVal {
	if local := o.changed[key]; local != nil {
		return local
	}
	var local *ZVal
	if o.scope.bindSlot != nil && (slot.RefCount() > 0 || slot.Guard() != nil) {
		local = o.scope.bindSlot(slot)
	}
	if slot.RefCount() > 0 {
		if local == nil {
			local = o.scope.refs[slot]
		}
	}
	if local == nil {
		if slot.RefCount() > 0 {
			local = o.scope.copySlot(slot, nil)
		} else {
			local = CopyZValKeepName(slot, slot.ReadValue())
			if child, ok := slot.ReadValue().(*ArrayValue); ok {
				local.StoreRaw(o.scope.Array(child))
			} else if o.scope.bindValue != nil {
				local.StoreRaw(o.scope.bindValue(slot.ReadValue()))
			}
		}
	}
	local = CopyReferenceBucket(local)
	// Packed positions become explicit keys before the parent is compacted.
	if local.IsPackedIntSlot() {
		local.Name = key
	}
	if o.changed == nil {
		o.changed = make(map[string]*ZVal)
	}
	o.changed[key] = local
	return local
}

func (o *OverlayArrayStore) get(key string) (*ZVal, bool) {
	if local := o.changed[key]; local != nil {
		return local, true
	}
	parent, ok := o.source(key)
	if !ok {
		return nil, false
	}
	return o.promote(key, parent), true
}

func (o *OverlayArrayStore) set(key string, value Value) {
	if local, ok := o.get(key); ok {
		local.StoreRaw(value)
		return
	}
	if o.changed == nil {
		o.changed = make(map[string]*ZVal)
	}
	if o.added == nil {
		o.added = make(map[string]int)
	}
	var slot *ZVal
	if key == "" {
		slot = NewEmptyStringKeyZVal(value)
	} else {
		slot = NewNamedZVal(key, value)
	}
	o.changed[key] = slot
	o.added[key] = len(o.order)
	o.order = append(o.order, key)
	o.length++
}

func (o *OverlayArrayStore) unset(key string) {
	_, local := o.changed[key]
	_, parent := o.source(key)
	if !local && !parent {
		return
	}
	if slot, ok := o.get(key); ok {
		slot.ReleaseRefSlot()
	}
	delete(o.changed, key)
	delete(o.added, key)
	if o.deleted == nil {
		o.deleted = make(map[string]bool)
	}
	o.deleted[key] = true
	o.length--
}

func (o *OverlayArrayStore) rangeSlots(yield func(int, *ZVal) bool) {
	position := 0
	for originalPosition, slot := range o.parent.entries {
		key := overlayKey(slot, originalPosition)
		if _, added := o.added[key]; added || o.deleted[key] {
			continue
		}
		if slot == nil {
			if !yield(position, nil) {
				return
			}
		} else if !yield(position, o.promote(key, slot)) {
			return
		}
		position++
	}
	for orderPosition, key := range o.order {
		if current, ok := o.added[key]; !ok || current != orderPosition {
			continue
		}
		if !yield(position, o.changed[key]) {
			return
		}
		position++
	}
}

func (o *OverlayArrayStore) position(key string) int {
	position := 0
	for original, slot := range o.parent.entries {
		candidate := overlayKey(slot, original)
		if _, added := o.added[candidate]; added || o.deleted[candidate] {
			continue
		}
		if candidate == key {
			return position
		}
		position++
	}
	for original, candidate := range o.order {
		if current, ok := o.added[candidate]; !ok || current != original {
			continue
		}
		if candidate == key {
			return position
		}
		position++
	}
	return -1
}

// Materialize is the explicit boundary for contiguous views and global edits.
// It keeps promoted references and child overlays, including detached aliases.
func (a *FlatArrayStore) Materialize() {
	if a.overlay == nil {
		return
	}
	o := a.overlay
	entries := make([]*ZVal, 0, o.length)
	o.rangeSlots(func(_ int, slot *ZVal) bool { entries = append(entries, slot); return true })
	a.entries = entries
	a.overlay = nil
	a.invalidateIndex()
}

func (a *ArrayValue) IsOverlay() bool { return a.overlay != nil }
