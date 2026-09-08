//go:build windows

package stream

import "syscall"

func setStreamNonblock(fd uintptr, nonblocking bool) error {
	return syscall.SetNonblock(syscall.Handle(fd), nonblocking)
}
