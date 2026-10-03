package data

import (
	"sync"
	"testing"
)

func TestArrayOverlayLocalEditsAndMaterialization(t *testing.T) {
	parent := NewArrayValue(nil).(*ArrayValue)
	parent.SetIntKey(-3, NewIntValue(1))
	parent.SetStringKey("", NewIntValue(2))
	parent.SetStringKey("x", NewIntValue(3))
	a := NewArrayOverlayScope().Array(parent)
	a.SetStringKey("x", NewIntValue(30))
	a.UnsetKey(NewStringValue(""))
	a.SetStringKey("", NewIntValue(20))
	a.AppendValue(NewIntValue(4))
	if !a.IsOverlay() || a.Len() != 4 || a.NextAppendIntKey() != -1 {
		t.Fatal("local edit materialized or lost automatic key state")
	}
	var keys []string
	for i, slot := range a.Range() {
		keys = append(keys, slot.PHPArrayKey(i).AsString())
	}
	if got := keys; len(got) != 4 || got[0] != "-3" || got[1] != "x" || got[2] != "" || got[3] != "-2" {
		t.Fatalf("order: %v", got)
	}
	if slot, _ := parent.LookupZValByStringKey("x"); slot.ReadValue().AsString() != "3" {
		t.Fatal("parent mutated")
	}
	a.EditReindexing(func(slots []*ZVal) { slots[0], slots[3] = slots[3], slots[0] })
	if a.IsOverlay() || a.NextAppendIntKey() != 4 {
		t.Fatal("bulk edit did not materialize and reindex")
	}
	if parent.Len() != 3 || parent.NextAppendIntKey() != -2 {
		t.Fatal("bulk edit mutated parent")
	}
}

func TestArrayOverlayNestedReferencesCyclesAndConcurrentRequests(t *testing.T) {
	child := NewArrayValue(nil).(*ArrayValue)
	child.SetStringKey("value", NewIntValue(1))
	parent := NewArrayValue(nil).(*ArrayValue)
	parent.SetStringKey("child", child)
	parent.SetStringKey("alias", child)
	parent.SetStringKey("self", parent)
	shared := NewNamedZVal("ref", NewIntValue(10))
	shared.AddRefSlot()
	parent.AppendEntries(shared)
	other := NewArrayValueFromSlots([]*ZVal{shared})
	var wg sync.WaitGroup
	for n := 0; n < 8; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scope := NewArrayOverlayScope()
			a := scope.Array(parent)
			b := scope.Array(other)
			x, _ := a.LookupZValByStringKey("child")
			y, _ := a.LookupZValByStringKey("alias")
			if x.ReadValue() != y.ReadValue() || x.ReadValue() == child {
				t.Error("nested alias not request-owned")
			}
			self, _ := a.LookupZValByStringKey("self")
			if self.ReadValue() != a {
				t.Error("cycle not preserved")
			}
			x.ReadValue().(*ArrayValue).SetStringKey("value", NewIntValue(9))
			r, _ := a.LookupZValByStringKey("ref")
			r.StoreRaw(NewIntValue(99))
			r2, _ := b.LookupZValByStringKey("ref")
			if r.ReferenceIdentity() != r2.ReferenceIdentity() || r2.ReadValue().AsString() != "99" {
				t.Error("request reference alias lost")
			}
			a.UnsetKey(NewStringValue("ref"))
			a.SetStringKey("ref", NewIntValue(50))
			if r.ReadValue().AsString() != "99" {
				t.Error("detached reference reused")
			}
		}()
	}
	wg.Wait()
	if z, _ := child.LookupZValByStringKey("value"); z.ReadValue().AsString() != "1" || shared.ReadValue().AsString() != "10" {
		t.Fatal("request modified parent graph")
	}
}

func BenchmarkArrayOverlayReadWrite(b *testing.B) {
	parent := NewArrayValue(nil).(*ArrayValue)
	for i := 0; i < 128; i++ {
		parent.SetStringKey("key"+IntArrayKeyName(i), NewIntValue(i))
	}
	for _, overlay := range []bool{false, true} {
		name := "clone"
		if overlay {
			name = "overlay"
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var a *ArrayValue
				if overlay {
					a = NewArrayOverlayScope().Array(parent)
				} else {
					a = CloneArrayValue(parent)
				}
				a.SetStringKey("key64", NewIntValue(1000))
				slot, _ := a.LookupZValByStringKey("key64")
				arrayAccessSink = slot.ReadValue().(*IntValue).Value
			}
		})
	}
}
