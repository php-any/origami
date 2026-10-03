package runtime_test

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
	"sync"
	"testing"
)

func TestRequestObjectGraphRetainedAliasesAndReferences(t *testing.T) {
	base := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	php.Load(base)
	value, ctl := base.LoadAndRun("../tests/php/request_object_graph_test.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	parent := value.(*data.ClassValue)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			defer runtime.BeginRequestOutput()()
			request := runtime.NewRequestVM(base).(*runtime.RequestVM)
			ctx := request.CreateContext(nil)
			scoped := request.RequestObjectScope().Object(parent)
			method, _ := scoped.GetMethod("advance")
			frame := data.WrapMethodFrame(ctx.CreateContext(method.GetVariables()), scoped, scoped.Class, scoped.Class)
			result, ctl := method.Call(frame)
			if ctl != nil {
				t.Error(ctl.AsString())
				return
			}
			if got := result.(data.Value).AsString(); got != "1:1:2:nested,gone" {
				t.Errorf("request graph result %s", got)
			}
		}()
	}
	group.Wait()
	child, _ := parent.GetProperty("child")
	count, _ := child.(*data.ClassValue).GetProperty("count")
	if count.AsString() != "0" {
		t.Fatal("retained child mutated")
	}
}
