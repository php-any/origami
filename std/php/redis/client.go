package redis

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/php-any/origami/data"
	oruntime "github.com/php-any/origami/runtime"
)

type clientState struct {
	mu                                       sync.Mutex
	conn                                     net.Conn
	reader                                   *bufio.Reader
	prefix, lastError, host                  string
	port, database, mode                     int
	timeout, readTimeout                     time.Duration
	queue                                    []prepared
	queueBytes                               int
	serializer, maxRetries, backoffAlgorithm int
	backoffBase, backoffCap                  time.Duration
	scanRetry, scanPrefix                    bool
	persistent                               bool
	persistentID                             string
	lazy                                     bool
	closed                                   bool
	auth                                     []string
	watching                                 bool
	poolDirty                                bool
}

func newState() *clientState {
	s := &clientState{port: 6379, maxRetries: 10}
	runtime.SetFinalizer(s, func(s *clientState) { s.mu.Lock(); defer s.mu.Unlock(); s.release() })
	return s
}
func (s *clientState) disconnect() {
	if s.conn != nil {
		_ = s.conn.Close()
	}
	s.conn = nil
	s.reader = nil
	s.mode = atomicMode
	s.queue = nil
	s.queueBytes = 0
}

// Cancellation interrupts blocking I/O even when the caller has no deadline.
// Wait for a running callback before clearing deadlines for the next operation.
func (s *clientState) ioScope(ctx context.Context) (func(), error) {
	if s.conn == nil {
		return nil, fmt.Errorf("Redis server went away (not connected)")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	timeout := s.readTimeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	deadline := time.Now().Add(timeout)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	c := s.conn
	if err := c.SetDeadline(deadline); err != nil {
		return nil, err
	}
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() { _ = c.SetDeadline(time.Now()); close(done) })
	return func() {
		if !stop() {
			<-done
		}
		_ = c.SetDeadline(time.Time{})
	}, nil
}
func (s *clientState) exchange(ctx context.Context, args []string) (any, error) {
	end, err := s.ioScope(ctx)
	if err != nil {
		// The server may already be in MULTI/WATCH from an earlier call.
		// Do not retain that state when a canceled request cannot finish it.
		s.disconnect()
		return nil, err
	}
	defer end()
	w := bufio.NewWriter(s.conn)
	if err = writeCommand(w, args); err == nil {
		err = w.Flush()
	}
	if err != nil {
		s.disconnect()
		return nil, err
	}
	v, err := readReply(s.reader)
	if err != nil {
		s.disconnect()
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}
	return v, err
}
func (s *clientState) remember(p prepared, v data.Value) {
	if b, ok := v.(*data.BoolValue); ok && b.Value {
		switch p.args[0] {
		case "AUTH":
			s.auth = append([]string(nil), p.args...)
		case "WATCH":
			s.watching = true
		case "UNWATCH":
			s.watching = false
		}
	}
	if p.args[0] == "CLIENT" && len(p.args) > 1 && strings.EqualFold(p.args[1], "SETNAME") {
		s.poolDirty = true
	}
	if p.args[0] == "SELECT" {
		if b, ok := v.(*data.BoolValue); ok && b.Value {
			s.database, _ = strconv.Atoi(p.args[1])
		}
	}
}
func (s *clientState) enqueue(p prepared) error {
	size := 0
	for _, a := range p.args {
		size += len(a) + 32
	}
	if size > maxReplyBytes-s.queueBytes {
		return fmt.Errorf("Redis batch exceeds 64 MiB limit")
	}
	s.queue = append(s.queue, p)
	s.queueBytes += size
	return nil
}
func (s *clientState) execute(ctx data.Context, p prepared) (data.GetValue, data.Control) {
	if s.mode != pipelineMode {
		if err := s.ensureConnected(oruntime.RequestContext()); err != nil {
			return failure(err)
		}
	}
	if (p.args[0] == "WATCH" || p.args[0] == "UNWATCH" || p.args[0] == "AUTH") && s.mode != atomicMode {
		return bad("%s requires atomic mode", p.args[0])
	}
	if s.mode == pipelineMode {
		if err := s.enqueue(p); err != nil {
			return failure(err)
		}
		return instance(ctx), nil
	}
	if s.mode == multiMode {
		// Check capacity before sending; a locally rejected command must never be
		// queued by the server without a corresponding result decoder.
		if err := s.enqueue(p); err != nil {
			return failure(err)
		}
		v, err := s.exchange(oruntime.RequestContext(), p.args)
		if err != nil {
			return failure(err)
		}
		if v != status("QUEUED") {
			s.queue = s.queue[:len(s.queue)-1]
			s.queueBytes = 0
			for _, q := range s.queue {
				for _, a := range q.args {
					s.queueBytes += len(a) + 32
				}
			}
			d, e := s.decode(v, false)
			if e != nil {
				return failure(e)
			}
			return d, nil
		}
		return instance(ctx), nil
	}
	v, err := s.exchange(oruntime.RequestContext(), p.args)
	if err != nil {
		s.lastError = err.Error()
		return failure(err)
	}
	d, err := s.convertValue(ctx, p, v)
	if err != nil {
		return failure(err)
	}
	s.remember(p, d)
	return d, nil
}
func (s *clientState) exec(ctx data.Context) (data.GetValue, data.Control) {
	mode := s.mode
	queue := s.queue
	if mode == atomicMode {
		return data.NewBoolValue(false), nil
	}
	s.mode = atomicMode
	s.queue = nil
	s.queueBytes = 0
	var replies []any
	if mode == multiMode {
		s.watching = false
		v, err := s.exchange(oruntime.RequestContext(), []string{"EXEC"})
		if err != nil {
			return failure(err)
		}
		if v == nil {
			return data.NewBoolValue(false), nil
		}
		if _, ok := v.(serverError); ok {
			d, e := s.decode(v, false)
			if e != nil {
				return failure(e)
			}
			return d, nil
		}
		var ok bool
		replies, ok = v.([]any)
		if !ok {
			s.disconnect()
			return bad("invalid EXEC response")
		}
	} else {
		if err := s.ensureConnected(oruntime.RequestContext()); err != nil {
			return failure(err)
		}
		end, err := s.ioScope(oruntime.RequestContext())
		if err != nil {
			s.disconnect()
			return failure(err)
		}
		defer end()
		w := bufio.NewWriter(s.conn)
		for _, p := range queue {
			if err = writeCommand(w, p.args); err != nil {
				break
			}
		}
		if err == nil {
			err = w.Flush()
		}
		if err != nil {
			s.disconnect()
			return failure(err)
		}
		budget := maxReplyBytes
		for range queue {
			v, e := readValue(s.reader, &budget, 0)
			if e != nil {
				s.disconnect()
				return failure(e)
			}
			replies = append(replies, v)
		}
	}
	if len(replies) != len(queue) {
		s.disconnect()
		return bad("Redis batch response count mismatch")
	}
	out := make([]data.Value, len(queue))
	var firstError error
	for i, p := range queue {
		d, e := s.convertValue(ctx, p, replies[i])
		if e != nil {
			if firstError == nil {
				firstError = e
			}
			continue
		}
		s.remember(p, d)
		out[i] = d
	}
	if firstError != nil {
		return failure(firstError)
	}
	return data.NewArrayValue(out), nil
}

type localFunc func(data.Context, *clientState, []data.Value) (data.GetValue, data.Control)

func localMethods() map[string]localFunc {
	m := map[string]localFunc{}
	m["__construct"] = construct
	m["scan"] = scan
	m["connect"] = connect
	m["pconnect"] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
		return openConnection(ctx, s, a, true)
	}
	m["close"] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
		if e := arity("close", a, 0, 0); e != nil {
			return failure(e)
		}
		connected := s.conn != nil
		s.disconnect()
		s.closed = true
		s.lazy = false
		return data.NewBoolValue(connected), nil
	}
	m["setoption"] = setOption
	m["getoption"] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
		if e := arity("getOption", a, 1, 1); e != nil {
			return failure(e)
		}
		n, e := integer(a[0])
		if e != nil {
			return failure(e)
		}
		switch n {
		case optPrefix:
			return data.NewStringValue(s.prefix), nil
		case optReadTimeout:
			return data.NewFloatValue(s.readTimeout.Seconds()), nil
		case optSerializer:
			return data.NewIntValue(s.serializer), nil
		case optScan:
			v := 0
			if s.scanRetry {
				v |= 1
			}
			if s.scanPrefix {
				v |= 2
			}
			return data.NewIntValue(v), nil
		case optMaxRetries:
			return data.NewIntValue(s.maxRetries), nil
		case optBackoffAlgorithm:
			return data.NewIntValue(s.backoffAlgorithm), nil
		case optBackoffBase:
			return data.NewIntValue(int(s.backoffBase / time.Millisecond)), nil
		case optBackoffCap:
			return data.NewIntValue(int(s.backoffCap / time.Millisecond)), nil
		case optCompression:
			return data.NewIntValue(0), nil
		}
		return data.NewBoolValue(false), nil
	}
	m["multi"] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
		if e := arity("multi", a, 0, 1); e != nil {
			return failure(e)
		}
		mode := multiMode
		if len(a) > 0 {
			n, e := integer(a[0])
			if e != nil {
				return failure(e)
			}
			mode = n
		}
		return begin(ctx, s, mode)
	}
	m["pipeline"] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
		if e := arity("pipeline", a, 0, 0); e != nil {
			return failure(e)
		}
		return begin(ctx, s, pipelineMode)
	}
	m["exec"] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
		if e := arity("exec", a, 0, 0); e != nil {
			return failure(e)
		}
		return s.exec(ctx)
	}
	m["discard"] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
		if e := arity("discard", a, 0, 0); e != nil {
			return failure(e)
		}
		mode := s.mode
		s.watching = false
		s.mode = atomicMode
		s.queue = nil
		s.queueBytes = 0
		if mode == pipelineMode {
			return data.NewBoolValue(true), nil
		}
		if mode == atomicMode {
			return data.NewBoolValue(false), nil
		}
		v, e := s.exchange(oruntime.RequestContext(), []string{"DISCARD"})
		if e != nil {
			return failure(e)
		}
		d, e := s.convert(prepared{result: "bool"}, v)
		if e != nil {
			return failure(e)
		}
		return d, nil
	}
	for _, name := range []string{"isconnected", "gethost", "getport", "getdbnum", "getmode", "gettimeout", "getreadtimeout", "getpersistentid", "getlasterror", "clearlasterror"} {
		n := name
		m[n] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
			if e := arity(n, a, 0, 0); e != nil {
				return failure(e)
			}
			switch n {
			case "isconnected":
				return data.NewBoolValue(s.conn != nil), nil
			case "gethost":
				if s.host == "" {
					return data.NewBoolValue(false), nil
				}
				return data.NewStringValue(s.host), nil
			case "getport":
				return data.NewIntValue(s.port), nil
			case "getdbnum":
				return data.NewIntValue(s.database), nil
			case "getmode":
				return data.NewIntValue(s.mode), nil
			case "gettimeout":
				return data.NewFloatValue(s.timeout.Seconds()), nil
			case "getreadtimeout":
				return data.NewFloatValue(s.readTimeout.Seconds()), nil
			case "getpersistentid":
				if s.persistent {
					return data.NewStringValue(s.persistentID), nil
				}
				return data.NewNullValue(), nil
			case "getlasterror":
				if s.lastError == "" {
					return data.NewNullValue(), nil
				}
				return data.NewStringValue(s.lastError), nil
			case "clearlasterror":
				s.lastError = ""
				return data.NewBoolValue(true), nil
			}
			return data.NewNullValue(), nil
		}
	}
	for _, name := range []string{"_prefix", "_serialize", "_unserialize", "_pack", "_unpack"} {
		n := name
		m[n] = func(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
			if e := arity(n, a, 1, 1); e != nil {
				return failure(e)
			}
			if n != "_prefix" && s.serializer == 1 {
				if n == "_serialize" || n == "_pack" {
					v, err := s.pack(ctx, a[0])
					if err != nil {
						return failure(err)
					}
					return data.NewStringValue(v), nil
				}
				str, err := scalar(a[0])
				if err != nil {
					return failure(err)
				}
				v, err := s.unpack(ctx, data.NewStringValue(str))
				if err != nil {
					return failure(err)
				}
				return v, nil
			}
			v, e := scalar(a[0])
			if e != nil {
				return failure(e)
			}
			if n == "_prefix" {
				v = s.prefix + v
			}
			return data.NewStringValue(v), nil
		}
	}
	return m
}
func duration(v data.Value) (time.Duration, error) {
	var seconds float64
	switch n := v.(type) {
	case *data.IntValue:
		seconds = float64(n.Value)
	case *data.FloatValue:
		seconds = n.Value
	default:
		return 0, fmt.Errorf("timeout must be numeric")
	}
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 || seconds > float64(math.MaxInt64)/float64(time.Second) {
		return 0, fmt.Errorf("invalid timeout")
	}
	return time.Duration(seconds * float64(time.Second)), nil
}
func connect(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
	return openConnection(ctx, s, a, false)
}
func openConnection(ctx data.Context, s *clientState, a []data.Value, persistent bool) (data.GetValue, data.Control) {
	return configureConnection(ctx, s, a, persistent, false)
}
func configureConnection(ctx data.Context, s *clientState, a []data.Value, persistent, lazy bool) (data.GetValue, data.Control) {
	if e := arity("connect", a, 1, 7); e != nil {
		return failure(e)
	}
	host, e := scalar(a[0])
	if e != nil {
		return failure(e)
	}
	port := 6379
	var timeout, readTimeout time.Duration
	if len(a) > 1 {
		port, e = integer(a[1])
		if e != nil {
			return failure(e)
		}
	}
	if len(a) > 2 {
		timeout, e = duration(a[2])
		if e != nil {
			return failure(e)
		}
	}
	id := ""
	if len(a) > 3 && !null(a[3]) && persistent {
		id, e = scalar(a[3])
		if e != nil {
			return failure(e)
		}
	}
	retryInterval := 0
	if len(a) > 4 {
		n, e := integer(a[4])
		if e != nil || n < 0 || n > math.MaxInt32 {
			return bad("retry_interval must be a non-negative millisecond interval")
		}
		retryInterval = n
	}
	if len(a) > 5 {
		readTimeout, e = duration(a[5])
		if e != nil {
			return failure(e)
		}
	}
	var auth []string
	if len(a) > 6 && !null(a[6]) {
		v, ok := a[6].(*data.ArrayValue)
		if !ok {
			return bad("connection context must be an array")
		}
		for i, slot := range v.Range() {
			if slot.PHPArrayKey(i).AsString() != "auth" {
				return bad("custom stream contexts are not supported")
			}
			p, err := prepare("auth", []data.Value{slot.ReadValue()}, "")
			if err != nil {
				return failure(err)
			}
			auth = p.args
		}
	}
	network := "tcp"
	address := ""
	secure := strings.HasPrefix(host, "tls://")
	target := strings.TrimPrefix(strings.TrimPrefix(host, "tcp://"), "tls://")
	if strings.HasPrefix(target, "/") {
		network = "unix"
		address = target
	} else {
		if port < 1 || port > 65535 {
			return bad("port must be between 1 and 65535")
		}
		if target == "" || strings.ContainsAny(target, "\x00/\r\n") {
			return bad("invalid Redis host")
		}
		address = net.JoinHostPort(strings.Trim(target, "[]"), strconv.Itoa(port))
	}
	// Validate first, then release the old socket. Pool entries are exclusive leases.
	s.release()
	s.host, s.port, s.timeout, s.readTimeout = host, port, timeout, readTimeout
	s.database, s.auth, s.watching = 0, auth, false
	s.persistent, s.persistentID, s.closed, s.lazy = persistent, id, false, false
	s.poolDirty = false
	s.backoffBase, s.backoffCap = time.Duration(retryInterval)*time.Millisecond, time.Duration(retryInterval)*time.Millisecond
	if lazy {
		s.lazy = true
		return data.NewNullValue(), nil
	}
	if persistent && s.borrow() {
		return data.NewBoolValue(true), nil
	}
	if err := s.dial(oruntime.RequestContext(), network, address, secure, target); err != nil {
		s.lastError = err.Error()
		return failure(err)
	}
	return data.NewBoolValue(true), nil
}
func (s *clientState) dial(ctx context.Context, network, address string, secure bool, target string) error {
	dialTimeout := s.timeout
	if dialTimeout == 0 {
		dialTimeout = 60 * time.Second
	}
	dialer := &net.Dialer{Timeout: dialTimeout}
	// A failed reconnect must not leave commands targeting the previous server.
	s.disconnect()
	var c net.Conn
	var e error
	if secure {
		if network != "tcp" {
			return fmt.Errorf("TLS requires TCP")
		}
		d := tls.Dialer{NetDialer: dialer, Config: &tls.Config{MinVersion: tls.VersionTLS12, ServerName: strings.Trim(target, "[]")}}
		c, e = d.DialContext(ctx, network, address)
	} else {
		c, e = dialer.DialContext(ctx, network, address)
	}
	if e != nil {
		s.lastError = e.Error()
		return e
	}
	s.conn = c
	s.reader = bufio.NewReaderSize(c, 64<<10)
	s.lastError = ""
	if len(s.auth) > 0 {
		v, err := s.exchange(ctx, s.auth)
		if err != nil {
			return err
		}
		if v != status("OK") {
			s.disconnect()
			return fmt.Errorf("Redis authentication failed: %v", v)
		}
	}
	if s.database != 0 {
		v, err := s.exchange(ctx, []string{"SELECT", strconv.Itoa(s.database)})
		if err != nil {
			return err
		}
		if v != status("OK") {
			s.disconnect()
			return fmt.Errorf("Redis SELECT failed: %v", v)
		}
	}
	return nil
}
func setOption(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
	if e := arity("setOption", a, 2, 2); e != nil {
		return failure(e)
	}
	n, e := integer(a[0])
	if e != nil {
		return failure(e)
	}
	switch n {
	case optPrefix:
		v, e := scalar(a[1])
		if e != nil {
			return failure(e)
		}
		s.prefix = v
	case optReadTimeout:
		d, e := duration(a[1])
		if e != nil {
			return failure(e)
		}
		s.readTimeout = d
	case optSerializer:
		v, err := integer(a[1])
		if err != nil || (v != 0 && v != 1) {
			return bad("unsupported Redis serializer")
		}
		s.serializer = v
	case optScan:
		v, err := integer(a[1])
		if err != nil {
			return failure(err)
		}
		switch v {
		case 0:
			s.scanRetry = false
		case 1:
			s.scanRetry = true
		case 2:
			s.scanPrefix = true
		case 3:
			s.scanPrefix = false
		default:
			return data.NewBoolValue(false), nil
		}
	case optMaxRetries, optBackoffAlgorithm, optBackoffBase, optBackoffCap:
		v, err := integer(a[1])
		if err != nil || v < 0 {
			return data.NewBoolValue(false), nil
		}
		switch n {
		case optMaxRetries:
			s.maxRetries = v
		case optBackoffAlgorithm:
			if v > 6 {
				return data.NewBoolValue(false), nil
			}
			s.backoffAlgorithm = v
		case optBackoffBase, optBackoffCap:
			if int64(v) > math.MaxInt64/int64(time.Millisecond)/1024 {
				return data.NewBoolValue(false), nil
			}
			if n == optBackoffBase {
				s.backoffBase = time.Duration(v) * time.Millisecond
			} else {
				s.backoffCap = time.Duration(v) * time.Millisecond
			}
		}
	case optCompression:
		v, e := integer(a[1])
		if e != nil || v != 0 {
			return bad("only SERIALIZER_NONE and COMPRESSION_NONE are supported")
		}
	default:
		return bad("unsupported Redis option %d", n)
	}
	return data.NewBoolValue(true), nil
}
func begin(ctx data.Context, s *clientState, mode int) (data.GetValue, data.Control) {
	if mode != multiMode && mode != pipelineMode {
		return bad("mode must be Redis::MULTI or Redis::PIPELINE")
	}
	if err := s.ensureConnected(oruntime.RequestContext()); err != nil {
		return failure(err)
	}
	if s.mode != atomicMode {
		return bad("nested Redis batches are not supported")
	}
	if mode == multiMode {
		v, e := s.exchange(oruntime.RequestContext(), []string{"MULTI"})
		if e != nil {
			return failure(e)
		}
		if v != status("OK") {
			d, e := s.decode(v, false)
			if e != nil {
				return failure(e)
			}
			return d, nil
		}
	}
	s.mode = mode
	return instance(ctx), nil
}
