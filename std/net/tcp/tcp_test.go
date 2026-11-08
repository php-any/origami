package tcp

import (
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func TestReadTimeoutAndCancellationReleaseConnection(t *testing.T) {
	ln, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer ln.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	entered := make(chan struct{}, 2)
	firstHandlerDone := make(chan struct{})
	var recordFirst sync.Once
	go func() {
		done <- serve(ctx, ln, map[string]int{"maxConnections": 1, "readTimeoutMs": 100, "writeTimeoutMs": 100, "connectionTimeoutMs": 1000}, func(ctx context.Context, c net.Conn) error {
			entered <- struct{}{}
			releaseSlot := false
			recordFirst.Do(func() { releaseSlot = true })
			if releaseSlot {
				defer close(firstHandlerDone)
			}
			b := make([]byte, 1)
			_, e := c.Read(b)
			return e
		})
	}()
	waitEntered := func(label string) {
		t.Helper()
		select {
		case <-entered:
		case <-time.After(3 * time.Second):
			t.Fatalf("%s: callback did not start", label)
		}
	}
	c, e := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	waitEntered("first connection")
	other, e := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer other.Close()
	other.SetReadDeadline(time.Now().Add(time.Second))
	if _, e = other.Read(make([]byte, 1)); e != io.EOF {
		t.Fatalf("excess connection not closed: %v", e)
	}
	c.SetReadDeadline(time.Now().Add(time.Second))
	if _, e = c.Read(make([]byte, 1)); e != io.EOF {
		t.Fatalf("idle connection not closed: %v", e)
	}
	select {
	case <-firstHandlerDone:
	case <-time.After(3 * time.Second):
		t.Fatal("first handler did not finish; maxConnections slot still held")
	}
	_ = c.Close()
	var next net.Conn
	replacementBy := time.Now().Add(5 * time.Second)
	for next == nil && time.Now().Before(replacementBy) {
		conn, dialErr := net.DialTimeout("tcp", ln.Addr().String(), 250*time.Millisecond)
		if dialErr != nil {
			continue
		}
		select {
		case <-entered:
			next = conn
		case <-time.After(400 * time.Millisecond):
			conn.Close()
		}
	}
	if next == nil {
		t.Fatal("replacement connection: callback did not start")
	}
	defer next.Close()
	cancel()
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown blocked")
	}
}
func TestWriteDeadline(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	c := &timedConn{Conn: a, write: 20 * time.Millisecond}
	_, e := c.Write([]byte("blocked"))
	if n, ok := e.(net.Error); !ok || !n.Timeout() {
		t.Fatalf("expected timeout: %v", e)
	}
}

func TestListenerFailureCancelsActiveCallbacks(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	entered := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- serve(ctx, listener, map[string]int{"maxConnections": 1, "connectionTimeoutMs": 30000}, func(ctx context.Context, conn net.Conn) error {
			close(entered)
			<-ctx.Done()
			return ctx.Err()
		})
	}()
	conn, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("callback did not start")
	}
	listener.Close()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("listener error was hidden")
		}
	case <-time.After(time.Second):
		t.Fatal("listener failure did not cancel active callbacks")
	}
}
