//go:build windows

package stream

import (
	"context"
	"testing"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestCancelStreamSelectTimeout(t *testing.T) {
	function := NewStreamSelectFunction()
	base := runtime.NewContext(runtime.NewVM(parser.NewParser())).CreateContext(function.GetVariables())
	for i := 0; i < 3; i++ {
		base.SetIndexZVal(i, data.NewZVal(data.NewArrayValue(nil)))
	}
	base.SetIndexZVal(3, data.NewZVal(data.NewIntValue(30)))
	request, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &requestTestContext{Context: base, request: request}
	done := make(chan any, 1)
	go func() {
		defer func() { done <- recover() }()
		_, _ = function.Call(ctx)
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case reason := <-done:
		if reason != data.ErrRequestCanceled {
			t.Fatalf("cancel result = %v", reason)
		}
	case <-time.After(time.Second):
		t.Fatal("stream_select kept waiting after cancel")
	}
}
