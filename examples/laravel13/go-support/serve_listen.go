//go:build !windows

package gosupport

import (
	"net"
)

func listenTCP(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
