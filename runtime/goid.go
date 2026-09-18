package runtime

import (
	"runtime"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/php-any/origami/runtime/getg"
)

// goidOff 是 runtime.g 里 goid 字段的字节偏移；^uintptr(0) 表示回退 Stack。
var goidOff atomic.Uintptr

func goid() uint64 {
	off := goidOff.Load()
	switch off {
	case 0:
		return goidInit()
	case ^uintptr(0):
		return goidFromStack()
	default:
		gp := getg.G()
		if gp == nil {
			return goidFromStack()
		}
		return *(*uint64)(unsafe.Add(gp, off))
	}
}

func goidInit() uint64 {
	id, off, ok := discoverGoidOffset()
	if !ok {
		goidOff.Store(^uintptr(0))
		return id
	}
	goidOff.Store(off)
	return id
}

func discoverGoidOffset() (uint64, uintptr, bool) {
	type result struct {
		id  uint64
		off uintptr
		ok  bool
	}
	ch := make(chan result, 1)
	go func() {
		id := goidFromStack()
		gp := getg.G()
		if gp == nil {
			ch <- result{id: id}
			return
		}
		var off uintptr
		matches := 0
		for o := uintptr(0); o < 384; o += 8 {
			if *(*uint64)(unsafe.Add(gp, o)) == id {
				matches++
				off = o
			}
		}
		ch <- result{id: id, off: off, ok: matches == 1}
	}()
	r := <-ch
	return r.id, r.off, r.ok
}

func goidFromStack() uint64 {
	var buf [32]byte
	n := runtime.Stack(buf[:], false)
	b := buf[:n]
	if len(b) < 12 || string(b[:10]) != "goroutine " {
		return 0
	}
	b = b[10:]
	i := 0
	for i < len(b) && b[i] >= '0' && b[i] <= '9' {
		i++
	}
	id, _ := strconv.ParseUint(unsafe.String(&b[0], i), 10, 64)
	return id
}
