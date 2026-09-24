//go:build !windows

package serve

import (
	"net"
)

func listenTCP(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
