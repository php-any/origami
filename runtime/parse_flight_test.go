package runtime

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestParseFlightSharesColdResult(t *testing.T) {
	cache := &parsedFileCache{}
	var parses atomic.Int32
	start, release := make(chan struct{}), make(chan struct{})
	entry := &parsedPHPFile{program: data.NewIntValue(42)}
	parse := func() (*parsedPHPFile, data.Control) { parses.Add(1); close(start); <-release; return entry, nil }
	var wg sync.WaitGroup
	results := make(chan *parsedPHPFile, 16)
	wg.Add(1)
	go func() { defer wg.Done(); got, _ := cache.load(context.Background(), "file", parse); results <- got }()
	<-start
	for i := 1; i < 16; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); got, _ := cache.load(context.Background(), "file", parse); results <- got }()
	}
	close(release)
	wg.Wait()
	close(results)
	for got := range results {
		if got != entry {
			t.Fatal("loaders did not share AST")
		}
	}
	if parses.Load() != 1 {
		t.Fatalf("parsed %d times", parses.Load())
	}
}

func TestParseFlightCancellationDoesNotCancelLeader(t *testing.T) {
	cache := &parsedFileCache{}
	start, release, leaderDone := make(chan struct{}), make(chan struct{}), make(chan struct{})
	defer func() {
		close(release)
		select {
		case <-leaderDone:
		case <-time.After(time.Second):
			t.Error("leader did not finish")
		}
	}()
	go func() {
		defer close(leaderDone)
		_, _ = cache.load(context.Background(), "file", func() (*parsedPHPFile, data.Control) {
			close(start)
			<-release
			return &parsedPHPFile{program: data.NewIntValue(1)}, nil
		})
	}()
	<-start
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan any, 1)
	go func() {
		defer func() { done <- recover() }()
		_, _ = cache.load(ctx, "file", func() (*parsedPHPFile, data.Control) { panic("waiter must not parse") })
	}()
	cancel()
	select {
	case result := <-done:
		if !data.IsRequestCanceled(result) {
			t.Fatalf("waiter panic = %v", result)
		}
	case <-time.After(time.Second):
		t.Fatal("waiter ignored cancellation")
	}
	select {
	case <-leaderDone:
		t.Fatal("waiter canceled shared leader")
	default:
	}
}

func TestParseFlightPanicAndRecursionReleaseState(t *testing.T) {
	cache := &parsedFileCache{}
	func() {
		defer func() {
			if recover() != "parse failed" {
				t.Fatal("leader panic was lost")
			}
		}()
		_, _ = cache.load(context.Background(), "file", func() (*parsedPHPFile, data.Control) { panic("parse failed") })
	}()
	_, control := cache.load(context.Background(), "file", func() (*parsedPHPFile, data.Control) {
		return cache.load(context.Background(), "file", func() (*parsedPHPFile, data.Control) { panic("recursive parse") })
	})
	if control == nil {
		t.Fatal("recursive parse was accepted")
	}
	entry := &parsedPHPFile{program: data.NewIntValue(1)}
	got, control := cache.load(context.Background(), "file", func() (*parsedPHPFile, data.Control) { return entry, nil })
	if got != entry || control != nil {
		t.Fatal("failed parse prevented later retry")
	}
	cache.flights.Range(func(_, _ any) bool { t.Fatal("flight retained after completion"); return false })
}

func TestClearParseCacheDoesNotRepublishOldGeneration(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	old := vm.parsedFiles.Load()
	start, release := make(chan struct{}), make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = old.load(context.Background(), "file", func() (*parsedPHPFile, data.Control) {
			close(start)
			<-release
			return &parsedPHPFile{program: data.NewIntValue(1)}, nil
		})
	}()
	<-start
	vm.ClearParsedFileCache()
	close(release)
	<-done
	if _, ok := vm.parsedFiles.Load().entries.Load("file"); ok {
		t.Fatal("old parse repopulated cleared cache")
	}
}

func TestFileLoadWaitCancellation(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	file := testPhpLoadPath(t, "waiting.php")
	_, _, finish := vm.beginPhpFileLoad(file)
	defer finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan any, 1)
	go func() {
		defer func() { done <- recover() }()
		defer BeginRequestOutput()()
		defer BeginRequestDeadline(ctx)()
		_ = vm.WaitPhpFileLoad(file)
	}()
	cancel()
	select {
	case result := <-done:
		if !data.IsRequestCanceled(result) {
			t.Fatalf("file wait panic = %v", result)
		}
	case <-time.After(time.Second):
		t.Fatal("autoload wait ignored cancellation")
	}
}

func BenchmarkParseFileCachedHit(b *testing.B) {
	vm := NewVM(parser.NewParser()).(*VM)
	file := filepath.Join(b.TempDir(), "benchmark-cached.php")
	if err := os.WriteFile(file, []byte("<?php echo 1;"), 0600); err != nil {
		b.Fatal(err)
	}
	if _, _, control := vm.ParseFileCached(file); control != nil {
		b.Fatal(control)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = vm.ParseFileCached(file)
	}
}
