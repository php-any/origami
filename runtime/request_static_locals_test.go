package runtime

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestFunctionStaticLocalsBelongToRequest(t *testing.T) {
	file := filepath.Join(t.TempDir(), "static-locals.php")
	if err := os.WriteFile(file, []byte("<?php function requestCounter() { static $n = 0; return ++$n; }"), 0600); err != nil {
		t.Fatal(err)
	}
	base := NewVM(parser.NewParser()).(*VM)
	if _, ctl := base.LoadAndRun(file); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	fn, ok := base.GetFunc("requestCounter")
	if !ok {
		t.Fatal("function missing")
	}
	check := func(vm data.VM, want int) {
		ctx := vm.CreateContext(fn.GetVariables())
		value, ctl := fn.Call(ctx)
		if ctl != nil {
			t.Error(ctl.AsString())
			return
		}
		got, ok := value.(*data.IntValue)
		if !ok || got.Value != want {
			t.Errorf("counter=%v want %d", value, want)
		}
	}
	check(base, 1)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			request := NewRequestVM(base)
			check(request, 1)
			check(request, 2)
		}()
	}
	group.Wait()
	check(base, 2)
}
