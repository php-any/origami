package runtime

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

type countedDefault struct {
	calls *int
	value int
}

func (d *countedDefault) GetValue(data.Context) (data.GetValue, data.Control) {
	*d.calls++
	return data.NewIntValue(d.value), nil
}

func TestInheritedDefaultsEvaluateOnceAndRelinkPerRequest(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	baseCalls, hiddenCalls := 0, 0
	parent := node.NewClassStatement(nil, "InitBase", "", nil, []data.Property{
		node.NewProperty(nil, "inherited", "public", false, &countedDefault{&baseCalls, 10}),
		node.NewProperty(nil, "hidden", "public", false, &countedDefault{&hiddenCalls, 1}),
	}, nil)
	child := node.NewClassStatement(nil, "InitChild", "InitBase", nil, []data.Property{
		node.NewProperty(nil, "hidden", "public", false, data.NewIntValue(2)),
	}, nil)
	for _, declaration := range []data.ClassStmt{parent, child} {
		if ctl := vm.AddClass(declaration); ctl != nil {
			t.Fatal(ctl.AsString())
		}
	}
	check := func(host data.VM, want int) {
		t.Helper()
		value, ctl := child.GetValue(host.CreateContext(nil))
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		object := value.(*data.ClassValue)
		inherited, _ := object.GetProperty("inherited")
		hidden, _ := object.GetProperty("hidden")
		if inherited.(*data.IntValue).Value != want || hidden.(*data.IntValue).Value != 2 {
			t.Fatal("initialization reused another request's parent or an overridden default")
		}
	}
	for i := 0; i < 3; i++ {
		check(vm, 10)
	}
	if baseCalls != 3 || hiddenCalls != 0 {
		t.Fatalf("defaults evaluated %d/%d times", baseCalls, hiddenCalls)
	}
	request := NewRequestVM(vm).(*RequestVM)
	replacement := node.NewClassStatement(nil, "InitBase", "", nil, []data.Property{
		node.NewProperty(nil, "inherited", "public", false, data.NewIntValue(20)),
	}, nil)
	if ctl := request.AddClass(replacement); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	check(request, 20)
	check(vm, 10)
}
