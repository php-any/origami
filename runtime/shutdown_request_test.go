package runtime

import (
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

type shutdownTestCallback struct{ fn func(data.Context) }

func (*shutdownTestCallback) GetName() string               { return "shutdown_test" }
func (*shutdownTestCallback) GetParams() []data.GetValue    { return nil }
func (*shutdownTestCallback) GetVariables() []data.Variable { return nil }
func (f *shutdownTestCallback) Call(ctx data.Context) (data.GetValue, data.Control) {
	f.fn(ctx)
	return data.NewNullValue(), nil
}

func TestRequestVMShutdownIsolatedAndOnce(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	first, second := NewRequestVM(base), NewRequestVM(base)
	var order []string
	callback := func(s string) data.Value {
		return data.NewFuncValue(&shutdownTestCallback{fn: func(ctx data.Context) {
			order = append(order, s)
			if s == "first" {
				first.RunShutdownCallbacks()
				first.AddShutdownCallback(data.NewFuncValue(&shutdownTestCallback{fn: func(data.Context) { order = append(order, "late") }}))
			}
		}})
	}
	base.AddShutdownCallback(callback("base"))
	first.AddShutdownCallback(callback("first"))
	second.AddShutdownCallback(callback("second"))
	first.RunShutdownCallbacks()
	first.RunShutdownCallbacks()
	second.RunShutdownCallbacks()
	base.RunShutdownCallbacks()
	if got := len(order); got != 4 {
		t.Fatalf("order=%v", order)
	}
	for i, want := range []string{"first", "late", "second", "base"} {
		if order[i] != want {
			t.Fatalf("order=%v", order)
		}
	}
}

func TestRequestShutdownFromBaseVMDoesNotLeak(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer BeginRequestOutput()()
			requestVM := NewRequestVM(base)
			calls := 0
			base.AddShutdownCallback(data.NewFuncValue(&shutdownTestCallback{fn: func(ctx data.Context) {
				if ctx.GetVM() != requestVM {
					t.Error("shutdown did not use request VM")
				}
				calls++
			}}))
			requestVM.RunShutdownCallbacks()
			base.RunShutdownCallbacks()
			if calls != 1 {
				t.Errorf("calls=%d", calls)
			}
		}()
	}
	wg.Wait()
	if len(base.shutdownCallbacks) != 0 {
		t.Fatal("request callbacks escaped to base VM")
	}
}

func TestStartupContextInheritsRequestExecutionState(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	startup := base.CreateContext(nil)
	startup.(*Context).EnterCall()
	startup.(*Context).LeaveCall()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer BeginRequestOutput()()
			child := startup.CreateContext(nil).(*Context)
			if child.call != currentRequestCallState() {
				t.Error("startup context retained its call state")
			}
			if depth := child.EnterCall(); depth != 1 {
				t.Errorf("depth=%d", depth)
			}
			child.LeaveCall()
			if ctx := startup.CreateBaseContext().(*Context); ctx.call != currentRequestCallState() {
				t.Error("base context retained startup call state")
			}
		}()
	}
	wg.Wait()
}
