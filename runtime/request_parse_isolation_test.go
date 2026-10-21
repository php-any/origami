package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestRequestParseOwnershipAndReload(t *testing.T) {
	file := filepath.Join(t.TempDir(), "Controller.php")
	write := func(version string) {
		t.Helper()
		source := fmt.Sprintf(`<?php namespace RequestOwned; class Controller { public static function value() { return '%s'; } } return Controller::value();`, version)
		if err := os.WriteFile(file, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	base := NewVM(parser.NewParser()).(*VM)
	write("v1")
	first := NewRequestVM(base).(*RequestVM)
	read := func(vm *RequestVM, want string) {
		t.Helper()
		value, ctl := vm.LoadInCallerContext(vm.CreateContext(nil), file)
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		if got := value.(data.Value).AsString(); got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	}
	read(first, "v1")
	key := normalizePhpFilePath(file)
	oldAST := first.parsedFiles[key].program
	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	write("v2")
	// New requests must reread even when size and timestamp are unchanged.
	if err := os.Chtimes(file, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	second := NewRequestVM(base).(*RequestVM)
	read(second, "v2")
	read(first, "v1")
	if second.parsedFiles[key].program == oldAST {
		t.Fatal("requests share AST")
	}
	if _, ok := base.GetClass("RequestOwned\\Controller"); ok {
		t.Fatal("class leaked into Base")
	}
	if len(base.parsedFiles.Load().generation().entries) != 0 {
		t.Fatal("request parse leaked into Base cache")
	}
	if _, ok := base.parser.GetClassPathManager().FindClassFile("RequestOwned\\Controller"); ok {
		t.Fatal("namespace leaked into Base")
	}
	third := NewRequestVM(base).(*RequestVM)
	if _, ok := third.requestClassPaths().FindClassFile("RequestOwned\\Controller"); ok {
		t.Fatal("namespace leaked into another request")
	}
	first.InvalidateParsedFile(file)
	if second.parsedFiles[key] == nil {
		t.Fatal("invalidating one request affected another")
	}
}

func TestConcurrentRequestParsesAreDistinct(t *testing.T) {
	file := filepath.Join(t.TempDir(), "Controller.php")
	if err := os.WriteFile(file, []byte(`<?php namespace ParallelOwned; class Controller {} return 42;`), 0600); err != nil {
		t.Fatal(err)
	}
	base := NewVM(parser.NewParser()).(*VM)
	const count = 12
	requests := make([]*RequestVM, count)
	for i := range requests {
		requests[i] = NewRequestVM(base).(*RequestVM)
	}
	var group sync.WaitGroup
	for _, request := range requests {
		group.Add(1)
		go func(vm *RequestVM) {
			defer group.Done()
			defer vm.Activate()()
			if _, ctl := vm.LoadAndRun(file); ctl != nil {
				t.Error(ctl.AsString())
			}
		}(request)
	}
	group.Wait()
	seen := make(map[data.GetValue]bool)
	for _, request := range requests {
		entry := request.parsedFiles[normalizePhpFilePath(file)]
		if entry == nil {
			t.Fatal("missing request parse")
		}
		if seen[entry.program] {
			t.Fatal("concurrent requests share AST")
		}
		seen[entry.program] = true
	}
	if _, ok := base.GetClass("ParallelOwned\\Controller"); ok {
		t.Fatal("class leaked")
	}
}

func TestRequestPathsInheritStartupAndKeepAdditionsLocal(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.php")
	if err := os.WriteFile(file, []byte(`<?php namespace StartupPath; class Controller {}`), 0600); err != nil {
		t.Fatal(err)
	}
	base := NewVM(parser.NewParser()).(*VM)
	base.AddNamespace("StartupPath", dir)
	first := NewRequestVM(base).(*RequestVM)
	if _, ctl := first.GetOrLoadClass("StartupPath\\Controller"); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if _, ok := first.GetClass("StartupPath\\Controller"); !ok {
		t.Fatal("startup path not inherited")
	}
	if _, ok := base.GetClass("StartupPath\\Controller"); ok {
		t.Fatal("request load through startup mapping leaked class")
	}
	first.AddNamespace("OnlyFirstRequest", dir)
	if _, ok := first.requestClassPaths().FindClassFile("OnlyFirstRequest\\Controller"); !ok {
		t.Fatal("local namespace unavailable")
	}
	second := NewRequestVM(base).(*RequestVM)
	if _, ok := second.requestClassPaths().FindClassFile("OnlyFirstRequest\\Controller"); ok {
		t.Fatal("explicit namespace addition leaked")
	}
}

func BenchmarkRequestOwnedParse(b *testing.B) {
	file := filepath.Join(b.TempDir(), "view.php")
	if err := os.WriteFile(file, []byte(`<?php $message = 'hello'; return $message;`), 0600); err != nil {
		b.Fatal(err)
	}
	base := NewVM(parser.NewParser()).(*VM)
	b.Run("new-request", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			request := NewRequestVM(base).(*RequestVM)
			if _, _, ctl := request.parseFileCached(file); ctl != nil {
				b.Fatal(ctl.AsString())
			}
		}
	})
	b.Run("same-request-cache", func(b *testing.B) {
		request := NewRequestVM(base).(*RequestVM)
		if _, _, ctl := request.parseFileCached(file); ctl != nil {
			b.Fatal(ctl.AsString())
		}
		b.ReportAllocs()
		for b.Loop() {
			if _, _, ctl := request.parseFileCached(file); ctl != nil {
				b.Fatal(ctl.AsString())
			}
		}
	})
}
