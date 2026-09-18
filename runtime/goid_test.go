package runtime

import (
	"strconv"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
)

func TestGoidUniqueAcrossGoroutines(t *testing.T) {
	const n = 64
	ids := make(chan uint64, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			ids <- goid()
		}()
	}
	wg.Wait()
	close(ids)
	seen := make(map[uint64]struct{}, n)
	for id := range ids {
		if id == 0 {
			t.Fatal("goid must not be 0")
		}
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate goid %d", id)
		}
		seen[id] = struct{}{}
	}
	if len(seen) != n {
		t.Fatalf("got %d unique goids want %d", len(seen), n)
	}
}

func TestGoidUsesStructOffset(t *testing.T) {
	_ = goid()
	off := goidOff.Load()
	if off == 0 || off == ^uintptr(0) {
		t.Fatalf("cheap goid offset not discovered: %v", off)
	}
	if goid() == 0 {
		t.Fatal("goid must not be 0 after offset init")
	}
}

func TestRequestStaticOverlayIsolatesGoroutines(t *testing.T) {
	const n = 16
	var wg sync.WaitGroup
	errs := make(chan string, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			restore := BeginRequestOutput()
			defer restore()
			marker := data.NewStringValue("req-" + strconv.Itoa(i))
			if !data.StoreRequestStatic("Illuminate\\Container\\Container", "instance", marker) {
				errs <- "store failed"
				return
			}
			got, ok := data.LoadRequestStatic("Illuminate\\Container\\Container", "instance")
			if !ok || got.AsString() != marker.AsString() {
				errs <- "load mismatch"
			}
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatal(e)
	}
	if _, ok := data.LoadRequestStatic("Illuminate\\Container\\Container", "instance"); ok {
		t.Fatal("overlay must not leak after restore")
	}
}

func TestGoidMatchesStackAndStaysStable(t *testing.T) {
	_ = goid()
	want := goidFromStack()
	got := goid()
	if got != want {
		t.Fatalf("goid()=%d stack=%d (getg offset is wrong)", got, want)
	}
	for i := 0; i < 1000; i++ {
		if id := goid(); id != want {
			t.Fatalf("goid changed on same goroutine: %d -> %d", want, id)
		}
	}
}
