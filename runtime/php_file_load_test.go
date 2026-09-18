package runtime

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/php-any/origami/parser"
)

func testPhpLoadPath(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte("<?php\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return normalizePhpFilePath(p)
}

func TestBeginPhpFileLoadFinishUnblocksWaiters(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	file := testPhpLoadPath(t, "stuck.php")
	_, _, finish := vm.beginPhpFileLoad(file)
	if finish == nil {
		t.Fatal("loader should return finish")
	}

	unblocked := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = vm.WaitPhpFileLoad(file)
		close(unblocked)
	}()

	select {
	case <-unblocked:
		t.Fatal("waiter must block until finish")
	case <-time.After(50 * time.Millisecond):
	}

	func() {
		defer finish()
		defer func() { _ = recover() }()
		panic("request canceled")
	}()

	select {
	case <-unblocked:
	case <-time.After(time.Second):
		t.Fatal("waiter still blocked after finish")
	}
	wg.Wait()
}

func TestBeginPhpFileLoadReentrantDoesNotWait(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	file := testPhpLoadPath(t, "reenter.php")
	loaded, wait, finish := vm.beginPhpFileLoad(file)
	if loaded || wait != nil || finish == nil {
		t.Fatalf("first load: loaded=%v wait=%v finish=%v", loaded, wait != nil, finish != nil)
	}
	defer finish()

	loaded, wait, finish2 := vm.beginPhpFileLoad(file)
	if finish2 != nil {
		t.Fatal("reentrant load must not become a second loader")
	}
	if wait != nil {
		t.Fatal("reentrant load must not wait on its own channel")
	}
	if !loaded {
		t.Fatal("reentrant load should report alreadyLoaded")
	}
	if vm.WaitPhpFileLoad(file) {
		t.Fatal("same goroutine WaitPhpFileLoad must not report cache complete")
	}
}

func TestClearPhpFileCacheUnblocksWaiters(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	file := testPhpLoadPath(t, "hot.php")
	_, _, finish := vm.beginPhpFileLoad(file)
	if finish == nil {
		t.Fatal("expected finish")
	}

	done := make(chan struct{})
	go func() {
		_ = vm.WaitPhpFileLoad(file)
		close(done)
	}()
	time.Sleep(30 * time.Millisecond)
	vm.ClearPhpFileCache()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("ClearPhpFileCache left waiters blocked")
	}
}
