package data

import (
	"context"
	"testing"
	"time"
)

func TestWaitRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan any, 1)
	go func() {
		defer func() { done <- recover() }()
		WaitRequest(ctx, 30*time.Second)
	}()
	cancel()
	select {
	case result := <-done:
		if !IsRequestCanceled(result) {
			t.Fatalf("panic = %v", result)
		}
	case <-time.After(time.Second):
		t.Fatal("wait remained blocked after cancel")
	}
	defer func() {
		if !IsRequestCanceled(recover()) {
			t.Fatal("zero wait ignored prior cancellation")
		}
	}()
	WaitRequest(ctx, 0)
}
