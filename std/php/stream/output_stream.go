package stream

import (
	"github.com/php-any/origami/data"
	"sync"
)

// php://output addresses the current PHP SAPI's output buffers. It owns no
// process stdout descriptor and closing it must not close the response writer.
type outputStream struct {
	mu     sync.RWMutex
	closed bool
}

func (s *outputStream) IsClosed() bool { s.mu.RLock(); defer s.mu.RUnlock(); return s.closed }
func (s *outputStream) Close() error   { s.mu.Lock(); s.closed = true; s.mu.Unlock(); return nil }
func init() {
	data.RegisterNativeRequestPolicy[*outputStream](data.NativeNew, func(*data.RequestObjectScope, *outputStream) *outputStream { return &outputStream{} })
}
