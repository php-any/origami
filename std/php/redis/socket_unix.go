//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly

package redis

import (
	"net"
	"syscall"
)

func socketIdle(conn net.Conn) bool {
	if wrapper, ok := conn.(interface{ NetConn() net.Conn }); ok {
		conn = wrapper.NetConn()
	}
	sc, ok := conn.(syscall.Conn)
	if !ok {
		return true
	}
	raw, err := sc.SyscallConn()
	if err != nil {
		return false
	}
	idle := false
	err = raw.Control(func(fd uintptr) {
		var b [1]byte
		_, _, e := syscall.Recvfrom(int(fd), b[:], syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
		idle = e == syscall.EAGAIN || e == syscall.EWOULDBLOCK
	})
	return err == nil && idle
}
