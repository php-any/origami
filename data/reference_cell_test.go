package data

import (
	"sync"
	"testing"
	"unsafe"
)

func TestReferenceCellSeparateBucketAndRequestIdentity(t *testing.T) {
	if unsafe.Sizeof(ZVal{}) != 48 {
		t.Fatalf("ZVal grew to %d bytes", unsafe.Sizeof(ZVal{}))
	}
	parent := NewArrayValue(nil).(*ArrayValue)
	slot := NewNamedZVal("left", NewIntValue(1))
	slot.AddRefSlot()
	parent.AppendEntries(slot)
	copy := CloneArrayValue(parent)
	copy.EditReindexing(func([]*ZVal) {})
	if parent.At(0).Name != "left" || !copy.At(0).IsPackedIntSlot() {
		t.Fatal("keys not bucket-owned")
	}
	copy.At(0).StoreRaw(NewIntValue(2))
	if slot.ReadValue().AsString() != "2" {
		t.Fatal("reference lost across copy")
	}
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			scope := NewArrayOverlayScope()
			a, b := scope.Array(parent), scope.Array(copy)
			a.At(0).StoreRaw(NewIntValue(3))
			if b.At(0).ReadValue().AsString() != "3" {
				t.Error("request cell alias lost")
			}
			if a.At(0).Name == b.At(0).Name {
				t.Error("request bucket keys aliased")
			}
		}()
	}
	group.Wait()
	if slot.ReadValue().AsString() != "2" {
		t.Fatal("worker reference mutated")
	}
}
