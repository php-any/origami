package runtime

import (
	"context"
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInvalidatedParseCannotRefillCache(t *testing.T) {
	cache := &parsedFileCache{}
	file := "changed.php"
	started, release := make(chan struct{}), make(chan struct{})
	completed := make(chan struct{})
	go func() {
		defer close(completed)
		_, _ = cache.load(context.Background(), file, func() (*parsedPHPFile, data.Control) {
			close(started)
			<-release
			return &parsedPHPFile{program: data.NewIntValue(1)}, nil
		})
	}()
	<-started
	flight, _ := syncMapLoad[*parsedFileFlight](&cache.flights, file)
	flight.invalidated.Store(true)
	cache.entries.Delete(file)
	close(release)
	<-completed
	if _, ok := cache.entries.Load(file); ok {
		t.Fatal("stale parse refilled the cache")
	}
	entry, ctl := cache.load(context.Background(), file, func() (*parsedPHPFile, data.Control) {
		return &parsedPHPFile{program: data.NewIntValue(2)}, nil
	})
	if ctl != nil || entry.program.(*data.IntValue).Value != 2 {
		t.Fatal("next load did not reparse")
	}
}

func TestIncludeDetectsExternalSourceReplacement(t *testing.T) {
	file := filepath.Join(t.TempDir(), "source.php")
	if err := os.WriteFile(file, []byte("<?php return 1;"), 0600); err != nil {
		t.Fatal(err)
	}
	p := parser.NewParser()
	vm := NewVM(p)
	program, ctl := p.ParseString(fmt.Sprintf("return require %q;", filepath.ToSlash(file)), "include.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	ctx := vm.CreateContext(p.GetVariables())
	for _, expected := range []int{1, 2} {
		if expected == 2 {
			if err := os.WriteFile(file, []byte("<?php return 2;"), 0600); err != nil {
				t.Fatal(err)
			}
			stamp := time.Now().Add(time.Second)
			if err := os.Chtimes(file, stamp, stamp); err != nil {
				t.Fatal(err)
			}
		}
		value, ctl := program.GetValue(ctx)
		if ctl != nil || value.(*data.IntValue).Value != expected {
			t.Fatalf("include got %v / %v, want %d", value, ctl, expected)
		}
	}
}
