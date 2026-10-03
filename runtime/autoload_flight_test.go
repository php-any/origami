package runtime

import (
	"context"
	"testing"
	"time"

	"github.com/php-any/origami/data"
)

func TestAutoloadFlightRecursionRetryAndCancellation(t *testing.T) {
	var flights autoloadFlights
	id := data.Symbols.Intern("FlightTest")
	leader, finish := flights.begin(context.Background(), id)
	if !leader {
		t.Fatal("first caller is not leader")
	}
	if recursive, _ := flights.begin(context.Background(), id); recursive {
		t.Fatal("recursive autoload entered")
	}
	canceled := make(chan any, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer func() { canceled <- recover() }()
		flights.begin(ctx, id)
	}()
	cancel()
	select {
	case result := <-canceled:
		if result != data.ErrRequestCanceled {
			t.Fatalf("cancellation = %v", result)
		}
	case <-time.After(time.Second):
		t.Fatal("waiter did not cancel")
	}
	waited := make(chan bool, 1)
	go func() { leader, _ := flights.begin(context.Background(), id); waited <- leader }()
	// Wait until the waiter has entered the loading branch.
	select {
	case <-waited:
		t.Fatal("loading flight did not wait")
	case <-time.After(10 * time.Millisecond):
	}
	finish(false)
	finish(true) // completion is idempotent, including deferred cleanup.
	select {
	case leader := <-waited:
		if leader {
			t.Fatal("waiter repeated the leader callback")
		}
	case <-time.After(time.Second):
		t.Fatal("waiter remained blocked")
	}
	if retry, complete := flights.begin(context.Background(), id); !retry {
		t.Fatal("failed load cannot retry")
	} else {
		complete(true)
	}
}
