//go:build amd64 || arm64

package getg

import "unsafe"

func G() unsafe.Pointer
