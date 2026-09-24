package serve

import (
	"context"
	"net"
	"syscall"
)

// Windows 默认 SO_REUSEADDR 允许同一端口被再次 bind；必须改成独占，占用时直接失败。
const soExclusiveAddrUse = ^syscall.SO_REUSEADDR

func listenTCP(addr string) (net.Listener, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var sockErr error
			if err := c.Control(func(fd uintptr) {
				handle := syscall.Handle(fd)
				sockErr = syscall.SetsockoptInt(handle, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 0)
				if sockErr != nil {
					return
				}
				sockErr = syscall.SetsockoptInt(handle, syscall.SOL_SOCKET, soExclusiveAddrUse, 1)
			}); err != nil {
				return err
			}
			return sockErr
		},
	}
	return lc.Listen(context.Background(), "tcp", addr)
}
