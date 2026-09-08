//go:build windows

package stream

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func applyFlock(file *os.File, op int) (bool, error) {
	handle := windows.Handle(file.Fd())
	ol := new(windows.Overlapped)
	nb := op&LockNB != 0
	base := op &^ LockNB

	switch base {
	case LockUN:
		err := windows.UnlockFileEx(handle, 0, 0xffffffff, 0xffffffff, ol)
		if err != nil {
			return false, err
		}
		return false, nil
	case LockSH, LockEX:
		flags := uint32(0)
		if base == LockEX {
			flags |= windows.LOCKFILE_EXCLUSIVE_LOCK
		}
		if nb {
			flags |= windows.LOCKFILE_FAIL_IMMEDIATELY
		}
		err := windows.LockFileEx(handle, flags, 0, 0xffffffff, 0xffffffff, ol)
		if err != nil {
			if nb && (err == windows.ERROR_LOCK_VIOLATION || err == windows.ERROR_IO_PENDING || err == syscall.EWOULDBLOCK) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	default:
		return false, syscall.EINVAL
	}
}
