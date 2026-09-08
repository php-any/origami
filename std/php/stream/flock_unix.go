//go:build !windows

package stream

import (
	"os"
	"syscall"
)

func phpFlockHow(op int) (int, bool) {
	nb := op&LockNB != 0
	base := op &^ LockNB
	var how int
	switch base {
	case LockSH:
		how = syscall.LOCK_SH
	case LockEX:
		how = syscall.LOCK_EX
	case LockUN:
		how = syscall.LOCK_UN
	default:
		return 0, false
	}
	if nb {
		how |= syscall.LOCK_NB
	}
	return how, true
}

func applyFlock(file *os.File, op int) (bool, error) {
	how, ok := phpFlockHow(op)
	if !ok {
		return false, syscall.EINVAL
	}
	err := syscall.Flock(int(file.Fd()), how)
	if err != nil {
		if op&LockNB != 0 && (err == syscall.EWOULDBLOCK || err == syscall.EAGAIN) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
