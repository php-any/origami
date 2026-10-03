package stream

import (
	"fmt"
	"syscall"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
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
		return int(fileDescriptor(stream.File)), true
	case *StreamInfoFromReader:
		if stream.IsClosed() || stream.Reader == nil {
			return 0, false
		}
		file, ok := stream.Reader.(interface {
			SyscallConn() (syscall.RawConn, error)
		})
		if !ok {
			return 0, false
		}
		return int(fileDescriptor(file)), true
	}
	return 0, false
}

func rangeStreamCollection(value data.Value, fn func(data.Value) bool) {
	switch collection := value.(type) {
	case *data.ArrayValue:
		for arraySlots141, arrayPosition141 := collection.View(), 0; arrayPosition141 < arraySlots141.Len(); arrayPosition141++ {
			zv := arraySlots141.At(arrayPosition141)
			if zv != nil && !fn(zv.ReadValue()) {
				return
			}
		}

	}
}

func countStreamCollection(value data.Value) int {
	n := 0
	rangeStreamCollection(value, func(stream data.Value) bool {
		if _, ok := streamFileDescriptor(stream); ok {
			n++
		}
		return true
	})
	return n
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

func streamSelectTimeout(ctx data.Context) time.Duration {
	seconds, _ := ctx.GetIndexValue(3)
	if _, infinite := seconds.(*data.NullValue); infinite {
		return -1
	}
	return time.Duration(intContextArg(ctx, 3))*time.Second +
		time.Duration(intContextArg(ctx, 4))*time.Microsecond
}

func validateSelectTimeout(ctx data.Context) data.Control {
	if intContextArg(ctx, 3) < 0 || intContextArg(ctx, 4) < 0 {
		return data.NewErrorThrowByName(nil, fmt.Errorf("stream_select(): timeout must be greater than or equal to 0"), "ValueError")
	}
	return nil
}

func filterSelectableStreams(value data.Value, keep func(int) bool) {
	switch collection := value.(type) {
	case *data.ArrayValue:
		collection.FilterSlots(func(_ int, slot *data.ZVal) bool {
			if slot == nil {
				return false
			}
			fd, ok := streamFileDescriptor(slot.ReadValue())
			return ok && keep(fd)
		})

	}
}

func (f *StreamSelectFunction) GetName() string { return "stream_select" }

var streamSelectFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "read", 0, nil, data.Mixed{}),
	node.NewParameterReference(nil, "write", 1, nil, data.Mixed{}),
	node.NewParameterReference(nil, "except", 2, nil, data.Mixed{}),
	node.NewParameter(nil, "seconds", 3, nil, nil),
	node.NewParameter(nil, "microseconds", 4, node.NewNullLiteral(nil), nil),
}

func (f *StreamSelectFunction) GetParams() []data.GetValue {
	return streamSelectFunctionGetParams
}

var streamSelectFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "read", 0, nil),
	node.NewVariable(nil, "write", 1, nil),
	node.NewVariable(nil, "except", 2, nil),
	node.NewVariable(nil, "seconds", 3, nil),
	node.NewVariable(nil, "microseconds", 4, nil),
}

func (f *StreamSelectFunction) GetVariables() []data.Variable {
	return streamSelectFunctionGetVariables
}
