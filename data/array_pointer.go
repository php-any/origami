package data

// The PHP internal pointer is independent of foreach iteration. An invalid
// pointer stays invalid until reset/end; moving it backwards cannot revive it.
func (a *ArrayValue) PointerValue() Value {
	if slot := a.At(a.iterator); slot != nil {
		return slot.Value
	}
	return NewBoolValue(false)
}
func (a *ArrayValue) PointerKey() Value {
	if slot := a.At(a.iterator); slot != nil {
		return slot.PHPArrayKey(a.iterator)
	}
	return NewNullValue()
}
func (a *ArrayValue) ResetPointer() Value { a.iterator = 0; return a.PointerValue() }
func (a *ArrayValue) EndPointer() Value   { a.iterator = a.Len() - 1; return a.PointerValue() }
func (a *ArrayValue) MovePointer(delta int) Value {
	if a.iterator >= 0 && a.iterator < a.Len() {
		a.iterator += delta
	}
	return a.PointerValue()
}

func (o *ObjectValue) PointerValue() Value {
	_, value, ok := o.property.GetByIndex(o.iterator)
	if ok {
		return value
	}
	return NewBoolValue(false)
}
func (o *ObjectValue) PointerKey() Value   { value, _ := o.Key(nil); return value }
func (o *ObjectValue) ResetPointer() Value { o.iterator = 0; return o.PointerValue() }
func (o *ObjectValue) EndPointer() Value   { o.iterator = o.property.Len() - 1; return o.PointerValue() }
func (o *ObjectValue) MovePointer(delta int) Value {
	if o.iterator >= 0 && o.iterator < o.property.Len() {
		o.iterator += delta
	}
	return o.PointerValue()
}
