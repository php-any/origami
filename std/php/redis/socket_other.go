//go:build !windows && !linux && !darwin && !freebsd && !openbsd && !netbsd && !dragonfly

package redis

import "net"

// Other targets retain ordinary I/O error handling without a speculative read.
func socketIdle(net.Conn) bool { return true }
