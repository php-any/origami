package data

import (
	"fmt"
	"os"
	"sync"
)

type PHPErrorInfo struct {
	Type          int
	Message, File string
	Line          int
}

// PHPErrorState is cold request state. It is never read by ordinary Call().
type PHPErrorState struct {
	mu        sync.Mutex
	reporting *int
	last      *PHPErrorInfo
	running   bool
}

func (s *PHPErrorState) Reporting() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reporting == nil {
		return 32767
	}
	return *s.reporting
}
func (s *PHPErrorState) SetReporting(level int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := 32767
	if s.reporting != nil {
		old = *s.reporting
	}
	s.reporting = &level
	return old
}
func (s *PHPErrorState) Last() *PHPErrorInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.last == nil {
		return nil
	}
	copy := *s.last
	return &copy
}
func (s *PHPErrorState) SetLast(info *PHPErrorInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = info
}
func (s *PHPErrorState) BeginHandler() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return false
	}
	s.running = true
	return true
}
func (s *PHPErrorState) EndHandler() { s.mu.Lock(); defer s.mu.Unlock(); s.running = false }
func (s *PHPErrorState) NewRequest() *PHPErrorState {
	copy := &PHPErrorState{}
	copy.SetReporting(s.Reporting())
	return copy
}

type PHPErrorHost interface{ PHPErrorState() *PHPErrorState }

func ErrorState(ctx Context) *PHPErrorState {
	if host, ok := ctx.GetVM().(PHPErrorHost); ok {
		return host.PHPErrorState()
	}
	return &PHPErrorState{}
}

type PHPErrorReporter interface {
	ReportPHPError(Context, int, string, From) Control
}

func EmitPHPError(ctx Context, level int, message string, from From) Control {
	if reporter, ok := ctx.GetVM().(PHPErrorReporter); ok {
		return reporter.ReportPHPError(ctx, level, message, from)
	}
	fmt.Fprintln(os.Stderr, message)
	return nil
}
