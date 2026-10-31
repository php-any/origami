package redis

import (
	"net"
	"syscall"
	"unsafe"
)

var wsaPoll = syscall.NewLazyDLL("ws2_32.dll").NewProc("WSAPoll")

// RESP2 idle sockets must have neither pending data nor an EOF/error. Polling
// with zero timeout avoids a PING round trip and does not require ACL +ping.
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
		poll := struct {
			socket          uintptr
			events, revents int16
		}{socket: fd, events: 0x0100}
		result, _, _ := wsaPoll.Call(uintptr(unsafe.Pointer(&poll)), 1, 0)
		idle = result == 0
	})
	return err == nil && idle
}
