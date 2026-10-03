package runtime

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestIgnoreAbortPreservesDeadlineAndRestoresClientCancellation(t *testing.T) {
	output := BeginRequestOutput()
	defer output()
	client, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	deadline := BeginRequestDeadline(client)
	defer deadline()
	SetRequestIgnoreUserAbort(true)
	ignored := RequestContext()
	cancel()
	if ignored.Err() != nil || RequestDeadlineExceeded() {
		t.Fatal("client cancellation was not ignored")
	}
	if _, ok := ignored.Deadline(); !ok {
		t.Fatal("execution deadline removed")
	}
	select {
	case <-ignored.Done():
	case <-time.After(time.Second):
		t.Fatal("execution deadline never fired")
	}
	if !errors.Is(ignored.Err(), context.DeadlineExceeded) {
		t.Fatal(ignored.Err())
	}
	SetRequestIgnoreUserAbort(false)
	if !errors.Is(RequestContext().Err(), context.DeadlineExceeded) {
		t.Fatal("restoring abort policy cleared an expired execution deadline")
	}
}

func TestIgnoreAbortKeepsPreviouslyBoundContext(t *testing.T) {
	output := BeginRequestOutput()
	defer output()
	client, cancel := context.WithCancel(context.Background())
	defer cancel()
	deadline := BeginRequestDeadline(client)
	defer deadline()
	captured := RequestContext()
	SetRequestIgnoreUserAbort(true)
	cancel()
	time.Sleep(10 * time.Millisecond)
	if captured.Err() != nil || captured != RequestContext() {
		t.Fatal("resources retained obsolete cancellation policy")
	}
	SetRequestIgnoreUserAbort(false)
	if captured.Err() != context.Canceled {
		t.Fatal("original context failed to observe restored abort policy")
	}
}
