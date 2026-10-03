//go:build !windows

package proc

import (
	"errors"
	"syscall"
	"testing"
	"time"
)

func watchDescendant(t *testing.T, pid int) func() bool {
	return func() bool {
		deadline := time.Now().Add(3 * time.Second)
		for time.Now().Before(deadline) {
			if errors.Is(syscall.Kill(pid, 0), syscall.ESRCH) {
				return true
			}
			time.Sleep(10 * time.Millisecond)
		}
		return false
	}
}
