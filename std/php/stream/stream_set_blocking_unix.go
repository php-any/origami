//go:build !windows

package stream

import "syscall"

func setStreamNonblock(fd uintptr, nonblocking bool) error {
	return syscall.SetNonblock(int(fd), nonblocking)
}
