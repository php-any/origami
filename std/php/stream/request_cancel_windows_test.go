//go:build windows

package stream

import (
	"context"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php/core"
	"os"
	"testing"
)

func TestStreamSelectEmptySetsAndCancellation(t *testing.T) {
	function := NewStreamSelectFunction()
	base := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
	for i := 0; i < 3; i++ {
		base.SetIndexZVal(i, data.NewZVal(data.NewArrayValue(nil)))
	}
	base.SetIndexZVal(3, data.NewZVal(data.NewIntValue(30)))
	_, ctl := function.Call(base)
	if thrown, ok := ctl.(*data.ThrowValue); !ok || thrown.Name != "ValueError" {
		t.Fatalf("empty select = %v", ctl)
	}
	request, cancel := context.WithCancel(context.Background())
	cancel()
	defer func() {
		if reason := recover(); !data.IsRequestCanceled(reason) {
			t.Fatalf("cancel = %v", reason)
		}
	}()
	_, _ = function.Call(&requestTestContext{Context: base, request: request})
}

func TestWindowsSelectCountsDistinctHandleAcrossSets(t *testing.T) {
	function := NewStreamSelectFunction()
	base := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
	file, err := os.Open("stream_select.go")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	resource := core.NewResourceValue(core.NewResourceClass("stream", NewStreamInfo(file, "r"), int(fileDescriptor(file))), base)
	for i := 0; i < 3; i++ {
		base.SetIndexZVal(i, data.NewZVal(data.NewArrayValue([]data.Value{resource, resource})))
	}
	base.SetIndexZVal(3, data.NewZVal(data.NewIntValue(0)))
	value, ctl := function.Call(base)
	if result, ok := value.(*data.IntValue); ctl != nil || !ok || result.Value != 1 {
		t.Fatalf("select = %v %v", value, ctl)
	}
	for i := 0; i < 3; i++ {
		value, _ := base.GetIndexValue(i)
		if value.(*data.ArrayValue).Len() != 2 {
			t.Fatal("PHP Windows retains duplicate handle membership")
		}
	}
}
