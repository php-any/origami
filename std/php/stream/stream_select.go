package stream

import (
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
	return time.Duration(intContextArg(ctx, 3))*time.Second +
		time.Duration(intContextArg(ctx, 4))*time.Microsecond
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
