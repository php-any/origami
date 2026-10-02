package spl

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/stream"
)

type requestTestContext struct {
	data.Context
	request context.Context
}

func (c *requestTestContext) GoContext() context.Context { return c.request }

func TestCancelSPLBlockedRead(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	request, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &requestTestContext{Context: runtime.NewContext(runtime.NewVM(parser.NewParser())), request: request}
	state := &sfoStateValue{stream: stream.NewStreamInfo(reader, "r"), request: request}
	core.BindOwnedResource(ctx, state)
	done := make(chan any, 1)
	go func() { defer func() { done <- recover() }(); _, _ = sfoReadRawLine(state) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case result := <-done:
		if !data.IsRequestCanceled(result) {
			t.Fatalf("SPL read cancel = %v", result)
		}
	case <-time.After(time.Second):
		t.Fatal("SPL read remained blocked")
	}
}

func TestSPLReadAheadReturnsReadError(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()
	_ = reader.Close()
	state := &sfoStateValue{stream: stream.NewStreamInfo(reader, "r")}
	if err := sfoLoadAllLines(state); err == nil {
		t.Fatal("read error was swallowed")
	}
	state.flags = SFO_SKIP_EMPTY
	sfoReadNext(state)
	if state.valid {
		t.Fatal("failed read produced valid line")
	}
}

func TestCancelSPLTempClosesAndRemovesFile(t *testing.T) {
	request, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &requestTestContext{Context: runtime.NewContext(runtime.NewVM(parser.NewParser())), request: request}
	value := data.NewClassValue(NewSplTempFileObjectClass(), ctx)
	if err := sfoOpenFileForTemp(ctx, value, "php://temp", "w+b", 0); err != nil {
		t.Fatal(err)
	}
	state := sfoGetState(value)
	path := state.temporaryPath
	if path == "" {
		t.Fatal("missing owned temporary file")
	}
	defer os.Remove(path)
	cancel()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); os.IsNotExist(err) && state.stream.IsClosed() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("SPL temporary file survived cancellation")
}

func TestCancelSPLPreservesUserFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(path, []byte("line\n"), 0600); err != nil {
		t.Fatal(err)
	}
	request, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &requestTestContext{Context: runtime.NewContext(runtime.NewVM(parser.NewParser())), request: request}
	value := data.NewClassValue(NewSplFileObjectClass(), ctx)
	if err := sfoOpenFile(ctx, value, path, "r", 0, false); err != nil {
		t.Fatal(err)
	}
	state := sfoGetState(value)
	cancel()
	deadline := time.Now().Add(time.Second)
	for !state.stream.IsClosed() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if !state.stream.IsClosed() {
		t.Fatal("SPL handle survived cancellation")
	}
	if content, err := os.ReadFile(path); err != nil || string(content) != "line\n" {
		t.Fatalf("user file changed: %q, %v", content, err)
	}
}
