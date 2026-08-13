package stream

import (
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// StreamSetBlockingFunction 实现 stream_set_blocking(resource, bool): bool。
type StreamSetBlockingFunction struct{}

func NewStreamSetBlockingFunction() data.FuncStmt {
	return &StreamSetBlockingFunction{}
}

func (f *StreamSetBlockingFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	streamValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	resource, ok := streamValue.(*core.ResourceValue)
	if !ok || resource.GetResource() == nil {
		return data.NewBoolValue(false), nil
	}

	blocking := true
	if value, ok := ctx.GetIndexValue(1); ok {
		if asBool, ok := value.(data.AsBool); ok {
			if parsed, err := asBool.AsBool(); err == nil {
				blocking = parsed
			}
		}
	}

	var fd uintptr
	switch stream := resource.GetResource().(type) {
	case *StreamInfo:
		if stream.IsClosed() || stream.File == nil {
			return data.NewBoolValue(false), nil
		}
		fd = stream.File.Fd()
	case *StreamInfoFromReader:
		if stream.IsClosed() || stream.Reader == nil {
			return data.NewBoolValue(false), nil
		}
		file, ok := stream.Reader.(interface{ Fd() uintptr })
		if !ok {
			return data.NewBoolValue(false), nil
		}
		fd = file.Fd()
	default:
		return data.NewBoolValue(false), nil
	}

	if err := syscall.SetNonblock(int(fd), !blocking); err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *StreamSetBlockingFunction) GetName() string { return "stream_set_blocking" }

func (f *StreamSetBlockingFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "enable", 1, nil, data.Bool{}),
	}
}

func (f *StreamSetBlockingFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, nil),
		node.NewVariable(nil, "enable", 1, data.Bool{}),
	}
}
