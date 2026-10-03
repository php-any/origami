//go:build !windows

package stream

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php/core"
)

func unixSelectContext(t *testing.T, file *os.File) (data.FuncStmt, data.Context) {
	t.Helper()
	function := NewStreamSelectFunction()
	ctx := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
	resource := core.NewResourceValue(core.NewResourceClass("stream", NewStreamInfo(file, "r"), int(fileDescriptor(file))), ctx)
	read := data.NewArrayValue(nil).(*data.ArrayValue)
	read.SetKey(data.NewStringValue("pipe"), resource)
	ctx.SetIndexZVal(0, data.NewZVal(read))
	ctx.SetIndexZVal(1, data.NewZVal(data.NewArrayValue(nil)))
	ctx.SetIndexZVal(2, data.NewZVal(data.NewArrayValue(nil)))
	ctx.SetIndexZVal(3, data.NewZVal(data.NewIntValue(0)))
	ctx.SetIndexZVal(4, data.NewZVal(data.NewIntValue(0)))
	return function, ctx
}

func TestUnixSelectPipeReadinessPreservesKeys(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	function, ctx := unixSelectContext(t, reader)
	value, ctl := function.Call(ctx)
	if ctl != nil || value.(*data.IntValue).Value != 0 {
		t.Fatalf("empty pipe select = %v %v", value, ctl)
	}
	read, _ := ctx.GetIndexValue(0)
	if read.(*data.ArrayValue).Len() != 0 {
		t.Fatal("timeout retained an unreadable stream")
	}
	if _, err := writer.Write([]byte("ready")); err != nil {
		t.Fatal(err)
	}
	function, ctx = unixSelectContext(t, reader)
	value, ctl = function.Call(ctx)
	if ctl != nil || value.(*data.IntValue).Value != 1 {
		t.Fatalf("ready pipe select = %v %v", value, ctl)
	}
	read, _ = ctx.GetIndexValue(0)
	array := read.(*data.ArrayValue)
	if array.Len() != 1 || array.View().At(0).Name != "pipe" {
		t.Fatal("select lost the input key")
	}
}

func TestUnixSelectWaitingPipeCancellation(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	function, ctx := unixSelectContext(t, reader)
	ctx.SetIndexZVal(3, data.NewZVal(data.NewNullValue()))
	request, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	result := make(chan any, 1)
	go func() {
		defer func() { result <- recover() }()
		_, _ = function.Call(&requestTestContext{Context: ctx, request: request})
	}()
	select {
	case reason := <-result:
		if !data.IsRequestCanceled(reason) {
			t.Fatalf("select cancellation = %v", reason)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("select did not observe cancellation")
	}
}
