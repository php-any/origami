package stream

import (
	"reflect"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"golang.org/x/sys/unix"
)

// StreamSelectFunction 使用系统 select 实现 stream_select。
type StreamSelectFunction struct{}

func NewStreamSelectFunction() data.FuncStmt {
	return &StreamSelectFunction{}
}

func streamFileDescriptor(value data.Value) (int, bool) {
	resource, ok := value.(*core.ResourceValue)
	if !ok {
		return 0, false
	}
	switch stream := resource.GetResource().(type) {
	case *StreamInfo:
		if stream.IsClosed() || stream.File == nil {
			return 0, false
		}
		return int(stream.File.Fd()), true
	case *StreamInfoFromReader:
		if stream.IsClosed() || stream.Reader == nil {
			return 0, false
		}
		file, ok := stream.Reader.(interface{ Fd() uintptr })
		if !ok {
			return 0, false
		}
		return int(file.Fd()), true
	}
	return 0, false
}

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

func rangeStreamCollection(value data.Value, fn func(data.Value) bool) {
	switch collection := value.(type) {
	case *data.ArrayValue:
		for _, zv := range collection.List {
			if zv != nil && !fn(zv.Value) {
				return
			}
		}
	case *data.ObjectValue:
		collection.RangeProperties(func(_ string, value data.Value) bool {
			return fn(value)
		})
	}
}

func filterStreamCollection(value data.Value, set *unix.FdSet) {
	switch collection := value.(type) {
	case *data.ArrayValue:
		filtered := collection.List[:0]
		for _, zv := range collection.List {
			if zv == nil {
				continue
			}
			fd, ok := streamFileDescriptor(zv.Value)
			if ok && fdIsSet(set, fd) {
				filtered = append(filtered, zv)
			}
		}
		collection.List = filtered
	case *data.ObjectValue:
		filtered := data.NewObjectValue()
		collection.RangeProperties(func(key string, value data.Value) bool {
			fd, ok := streamFileDescriptor(value)
			if ok && fdIsSet(set, fd) {
				filtered.SetProperty(key, value)
			}
			return true
		})
		*collection = *filtered
	}
}

func intContextArg(ctx data.Context, index int) int {
	value, ok := ctx.GetIndexValue(index)
	if !ok {
		return 0
	}
	asInt, ok := value.(data.AsInt)
	if !ok {
		return 0
	}
	result, _ := asInt.AsInt()
	return result
}

func (f *StreamSelectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	read, _ := ctx.GetIndexValue(0)
	write, _ := ctx.GetIndexValue(1)
	except, _ := ctx.GetIndexValue(2)

	var readSet, writeSet, exceptSet unix.FdSet
	maxFD := -1
	addStreamCollection(&readSet, read, &maxFD)
	addStreamCollection(&writeSet, write, &maxFD)
	addStreamCollection(&exceptSet, except, &maxFD)
	if maxFD < 0 {
		return data.NewIntValue(0), nil
	}

	timeout := time.Duration(intContextArg(ctx, 3))*time.Second +
		time.Duration(intContextArg(ctx, 4))*time.Microsecond
	timeval := unix.NsecToTimeval(timeout.Nanoseconds())
	ready, err := unix.Select(maxFD+1, &readSet, &writeSet, &exceptSet, &timeval)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	filterStreamCollection(read, &readSet)
	filterStreamCollection(write, &writeSet)
	filterStreamCollection(except, &exceptSet)
	return data.NewIntValue(ready), nil
}

func (f *StreamSelectFunction) GetName() string { return "stream_select" }

func (f *StreamSelectFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "read", 0, nil, data.Mixed{}),
		node.NewParameterReference(nil, "write", 1, nil, data.Mixed{}),
		node.NewParameterReference(nil, "except", 2, nil, data.Mixed{}),
		node.NewParameter(nil, "seconds", 3, nil, nil),
		node.NewParameter(nil, "microseconds", 4, node.NewNullLiteral(nil), nil),
	}
}

func (f *StreamSelectFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "read", 0, nil),
		node.NewVariable(nil, "write", 1, nil),
		node.NewVariable(nil, "except", 2, nil),
		node.NewVariable(nil, "seconds", 3, nil),
		node.NewVariable(nil, "microseconds", 4, nil),
	}
}
