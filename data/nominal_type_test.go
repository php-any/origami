package data

import (
	"errors"
	"fmt"
	"testing"
)

type nominalThrowContext struct {
	Context
	vm VM
}

func (c *nominalThrowContext) GetVM() VM { return c.vm }

func TestInternalThrowUsesRegisteredNominalDeclaration(t *testing.T) {
	parent := "Exception"
	class := &nominalClassFixture{name: "Native\\DeclaredException", parent: &parent}
	vm := &nominalVMFixture{classes: map[string]ClassStmt{class.name: class}}
	ctx := &nominalThrowContext{vm: vm}
	thrown := NewErrorThrowByName(nil, errors.New("native"), class.name).(*ThrowValue)
	for _, name := range []string{class.name, "Exception", "Throwable"} {
		if !NominalValueMatches(thrown, name, ctx) {
			t.Fatalf("registered throw does not match %s", name)
		}
	}
	if NominalValueMatches(thrown, "Error", ctx) {
		t.Fatal("Exception matched Error branch")
	}
}

type nominalClassFixture struct {
	ClassStmt
	name       string
	parent     *string
	implements []string
}

func (c *nominalClassFixture) GetName() string         { return c.name }
func (c *nominalClassFixture) GetExtend() *string      { return c.parent }
func (c *nominalClassFixture) GetImplements() []string { return c.implements }

type nominalInterfaceFixture struct {
	InterfaceStmt
	name    string
	parents []string
}

func (i *nominalInterfaceFixture) GetName() string      { return i.name }
func (i *nominalInterfaceFixture) GetExtends() []string { return i.parents }

type nominalVMFixture struct {
	VM
	classes    map[string]ClassStmt
	interfaces map[string]InterfaceStmt
}

func (v *nominalVMFixture) GetClass(name string) (ClassStmt, bool) {
	c, ok := v.classes[name]
	return c, ok
}
func (v *nominalVMFixture) GetInterface(name string) (InterfaceStmt, bool) {
	i, ok := v.interfaces[name]
	return i, ok
}
func (v *nominalVMFixture) LoadPkg(string) (GetValue, Control) {
	panic("nominal comparison must never autoload")
}

func TestNominalAncestry(t *testing.T) {
	vm := &nominalVMFixture{classes: map[string]ClassStmt{}, interfaces: map[string]InterfaceStmt{}}
	for i := 0; i < 40; i++ {
		name := fmt.Sprintf("I%d", i)
		var parents []string
		if i != 0 {
			parents = []string{fmt.Sprintf("I%d", i-1)}
		}
		vm.interfaces[name] = &nominalInterfaceFixture{name: name, parents: parents}
	}
	base := &nominalClassFixture{name: "Base", implements: []string{"I39"}}
	parent := "Base"
	child := &nominalClassFixture{name: "Child", parent: &parent}
	vm.classes["Base"] = base
	for _, target := range []string{"child", "\\BASE", "i0", "I39"} {
		if !NominalIsA(child, target, vm) {
			t.Errorf("Child should satisfy %s", target)
		}
	}
	if NominalIsA(child, "Missing", vm) {
		t.Fatal("unknown target matched")
	}
	vm.interfaces["CycleA"] = &nominalInterfaceFixture{name: "CycleA", parents: []string{"CycleB"}}
	vm.interfaces["CycleB"] = &nominalInterfaceFixture{name: "CycleB", parents: []string{"CycleA", "I0"}}
	if InterfaceIsA("CycleA", "Missing", vm) {
		t.Fatal("cyclic interface matched unknown target")
	}
	if !InterfaceIsA("CycleA", "I0", vm) {
		t.Fatal("cycle hid another parent")
	}
	parent = "Child"
	vm.classes["Child"] = child
	if NominalIsA(child, "Missing", vm) {
		t.Fatal("cyclic class matched unknown target")
	}
}

func TestInternalTypeHierarchy(t *testing.T) {
	for _, tc := range []struct {
		source, target string
		want           bool
	}{
		{"ValueError", "Exception", false}, {"ValueError", "Error", true},
		{"ParseError", "CompileError", true}, {"ParseError", "Throwable", true},
		{"BadMethodCallException", "BadFunctionCallException", true},
		{"FiberError", "error", true}, {"\\typeerror", "\\THROWABLE", true},
		{"RuntimeException", "Ns\\Exception", false}, {"Ns\\Error", "Error", false},
		{"Ns\\Iterator", "Traversable", false},
	} {
		if got := InternalTypeIsA(tc.source, tc.target); got != tc.want {
			t.Errorf("%s is %s: %v", tc.source, tc.target, got)
		}
	}
	if TypeNameEqual("ÄName", "äname") {
		t.Fatal("PHP only folds ASCII identifiers")
	}
}

var nominalBenchmarkResult bool

// The same fixture runs against the old Class.Is implementation for comparison.
func BenchmarkNominalClassIs(b *testing.B) {
	vm := &nominalVMFixture{classes: map[string]ClassStmt{}, interfaces: map[string]InterfaceStmt{}}
	vm.interfaces["Root"] = &nominalInterfaceFixture{name: "Root"}
	vm.interfaces["Middle"] = &nominalInterfaceFixture{name: "Middle", parents: []string{"Root"}}
	vm.interfaces["Leaf"] = &nominalInterfaceFixture{name: "Leaf", parents: []string{"Middle"}}
	parent := "Base"
	base := &nominalClassFixture{name: "Base", implements: []string{"Leaf"}}
	child := &nominalClassFixture{name: "Child", parent: &parent}
	vm.classes["Base"] = base
	object := NewClassValue(child, nil)
	object.vm = vm
	for _, target := range []string{"Child", "Base", "Root", "Missing"} {
		b.Run(target, func(b *testing.B) {
			typ := Class{Name: target}
			b.ReportAllocs()
			for b.Loop() {
				nominalBenchmarkResult = typ.Is(object)
			}
		})
	}
}
