package redis

import (
	"bufio"
	"context"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/php-any/origami/data"
)

func construct(ctx data.Context, s *clientState, a []data.Value) (data.GetValue, data.Control) {
	if err := arity("__construct", a, 0, 1); err != nil {
		return failure(err)
	}
	if len(a) == 0 || null(a[0]) {
		return data.NewNullValue(), nil
	}
	opts, ok := a[0].(*data.ArrayValue)
	if !ok {
		return bad("Redis constructor options must be an array")
	}
	args := []data.Value{data.NewStringValue("127.0.0.1"), data.NewIntValue(6379), data.NewIntValue(0), data.NewNullValue(), data.NewIntValue(0), data.NewIntValue(0), data.NewArrayValue(nil)}
	persistent := false
	var backoff *data.ArrayValue
	for i, slot := range opts.Range() {
		v := slot.ReadValue()
		switch strings.ToLower(slot.PHPArrayKey(i).AsString()) {
		case "host":
			args[0] = v
		case "port":
			args[1] = v
		case "connecttimeout":
			args[2] = v
		case "retryinterval":
			args[4] = v
		case "readtimeout":
			args[5] = v
		case "persistent":
			if _, ok := v.(*data.StringValue); ok {
				persistent = true
				args[3] = v
			} else if b, ok := v.(data.AsBool); ok {
				persistent, _ = b.AsBool()
			}
		case "auth":
			args[6].(*data.ArrayValue).SetStringKey("auth", v)
		case "backoff":
			backoff, ok = v.(*data.ArrayValue)
			if !ok {
				return bad("backoff must be an array")
			}
		case "ssl":
			return bad("custom TLS context options are not supported")
		default:
			return bad("unknown Redis constructor option %s", slot.PHPArrayKey(i).AsString())
		}
	}
	if _, ctl := configureConnection(ctx, s, args, persistent, true); ctl != nil {
		return nil, ctl
	}
	if backoff != nil {
		for i, slot := range backoff.Range() {
			option, ok := map[string]int{"algorithm": optBackoffAlgorithm, "base": optBackoffBase, "cap": optBackoffCap}[strings.ToLower(slot.PHPArrayKey(i).AsString())]
			if !ok {
				return bad("unknown backoff option")
			}
			v, ctl := setOption(ctx, s, []data.Value{data.NewIntValue(option), slot.ReadValue()})
			if ctl != nil {
				return nil, ctl
			}
			if b, ok := v.(*data.BoolValue); !ok || !b.Value {
				return bad("invalid backoff option")
			}
		}
	}
	return data.NewNullValue(), nil
}

// Idle sockets are leased exclusively, never shared by live PHP objects. A
// bounded process pool spans request VMs without retaining their contexts.
type poolKey struct {
	host, id, auth       string
	port, db             int
	timeout, readTimeout time.Duration
}
type idleConnection struct {
	key    poolKey
	conn   net.Conn
	reader *bufio.Reader
	timer  *time.Timer
}

var persistentPool = struct {
	sync.Mutex
	entries []*idleConnection
}{}

func (s *clientState) poolKey() poolKey {
	// Length-prefix credentials so embedded delimiters cannot collide.
	var auth strings.Builder
	for _, value := range s.auth {
		fmt.Fprintf(&auth, "%d:%s", len(value), value)
	}
	return poolKey{s.host, s.persistentID, auth.String(), s.port, s.database, s.timeout, s.readTimeout}
}
func (s *clientState) release() {
	if s.conn == nil || !s.persistent || s.mode != atomicMode || s.watching || s.poolDirty || s.reader.Buffered() != 0 {
		s.disconnect()
		return
	}
	persistentPool.Lock()
	defer persistentPool.Unlock()
	if len(persistentPool.entries) >= 32 {
		s.disconnect()
		return
	}
	entry := &idleConnection{key: s.poolKey(), conn: s.conn, reader: s.reader}
	entry.timer = time.AfterFunc(time.Minute, func() {
		persistentPool.Lock()
		defer persistentPool.Unlock()
		for i, candidate := range persistentPool.entries {
			if candidate == entry {
				persistentPool.entries = append(persistentPool.entries[:i], persistentPool.entries[i+1:]...)
				entry.conn.Close()
				break
			}
		}
	})
	persistentPool.entries = append(persistentPool.entries, entry)
	s.conn, s.reader = nil, nil
}
func (s *clientState) borrow() bool {
	key := s.poolKey()
	persistentPool.Lock()
	defer persistentPool.Unlock()
	for i := len(persistentPool.entries) - 1; i >= 0; i-- {
		entry := persistentPool.entries[i]
		if entry.key != key {
			continue
		}
		persistentPool.entries = append(persistentPool.entries[:i], persistentPool.entries[i+1:]...)
		entry.timer.Stop()
		if !socketIdle(entry.conn) {
			entry.conn.Close()
			continue
		}
		s.conn, s.reader = entry.conn, entry.reader
		return true
	}
	return false
}
func (s *clientState) redial(ctx context.Context) error {
	target := strings.TrimPrefix(strings.TrimPrefix(s.host, "tcp://"), "tls://")
	network, address := "tcp", net.JoinHostPort(strings.Trim(target, "[]"), strconv.Itoa(s.port))
	if strings.HasPrefix(target, "/") {
		network, address = "unix", target
	}
	return s.dial(ctx, network, address, strings.HasPrefix(s.host, "tls://"), target)
}
func (s *clientState) ensureConnected(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		s.disconnect()
		return err
	}
	if s.closed || s.host == "" && s.conn == nil {
		return fmt.Errorf("Redis server went away (not connected)")
	}
	if s.conn != nil && (s.maxRetries == 0 || socketIdle(s.conn)) {
		return nil
	}
	if s.mode == multiMode || s.watching {
		s.disconnect()
		return fmt.Errorf("Connection lost in MULTI/WATCH mode")
	}
	if s.lazy {
		s.lazy = false
		if s.persistent && s.borrow() {
			return nil
		}
		return s.redial(ctx)
	}
	s.disconnect()
	var err error = fmt.Errorf("Redis connection lost")
	previous := time.Duration(0)
	for i := 0; i < s.maxRetries; i++ {
		delay := s.backoff(i, &previous)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		if err = s.redial(ctx); err == nil {
			return nil
		}
	}
	s.lastError = err.Error()
	return err
}
func (s *clientState) backoff(index int, previous *time.Duration) time.Duration {
	random := func(lo, hi time.Duration) time.Duration {
		if lo > hi {
			lo, hi = hi, lo
		}
		if lo == hi {
			return lo
		}
		return lo + time.Duration(rand.Int64N(int64(hi-lo)+1))
	}
	base, cap := s.backoffBase, s.backoffCap
	exponential := min(cap, base*time.Duration(1<<min(index, 10)))
	switch s.backoffAlgorithm {
	case 1:
		*previous = random(base, min(*previous, time.Duration(1<<61))*3)
		return min(cap, *previous)
	case 2:
		return random(0, exponential)
	case 3:
		return exponential/2 + random(0, exponential)/2
	case 4:
		return exponential
	case 5:
		return min(cap, random(0, base))
	case 6:
		return min(cap, base)
	default:
		if index == 0 {
			return min(cap, random(0, base))
		}
		return min(cap, base)
	}
}
