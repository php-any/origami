package data

import (
	"testing"
	"unsafe"
)

func TestArrayStorageLayoutBudget(t *testing.T) {
	if size := unsafe.Sizeof(ZVal{}); size > 48 {
		t.Fatalf("ZVal grew beyond 48 bytes: %d", size)
	}
	if size := unsafe.Sizeof(ArrayValue{}); size > 88 {
		t.Fatalf("flat store added array allocation size: %d", size)
	}
}

func TestArrayReferenceSurvivesRemovalAndReordering(t *testing.T) {
	a := NewArrayValue([]Value{NewIntValue(10), NewIntValue(20), NewIntValue(30)}).(*ArrayValue)
	ref := &ArraySlotRef{Slot: a.At(1)}
	a.UnsetKey(NewIntValue(0))
	if ref.AsString() != "20" {
		t.Fatal("reference moved to another slot after unset")
	}
	a.UnsetKey(NewIntValue(1))
	a.AppendValue(NewIntValue(40))
	if ref.AsString() != "20" {
		t.Fatal("detached reference reused a different array entry")
	}
	ref.Slot.Value = NewIntValue(99)
	for _, slot := range a.Range() {
		if slot.Value.AsString() == "99" {
			t.Fatal("detached reference still writes to the array")
		}
	}
}

func TestArrayBulkReorderInvalidatesKeys(t *testing.T) {
	a := NewArrayValue(nil).(*ArrayValue)
	a.SetStringKey("one", NewIntValue(1))
	a.SetStringKey("two", NewIntValue(2))
	_, _ = a.LookupZValByStringKey("one")
	reverse, _ := a.GetMethod("reverse")
	_, ctl := reverse.Call(nil)
	if ctl != nil {
		t.Fatal(ctl)
	}
	for key, want := range map[string]string{"one": "1", "two": "2"} {
		slot, ok := a.LookupZValByStringKey(key)
		if !ok || slot.Value.AsString() != want {
			t.Fatalf("reordered key %q returned the wrong slot", key)
		}
	}
	if a.At(0).Name != "two" {
		t.Fatal("bulk replacement lost insertion order")
	}
}

func TestArrayEditPanicInvalidatesKeys(t *testing.T) {
	a := NewArrayValue(nil).(*ArrayValue)
	a.SetStringKey("one", NewIntValue(1))
	a.SetStringKey("two", NewIntValue(2))
	_, _ = a.LookupZValByStringKey("one")
	func() {
		defer func() { _ = recover() }()
		a.EditSlots(func(slots []*ZVal) {
			slots[0], slots[1] = slots[1], slots[0]
			panic("comparison failed")
		})
	}()
	if slot, ok := a.LookupZValByStringKey("one"); !ok || slot.Value.AsString() != "1" {
		t.Fatal("partial mutation retained a stale key index after panic")
	}
}

var arrayAccessSink int

// Compare the API with the legacy slice loop on the same data and compiler.
func BenchmarkArrayTraversal(b *testing.B) {
	for _, size := range []int{8, 128} {
		values := make([]Value, size)
		for i := range values {
			values[i] = NewIntValue(i)
		}
		a := NewArrayValue(values).(*ArrayValue)
		b.Run(IntArrayKeyName(size)+"/slice", func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				sum := 0
				for i, slot := range a.entries {
					sum += i + slot.Value.(*IntValue).Value
				}
				arrayAccessSink = sum
			}
		})
		b.Run(IntArrayKeyName(size)+"/range", func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				sum := 0
				for i, slot := range a.Range() {
					sum += i + slot.Value.(*IntValue).Value
				}
				arrayAccessSink = sum
			}
		})
		b.Run(IntArrayKeyName(size)+"/span", func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				sum := 0
				for slots, i := a.View(), 0; i < slots.Len(); i++ {
					sum += i + slots.At(i).Value.(*IntValue).Value
				}
				arrayAccessSink = sum
			}
		})
	}
}
