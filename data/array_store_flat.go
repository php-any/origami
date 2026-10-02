package data

// FlatArrayStore owns the ordered entries, key index and PHP automatic integer
// key state. It is embedded privately in ArrayValue, so flat access needs no
// interface dispatch, lock, additional pointer or allocation.
type FlatArrayStore struct {
	entries        []*ZVal
	keyIndex       map[string]int
	idxLen         int
	packed         bool
	appendKeyKnown bool
	intKeySeen     bool
	nextIntKey     int
	iterator       int
}

// The alias keeps the embedded field private while promoting the store's API.
type flatArrayStore = FlatArrayStore

func NewArrayValueFromSlotsWithProvenance(slots []*ZVal, class string) *ArrayValue {
	a := NewArrayValueFromSlots(slots)
	a.IndirectOverloadClass = class
	return a
}
