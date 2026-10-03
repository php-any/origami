package data

import (
	"sync"
	"testing"
)

type descriptorClassFixture struct{ nominalClassFixture }

func (*descriptorClassFixture) GetMethods() []Method        { return nil }
func (*descriptorClassFixture) GetPropertyList() []Property { return nil }

type descriptorInterfaceFixture struct{ nominalInterfaceFixture }

func (*descriptorInterfaceFixture) GetMethods() []Method { return nil }

type descriptorAliasFixture struct {
	ClassStmt
	name string
}

func (c *descriptorAliasFixture) GetName() string          { return c.name }
func (c *descriptorAliasFixture) OriginalClass() ClassStmt { return c.ClassStmt }
func TestDescriptorPublicationDependenciesAndSnapshots(t *testing.T) {
	r := NewClassRegistry(nil)
	parent := "DescriptorBase"
	child := &descriptorClassFixture{nominalClassFixture: nominalClassFixture{name: "DescriptorChild", parent: &parent}}
	r.PublishClass(child)
	childID := Symbols.Intern(child.GetName())
	rootID := Symbols.Intern("DescriptorRoot")
	old, _ := r.Descriptor(childID)
	base := &descriptorClassFixture{nominalClassFixture: nominalClassFixture{name: parent, implements: []string{"DescriptorLeaf"}}}
	r.PublishClass(base)
	r.PublishInterface(&descriptorInterfaceFixture{nominalInterfaceFixture{name: "DescriptorLeaf", parents: []string{"DescriptorRoot"}}})
	r.PublishInterface(&descriptorInterfaceFixture{nominalInterfaceFixture{name: "DescriptorRoot"}})
	current, _ := r.Descriptor(childID)
	if old.IsA(ClassID(rootID)) || !current.IsA(ClassID(rootID)) {
		t.Fatal("dependency publication changed old descriptor or missed ancestor")
	}
	frozen := r.Snapshot()
	lateParent := "RequestLateDescriptorParent"
	r.PublishClass(&descriptorClassFixture{nominalClassFixture: nominalClassFixture{name: "RequestLateDescriptorChild", parent: &lateParent}})
	lateFrozen := r.Snapshot()
	lateRequest := NewClassRegistry(lateFrozen)
	lateRequest.PublishClass(&descriptorClassFixture{nominalClassFixture: nominalClassFixture{name: lateParent, implements: []string{"RequestLateAncestor"}}})
	lateChild, _ := lateRequest.Descriptor(Symbols.Intern("RequestLateDescriptorChild"))
	if !lateChild.IsA(ClassID(Symbols.Intern("RequestLateAncestor"))) {
		t.Fatal("inherited dependency snapshot was not refreshed")
	}
	oldChild, _ := r.Descriptor(Symbols.Intern("RequestLateDescriptorChild"))
	if oldChild.IsA(ClassID(Symbols.Intern("RequestLateAncestor"))) {
		t.Fatal("request linking changed worker descriptor")
	}
	request := NewClassRegistry(frozen)
	alias := &descriptorAliasFixture{ClassStmt: base, name: "DescriptorAlias"}
	request.PublishClass(alias)
	if matched, linked := request.IsA(child, Symbols.Intern(alias.name)); !matched || !linked {
		t.Fatal("alias lost canonical class identity")
	}
	if _, ok := frozen.FindClass(alias.name); ok {
		t.Fatal("request declaration changed frozen registry")
	}
	if got, ok := request.Snapshot().FindClass("\\descriptoralias"); !ok || got != alias {
		t.Fatal("case folded alias lookup failed")
	}
	missing := Symbols.Intern("DescriptorMissing")
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			local := NewClassRegistry(frozen)
			local.PublishClass(alias)
			for j := 0; j < 50; j++ {
				if matched, _ := local.IsA(child, missing); matched {
					t.Error("unknown ancestor matched")
				}
				if !current.IsA(ClassID(rootID)) {
					t.Error("immutable descriptor changed")
				}
			}
		}()
	}
	group.Wait()
}
func TestIDMapPersistentPresence(t *testing.T) {
	var original IDMap[int]
	one := original.With(0, 0).With(64, 7)
	two := one.With(64, 9).With(300, 11)
	if _, ok := original.Get(0); ok {
		t.Fatal("empty map mutated")
	}
	if value, ok := one.Get(0); !ok || value != 0 {
		t.Fatal("zero value lacks presence")
	}
	if value, _ := one.Get(64); value != 7 {
		t.Fatal("old page mutated")
	}
	if value, _ := two.Get(64); value != 9 {
		t.Fatal("new page missing")
	}
}
func BenchmarkDescriptorNominalMatch(b *testing.B) {
	r := NewClassRegistry(nil)
	child := &descriptorClassFixture{nominalClassFixture: nominalClassFixture{name: "BenchDescriptor", implements: []string{"BenchDescriptorLeaf"}}}
	r.PublishClass(child)
	r.PublishInterface(&descriptorInterfaceFixture{nominalInterfaceFixture{name: "BenchDescriptorLeaf", parents: []string{"BenchDescriptorRoot"}}})
	target := Symbols.Intern("BenchDescriptorRoot")
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		nominalBenchmarkResult, _ = r.IsA(child, target)
	}
}
