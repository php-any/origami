// Package tcp provides bounded stream servers without imposing a wire protocol.
package tcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/net/internal/script"
	"github.com/php-any/origami/utils"
)

func Load(vm data.VM) { vm.AddClass(NewServer()); vm.AddClass(connection(nil)) }
func NewServer() data.ClassStmt {
	address := ""
	opts := map[string]int{"maxConnections": 128, "readTimeoutMs": 10000, "writeTimeoutMs": 10000, "connectionTimeoutMs": 60000}
	c := &script.Class{Name: "Net\\Tcp\\Server", Factory: NewServer}
	c.Methods = []data.Method{
		&script.Method{Name: "__construct", Args: []string{"address"}, Run: func(ctx data.Context) (data.GetValue, error) {
			var e error
			address, e = utils.ConvertFromIndex[string](ctx, 0)
			return nil, e
		}},
		&script.Method{Name: "configure", Args: []string{"options"}, Run: func(ctx data.Context) (data.GetValue, error) {
			v, e := script.Options(ctx, 0, opts)
			if e == nil {
				opts = v
			}
			return nil, e
		}},
		&script.Method{Name: "serve", Args: []string{"callback"}, Run: func(ctx data.Context) (data.GetValue, error) {
			fn, e := script.Callback(ctx, 0)
			if e != nil {
				return nil, e
			}
			ln, e := net.Listen("tcp", address)
			if e != nil {
				return nil, e
			}
			defer ln.Close()
			return nil, serve(ctx.GoContext(), ln, opts, func(gc context.Context, n net.Conn) error {
				return script.Invoke(ctx.GetVM(), fn, gc, func(c data.Context) []data.Value { return []data.Value{data.NewProxyValue(connection(n), c)} })
			})
		}},
	}
	return c
}
func serve(ctx context.Context, ln net.Listener, opts map[string]int, callback func(context.Context, net.Conn) error) error {
	ctx, cancel := context.WithCancel(ctx)
	slots := make(chan struct{}, opts["maxConnections"])
	var wg sync.WaitGroup
	stop := context.AfterFunc(ctx, func() { ln.Close() })
	defer stop()
	defer wg.Wait()
	defer cancel() // Unblock active connections before waiting after an Accept failure.
	for {
		conn, e := ln.Accept()
		if e != nil {
			if ctx.Err() != nil {
				return nil
			}
			return e
		}
		select {
		case slots <- struct{}{}:
		default:
			conn.Close()
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			defer conn.Close()
			gc, cancel := context.WithTimeout(ctx, time.Duration(opts["connectionTimeoutMs"])*time.Millisecond)
			defer cancel()
			stop := context.AfterFunc(gc, func() { conn.Close() })
			defer stop()
			wrapped := &timedConn{Conn: conn, read: time.Duration(opts["readTimeoutMs"]) * time.Millisecond, write: time.Duration(opts["writeTimeoutMs"]) * time.Millisecond}
			// Script owns application error reporting; an uncaught error closes this stream.
			_ = callback(gc, wrapped)
		}()
	}
}

type timedConn struct {
	net.Conn
	read, write time.Duration
}

func (c *timedConn) Read(p []byte) (int, error) {
	if e := c.Conn.SetReadDeadline(time.Now().Add(c.read)); e != nil {
		return 0, e
	}
	return c.Conn.Read(p)
}
func (c *timedConn) Write(p []byte) (int, error) {
	if e := c.Conn.SetWriteDeadline(time.Now().Add(c.write)); e != nil {
		return 0, e
	}
	return c.Conn.Write(p)
}
func connection(n net.Conn) data.ClassStmt {
	c := &script.Class{Name: "Net\\Tcp\\Conn", Source: n}
	c.Methods = []data.Method{
		&script.Method{Name: "read", Args: []string{"maxBytes"}, Run: func(ctx data.Context) (data.GetValue, error) {
			if n == nil {
				return nil, fmt.Errorf("uninitialized connection")
			}
			size, e := utils.ConvertFromIndex[int](ctx, 0)
			if e != nil {
				return nil, e
			}
			if size < 1 || size > 16<<20 {
				return nil, fmt.Errorf("read size must be 1..16777216")
			}
			b := make([]byte, size)
			k, e := n.Read(b)
			if errors.Is(e, io.EOF) {
				e = nil
			}
			return data.NewStringValue(string(b[:k])), e
		}},
		&script.Method{Name: "write", Args: []string{"bytes"}, Run: func(ctx data.Context) (data.GetValue, error) {
			if n == nil {
				return nil, fmt.Errorf("uninitialized connection")
			}
			b, e := utils.ConvertFromIndex[string](ctx, 0)
			if e != nil {
				return nil, e
			}
			total := 0
			for total < len(b) {
				k, e := io.WriteString(n, b[total:])
				total += k
				if e != nil {
					return nil, e
				}
				if k == 0 {
					return nil, io.ErrShortWrite
				}
			}
			return data.NewIntValue(total), nil
		}},
		&script.Method{Name: "close", Run: func(data.Context) (data.GetValue, error) {
			if n == nil {
				return nil, nil
			}
			return nil, n.Close()
		}},
	}
	return c
}
