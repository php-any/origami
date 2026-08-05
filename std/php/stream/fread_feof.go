package stream

import (
	"io"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

func streamResource(ctx data.Context) interface{} {
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil
	}
	resource, ok := value.(*core.ResourceValue)
	if !ok {
		return nil
	}
	return resource.GetResource()
}

type FreadFunction struct{}

func NewFreadFunction() data.FuncStmt { return &FreadFunction{} }

func (f *FreadFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	resource := streamResource(ctx)
	reader, ok := resource.(io.Reader)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	lengthValue, ok := ctx.GetIndexValue(1)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	asInt, ok := lengthValue.(data.AsInt)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	length, err := asInt.AsInt()
	if err != nil || length < 1 {
		return data.NewBoolValue(false), nil
	}

	buffer := make([]byte, length)
	n, readErr := reader.Read(buffer)
	if readErr != nil && readErr != io.EOF && !isBenignPipeReadError(readErr) && n == 0 {
		// 非阻塞流暂无数据（EAGAIN/EWOULDBLOCK）与 PHP 一样返回空字符串。
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(string(buffer[:n])), nil
}

func (f *FreadFunction) GetName() string { return "fread" }
func (f *FreadFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "length", 1, nil, data.Int{}),
	}
}
func (f *FreadFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, nil),
		node.NewVariable(nil, "length", 1, data.Int{}),
	}
}

type FeofFunction struct{}

func NewFeofFunction() data.FuncStmt { return &FeofFunction{} }

func (f *FeofFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	switch stream := streamResource(ctx).(type) {
	case *StreamInfoFromReader:
		return data.NewBoolValue(stream.IsClosed() || stream.IsEOF()), nil
	case *StreamInfo:
		if stream.IsClosed() || stream.File == nil {
			return data.NewBoolValue(true), nil
		}
		current, err := stream.File.Seek(0, io.SeekCurrent)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		stat, err := stream.File.Stat()
		return data.NewBoolValue(err == nil && current >= stat.Size()), nil
	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *FeofFunction) GetName() string { return "feof" }
func (f *FeofFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "stream", 0, nil, nil)}
}
func (f *FeofFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "stream", 0, nil)}
}
