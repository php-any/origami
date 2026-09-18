package runtime

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestRequestDeadlineIsolatedPerGoroutine(t *testing.T) {
	restore := BeginRequestOutput()
	defer restore()

	if RequestDeadlineExceeded() {
		t.Fatal("no deadline yet")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stop := BeginRequestDeadline(ctx)
	defer stop()
	if RequestDeadlineExceeded() {
		t.Fatal("deadline should not be exceeded")
	}
	if !InHTTPRequest() {
		t.Fatal("request goroutine should count as HTTP")
	}
	cancel()
	if !RequestDeadlineExceeded() {
		t.Fatal("canceled deadline must be exceeded")
	}
}

func TestConcurrentRequestDeadlinesDoNotClobber(t *testing.T) {
	const n = 32
	var wg sync.WaitGroup
	errs := make(chan string, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			restore := BeginRequestOutput()
			defer restore()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stop := BeginRequestDeadline(ctx)
			defer stop()
			if !SetRequestPHPDeadline(30) {
				errs <- "php deadline not isolated"
				return
			}
			if RequestDeadlineExceeded() {
				errs <- "deadline fired too early"
				return
			}
			if i%2 == 0 {
				cancel()
				if !RequestDeadlineExceeded() {
					errs <- "canceled request must be exceeded"
				}
			} else if RequestDeadlineExceeded() {
				errs <- "sibling cancel leaked"
			}
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatal(e)
	}
}

func TestSharedVMCallStackIsolatedAcrossGoroutines(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	const n = 32
	var wg sync.WaitGroup
	errs := make(chan string, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			restore := BeginRequestOutput()
			defer restore()

			marker := "req-" + strconv.Itoa(i)
			base.PushCallFrame(data.CallFrame{Function: marker})
			defer base.PopCallFrame()

			if d := base.EnterCall(); d != 1 {
				errs <- "depth=" + strconv.Itoa(d) + " want 1 for " + marker
				base.LeaveCall()
				return
			}
			got := base.SnapshotCallStack()
			if len(got) != 1 || got[0].Function != marker {
				errs <- "stack leaked or wrong: " + marker
			}
			base.LeaveCall()
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatal(e)
	}
}

func TestContextPoolDoesNotMutateEscapedZVal(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	parent := vm.CreateContext(nil).(*Context)
	vars := []data.Variable{data.NewVariable("x", 0, nil)}

	frame := parent.CreateContext(vars).(*Context)
	slot := frame.GetIndexZVal(0)
	slot.Value = data.NewIntValue(42)
	slot.Defined = true
	frame.MarkEscaped()
	frame.ReleasePooled()

	if iv, ok := slot.Value.(*data.IntValue); !ok || iv.Value != 42 {
		t.Fatalf("escaped ZVal mutated after ReleasePooled: %#v", slot.Value)
	}

	// 池应仍能给出新帧，且新槽不是 escaped 那个 ZVal
	next := parent.CreateContext(vars).(*Context)
	ns := next.GetIndexZVal(0)
	if ns == slot {
		t.Fatal("pooled frame reused escaped ZVal pointer")
	}
	next.ReleasePooled()
}

func TestTempVMAndSharedVMCallStacksDoNotMix(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	tmp := NewTempVM(base).(*TempVM)

	restore := BeginRequestOutput()
	defer restore()

	base.PushCallFrame(data.CallFrame{Function: "shared"})
	defer base.PopCallFrame()
	tmp.PushCallFrame(data.CallFrame{Function: "temp"})
	defer tmp.PopCallFrame()

	if got := tmp.SnapshotCallStack(); len(got) != 1 || got[0].Function != "temp" {
		t.Fatalf("TempVM stack: %#v", got)
	}
	if got := base.SnapshotCallStack(); len(got) != 1 || got[0].Function != "shared" {
		t.Fatalf("shared VM stack: %#v", got)
	}
}
