package data

import (
	"fmt"
	"sync"
	"testing"
	"unsafe"
)

func TestTypeArenaCanonicalCompounds(t *testing.T) {
	a := NewTypeArena()
	first := a.Nominal("App\\Contract")
	if a.Nominal("\\aPP\\cONTRACT") != first {
		t.Fatal("ASCII nominal identity differs")
	}
	if a.Nominal("É") == a.Nominal("é") {
		t.Fatal("non-ASCII names folded")
	}
	union := a.Union(TypeString, TypeInt, TypeString)
	if union != a.Union(TypeInt, TypeString) || union != a.Union(union, TypeNever) {
		t.Fatal("union was not interned and normalized")
	}
	if a.Union(union, TypeMixed) != TypeMixed || a.Union(TypeTrue, TypeFalse) != TypeBool {
		t.Fatal("union simplification")
	}
	intersection := a.Intersection(first, a.Nominal("Other"))
	if intersection != a.Intersection(intersection, first) {
		t.Fatal("intersection flattening")
	}
	if !a.Matches(union, NewIntValue(3), nil) || !a.Matches(union, NewStringValue("3"), nil) || a.Matches(union, NewNullValue(), nil) {
		t.Fatal("compound exact checks")
	}
	if a.Matches(TypeFalse, NewBoolValue(true), nil) || !a.Matches(TypeFalse, NewBoolValue(false), nil) {
		t.Fatal("literal false widened to bool")
	}
	if unsafe.Sizeof(TypeRef(0)) != 4 || unsafe.Sizeof(typeNode{}) > 24 {
		t.Fatal("arena layout expanded")
	}
}
func TestTypeArenaConcurrentPublication(t *testing.T) {
	a := NewTypeArena()
	stable := a.Union(TypeInt, TypeFloat, TypeString, TypeArray, TypeObject, TypeBool, TypeNull, TypeCallable)
	view := a.Members(stable)
	old := a.snapshot.Load()
	var group sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		group.Add(1)
		go func(worker int) {
			defer group.Done()
			for i := 0; i < 200; i++ {
				name := fmt.Sprintf("Worker%dClass%d", worker, i)
				ref := a.Nominal(name)
				a.Union(ref, TypeNull)
				if !a.Matches(stable, NewIntValue(i), nil) || a.Name(ref) != name || view.Len() != 8 {
					t.Error("published metadata changed")
				}
			}
		}(worker)
	}
	group.Wait()
	if len(old.nodes) != 1 || view.At(0) != TypeNull {
		t.Fatal("old snapshot or borrowed view mutated")
	}
}

var typeArenaBenchSink bool

func BenchmarkTypeArenaMatch(b *testing.B) {
	a := NewTypeArena()
	union := a.Union(TypeInt, TypeFloat, TypeString, TypeArray, TypeObject, TypeBool, TypeNull, TypeCallable)
	value := NewIntValue(3)
	for _, test := range []struct {
		name string
		ref  TypeRef
	}{{"scalar", TypeInt}, {"eight_member_union", union}} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				typeArenaBenchSink = a.Matches(test.ref, value, nil)
			}
		})
	}
}
