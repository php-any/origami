package script

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestListenerSlotReleasedOnce(t *testing.T) {
	raw, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	ln := LimitListener(raw, 1)
	defer ln.Close()
	accepted := make(chan net.Conn, 1)
	go func() { c, _ := ln.Accept(); accepted <- c }()
	client, e := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer client.Close()
	server := <-accepted
	server.Close()
	server.Close()
	go func() { c, _ := ln.Accept(); accepted <- c }()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	next, e := (&net.Dialer{}).DialContext(ctx, "tcp", ln.Addr().String())
	if e != nil {
		t.Fatal(e)
	}
	defer next.Close()
	select {
	case c := <-accepted:
		c.Close()
	case <-ctx.Done():
		t.Fatal("slot leaked")
	}
}
