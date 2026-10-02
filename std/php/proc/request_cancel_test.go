package proc

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

func TestProcessCancellationChild(t *testing.T) {
	if os.Args[len(os.Args)-1] != "origami-cancel-child" {
		return
	}
	time.Sleep(30 * time.Second)
}

func TestCancelProcOpenReapsChildAndPipes(t *testing.T) {
	function := NewProcOpenFunction()
	base := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
	base.SetIndexZVal(0, data.NewZVal(data.NewArrayValue([]data.Value{
		data.NewStringValue(os.Args[0]), data.NewStringValue("-test.run=^TestProcessCancellationChild$"), data.NewStringValue("--"), data.NewStringValue("origami-cancel-child"),
	})))
	base.SetIndexZVal(1, data.NewZVal(data.NewNullValue()))
	base.SetIndexZVal(2, data.NewZVal(data.NewArrayValue(nil)))
	request, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &requestTestContext{Context: base, request: request}
	result, control := function.Call(ctx)
	if control != nil {
		t.Fatal(control)
	}
	resource, ok := result.(*core.ResourceValue)
	if !ok {
		t.Fatalf("proc_open = %T", result)
	}
	process := resource.GetResource().(*ProcessInfo)
	defer process.Close()
	pipes, _ := base.GetIndexValue(2)
	stdout, _ := pipes.(*data.ArrayValue).FindSlotByIntKey(1)
	reader := stdout.Value.(*core.ResourceValue).GetResource().(io.ReadCloser)
	defer reader.Close()
	readDone := make(chan struct{})
	go func() { _, _ = io.ReadAll(reader); close(readDone) }()
	cancel()
	select {
	case <-process.done:
	case <-time.After(3 * time.Second):
		t.Fatal("canceled subprocess was not reaped")
	}
	select {
	case <-readDone:
	case <-time.After(time.Second):
		t.Fatal("canceled process pipe remained blocked")
	}
	if process.GetRunning() {
		t.Fatal("reaped process still reported running")
	}
}
