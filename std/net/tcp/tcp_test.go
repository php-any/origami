package tcp

import (
	"context"
	"io"
	"net"
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
	go func() {
		done <- serve(ctx, ln, map[string]int{"maxConnections": 1, "readTimeoutMs": 100, "writeTimeoutMs": 100, "connectionTimeoutMs": 1000}, func(ctx context.Context, c net.Conn) error {
			entered <- struct{}{}
			b := make([]byte, 1)
			_, e := c.Read(b)
			return e
		})
	}()
	c, e := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	<-entered
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
	next, e := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer next.Close()
	<-entered
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
