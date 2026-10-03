package runtime

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// One stable context per execution phase lets resources registered before
// ignore_user_abort changes observe the final policy. PHP call frames stay flat.
type requestCancellation struct {
	parent context.Context
	done   chan struct{}
	once   sync.Once
	ignore atomic.Bool
	cause  atomic.Int32
	stop   func() bool
	timer  *time.Timer
}

func newRequestCancellation(parent context.Context, ignore bool) *requestCancellation {
	state := &requestCancellation{parent: parent, done: make(chan struct{})}
	state.ignore.Store(ignore)
	if deadline, ok := parent.Deadline(); ok {
		state.timer = time.AfterFunc(time.Until(deadline), func() { state.cancel(context.DeadlineExceeded) })
	}
	if parent.Done() != nil {
		state.stop = context.AfterFunc(parent, func() {
			if !state.ignore.Load() || parent.Err() == context.DeadlineExceeded {
				state.cancel(parent.Err())
			}
		})
	}
	return state
}
func (c *requestCancellation) Deadline() (time.Time, bool) { return c.parent.Deadline() }
func (c *requestCancellation) Done() <-chan struct{}       { return c.done }
func (c *requestCancellation) Value(key any) any           { return c.parent.Value(key) }
func (c *requestCancellation) Err() error {
	if c.cause.Load() == 0 {
		select {
		case <-c.parent.Done():
			if !c.ignore.Load() || c.parent.Err() == context.DeadlineExceeded {
				c.cancel(c.parent.Err())
			}
		default:
		}
	}
	switch c.cause.Load() {
	case 1:
		return context.Canceled
	case 2:
		return context.DeadlineExceeded
	}
	return nil
}
func (c *requestCancellation) cancel(err error) {
	c.once.Do(func() {
		code := int32(1)
		if err == context.DeadlineExceeded {
			code = 2
		}
		c.cause.Store(code)
		close(c.done)
	})
}
func (c *requestCancellation) setIgnore(ignore bool) {
	c.ignore.Store(ignore)
	if err := c.parent.Err(); err != nil && (!ignore || err == context.DeadlineExceeded) {
		c.cancel(err)
	}
}
func (c *requestCancellation) close() {
	if c.stop != nil {
		c.stop()
	}
	if c.timer != nil {
		c.timer.Stop()
	}
	c.cancel(context.Canceled)
}
