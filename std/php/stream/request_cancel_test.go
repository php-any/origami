package stream

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php/core"
)

type requestTestContext struct {
	data.Context
	request context.Context
}

func (c *requestTestContext) GoContext() context.Context { return c.request }

func TestCancelClosesBlockedStreamReads(t *testing.T) {
	for _, pipeReader := range []bool{false, true} {
		name := "file"
		if pipeReader {
			name = "reader"
		}
		t.Run(name, func(t *testing.T) {
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			defer writer.Close()
			request, cancel := context.WithCancel(context.Background())
			defer cancel()
			vm := runtime.NewVM(parser.NewParser())
			ctx := &requestTestContext{Context: runtime.NewContext(vm), request: request}
			var stream io.ReadCloser = NewStreamInfo(reader, "r")
			if pipeReader {
				stream = NewStreamInfoFromReader(reader, "r")
			}
			_ = core.NewResourceValue(core.NewResourceClass("stream", stream, int(fileDescriptor(reader))), ctx)
			done := make(chan error, 1)
			go func() { _, err := stream.Read(make([]byte, 1)); done <- err }()
			// Let Read enter the OS call. Closing the stream must be able to
			// interrupt it, rather than wait for the read's metadata lock.
			time.Sleep(20 * time.Millisecond)
			cancel()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("canceled read succeeded")
				}
			case <-time.After(time.Second):
				t.Fatal("request cancellation did not release blocked read")
			}
		})
	}
}

func TestStreamRegisteredAfterCancellationCloses(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	request, cancel := context.WithCancel(context.Background())
	cancel()
	ctx := &requestTestContext{Context: runtime.NewContext(runtime.NewVM(parser.NewParser())), request: request}
	stream := NewStreamInfo(reader, "r")
	_ = core.NewResourceValue(core.NewResourceClass("stream", stream, int(fileDescriptor(reader))), ctx)
	deadline := time.Now().Add(time.Second)
	for !stream.IsClosed() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if !stream.IsClosed() {
		t.Fatal("late resource escaped canceled request")
	}
}

func TestCanceledPHPReadUnwindsInsteadOfReturningEmpty(t *testing.T) {
	for _, function := range []data.FuncStmt{NewFreadFunction(), NewStreamGetContentsFunction()} {
		t.Run(function.GetName(), func(t *testing.T) {
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			defer writer.Close()
			request, cancel := context.WithCancel(context.Background())
			defer cancel()
			base := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
			ctx := &requestTestContext{Context: base, request: request}
			resource := core.NewResourceValue(core.NewResourceClass("stream", NewStreamInfo(reader, "r"), int(fileDescriptor(reader))), ctx)
			base.SetIndexZVal(0, data.NewZVal(resource))
			if function.GetName() == "fread" {
				base.SetIndexZVal(1, data.NewZVal(data.NewIntValue(1)))
			} else {
				base.SetIndexZVal(1, data.NewZVal(data.NewNullValue()))
				base.SetIndexZVal(2, data.NewZVal(data.NewIntValue(-1)))
			}
			done := make(chan any, 1)
			go func() {
				defer func() { done <- recover() }()
				_, _ = function.Call(ctx)
			}()
			time.Sleep(20 * time.Millisecond)
			cancel()
			select {
			case result := <-done:
				if !data.IsRequestCanceled(result) {
					t.Fatalf("PHP read cancellation = %v", result)
				}
			case <-time.After(time.Second):
				t.Fatal("PHP read remained blocked")
			}
		})
	}
}

func TestTemporaryStreamRemovedOnClose(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	fn := NewFopenFunction()
	ctx := runtime.NewContext(vm).CreateContext(fn.GetVariables())
	ctx.SetIndexZVal(0, data.NewZVal(data.NewStringValue("php://temp")))
	ctx.SetIndexZVal(1, data.NewZVal(data.NewStringValue("w+")))
	value, ctl := fn.Call(ctx)
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	stream := value.(*core.ResourceValue).GetResource().(*StreamInfo)
	path := stream.File.Name()
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("temporary stream retained: %v", err)
	}
}
