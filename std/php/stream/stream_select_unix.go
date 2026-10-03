//go:build !windows

package stream

import (
	"fmt"
	"reflect"
	"time"

	"github.com/php-any/origami/data"
	"golang.org/x/sys/unix"
)

func fdSetBit(set *unix.FdSet, fd int) bool {
	bits := reflect.ValueOf(set).Elem().FieldByName("Bits")
	if !bits.IsValid() || bits.Len() == 0 {
		return false
	}
	wordBits := int(bits.Index(0).Type().Bits())
	index := fd / wordBits
	if index < 0 || index >= bits.Len() {
		return false
	}
	word := bits.Index(index)
	word.SetInt(word.Int() | int64(uint64(1)<<uint(fd%wordBits)))
	return true
}

func fdIsSet(set *unix.FdSet, fd int) bool {
	bits := reflect.ValueOf(set).Elem().FieldByName("Bits")
	if !bits.IsValid() || bits.Len() == 0 {
		return false
	}
	wordBits := int(bits.Index(0).Type().Bits())
	index := fd / wordBits
	if index < 0 || index >= bits.Len() {
		return false
	}
	return uint64(bits.Index(index).Int())&(uint64(1)<<uint(fd%wordBits)) != 0
}

func addStreamCollection(set *unix.FdSet, value data.Value, maxFD *int) {
	rangeStreamCollection(value, func(stream data.Value) bool {
		fd, ok := streamFileDescriptor(stream)
		if ok && fdSetBit(set, fd) && fd > *maxFD {
			*maxFD = fd
		}
		return true
	})
}

func filterStreamCollection(value data.Value, set *unix.FdSet) {
	switch collection := value.(type) {
	case *data.ArrayValue:
		collection.FilterSlots(func(_ int, slot *data.ZVal) bool {
			if slot == nil {
				return false
			}
			fd, ok := streamFileDescriptor(slot.ReadValue())
			return ok && fdIsSet(set, fd)
		})

	}
}

func (f *StreamSelectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	request := ctx.GoContext()
	if request.Err() != nil {
		panic(data.ErrRequestCanceled)
	}
	if ctl := validateSelectTimeout(ctx); ctl != nil {
		return nil, ctl
	}
	read := data.CowSeparateIndex(ctx, 0)
	write := data.CowSeparateIndex(ctx, 1)
	except := data.CowSeparateIndex(ctx, 2)

	var readSet, writeSet, exceptSet unix.FdSet
	maxFD := -1
	addStreamCollection(&readSet, read, &maxFD)
	addStreamCollection(&writeSet, write, &maxFD)
	addStreamCollection(&exceptSet, except, &maxFD)
	if maxFD < 0 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("No stream arrays were passed"), "ValueError")
	}

	timeout := streamSelectTimeout(ctx)
	infinite := timeout < 0
	end := time.Now().Add(timeout)
	for {
		wait := timeout
		if wait < 0 {
			wait = 50 * time.Millisecond
		}
		if request.Done() != nil && wait > 50*time.Millisecond {
			wait = 50 * time.Millisecond
		}
		timeval := unix.NsecToTimeval(wait.Nanoseconds())
		r, w, e := readSet, writeSet, exceptSet
		ready, err := unix.Select(maxFD+1, &r, &w, &e, &timeval)
		if request.Err() != nil {
			panic(data.ErrRequestCanceled)
		}
		if err != nil {
			if err == unix.EINTR && (infinite || time.Now().Before(end)) {
				if !infinite {
					timeout = time.Until(end)
				}
				continue
			}
			return data.NewBoolValue(false), nil
		}
		if ready > 0 || (!infinite && !time.Now().Before(end)) {
			filterStreamCollection(read, &r)
			filterStreamCollection(write, &w)
			filterStreamCollection(except, &e)
			return data.NewIntValue(ready), nil
		}
		if !infinite {
			timeout = time.Until(end)
		}
	}
}
