package redis

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/php-any/origami/data"
)

func TestRESPBinaryAndMalformedReplies(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if e := writeCommand(w, []string{"SET", "key\r\n", "a\x00b"}); e != nil {
		t.Fatal(e)
	}
	w.Flush()
	v, e := readReply(bufio.NewReader(&b))
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(v, []any{"SET", "key\r\n", "a\x00b"}) {
		t.Fatalf("binary roundtrip: %#v", v)
	}
	for _, raw := range []string{"$2147483647\r\n", "*2147483647\r\n", "$-2\r\n", ":abc\r\n", "$1\r\nx!!", "+OK\n", strings.Repeat("*1\r\n", 130) + "+OK\r\n"} {
		if _, e := readReply(bufio.NewReader(strings.NewReader(raw))); e == nil {
			t.Fatalf("malformed reply accepted: %.40q", raw)
		}
	}
	v, e = readReply(bufio.NewReader(strings.NewReader("*3\r\n$-1\r\n$0\r\n\r\n-WRONGTYPE example\r\n")))
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(v, []any{nil, "", serverError("WRONGTYPE example")}) {
		t.Fatalf("reply distinctions lost: %#v", v)
	}
}

func TestRequestCancellationClosesUncertainConnection(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()
	s := newState()
	s.conn = client
	s.reader = bufio.NewReader(client)
	defer s.disconnect()
	requestRead := make(chan struct{})
	go func() {
		_, _ = readReply(bufio.NewReader(server))
		close(requestRead)
		_, _ = io.Copy(io.Discard, server)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() { _, e := s.exchange(ctx, []string{"GET", "blocked"}); result <- e }()
	<-requestRead
	cancel()
	select {
	case e := <-result:
		if e == nil {
			t.Fatal("cancel swallowed")
		}
	case <-time.After(time.Second):
		t.Fatal("request cancellation did not interrupt I/O")
	}
	if s.conn != nil {
		t.Fatal("connection with uncertain protocol position was reused")
	}
}

func TestPrefixAndCommandArguments(t *testing.T) {
	str := data.NewStringValue
	opts := data.NewArrayValue(nil).(*data.ArrayValue)
	opts.AppendValue(str("NX"))
	opts.SetStringKey("px", data.NewIntValue(1000))
	p, e := prepare("set", []data.Value{str("key"), str("value"), opts}, "p:")
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(p.args, []string{"SET", "p:key", "value", "NX", "PX", "1000"}) {
		t.Fatal(p.args)
	}
	p, e = prepare("eval", []data.Value{str("return ARGV[1]"), data.NewArrayValue([]data.Value{str("key"), str("value")}), data.NewIntValue(1)}, "p:")
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(p.args, []string{"EVAL", "return ARGV[1]", "1", "p:key", "value"}) {
		t.Fatal(p.args)
	}
	p, e = prepare("lrem", []data.Value{str("list"), str("value"), data.NewIntValue(2)}, "p:")
	if e != nil || !reflect.DeepEqual(p.args, []string{"LREM", "p:list", "2", "value"}) {
		t.Fatal(p, e)
	}
	for _, cmd := range []string{"HELLO", "MULTI", "SUBSCRIBE", "CLIENT"} {
		if _, e = prepare("rawcommand", []data.Value{str(cmd)}, ""); e == nil {
			t.Fatalf("raw state change %s accepted", cmd)
		}
	}
}

func TestReadTimeoutDropsConnectionWithoutReplay(t *testing.T) {
	c, server := net.Pipe()
	defer server.Close()
	s := newState()
	s.conn = c
	s.reader = bufio.NewReader(c)
	s.readTimeout = 20 * time.Millisecond
	defer s.disconnect()
	go io.Copy(io.Discard, server)
	if _, err := s.exchange(context.Background(), []string{"INCR", "key"}); err == nil {
		t.Fatal("timeout swallowed")
	}
	if s.conn != nil {
		t.Fatal("timed out connection retained")
	}
}

func TestAlreadyCanceledExecDoesNotLeakTransactionState(t *testing.T) {
	c, server := net.Pipe()
	defer server.Close()
	s := newState()
	s.conn = c
	s.reader = bufio.NewReader(c)
	s.mode = multiMode
	defer s.disconnect()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.exchange(ctx, []string{"EXEC"}); err != context.Canceled {
		t.Fatalf("canceled EXEC: %v", err)
	}
	if s.conn != nil || s.mode != atomicMode {
		t.Fatal("canceled request retained connection/transaction state")
	}
}
