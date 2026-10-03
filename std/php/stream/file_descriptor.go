package stream

import (
	"syscall"
)

// Fd puts Unix Go files into blocking mode, preventing Close from interrupting
// a pending Read. RawConn.Control reads the descriptor without that side effect.
func fileDescriptor(file interface {
	SyscallConn() (syscall.RawConn, error)
}) uintptr {
	raw, err := file.SyscallConn()
	if err != nil {
		return ^uintptr(0)
	}
	fd := ^uintptr(0)
	if err := raw.Control(func(value uintptr) { fd = value }); err != nil {
		return ^uintptr(0)
	}
	return fd
}
