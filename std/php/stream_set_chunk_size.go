package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StreamSetChunkSizeFunction 实现 stream_set_chunk_size(resource $stream, int $size): int
// 简化：不改变底层缓冲，原样返回 $size。
type StreamSetChunkSizeFunction struct{}

func NewStreamSetChunkSizeFunction() data.FuncStmt {
	return &StreamSetChunkSizeFunction{}
}

func (f *StreamSetChunkSizeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	size, _ := ctx.GetIndexValue(1)
	if size == nil {
		return data.NewIntValue(0), nil
	}
	if iv, ok := size.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return data.NewIntValue(n), nil
		}
	}
	return data.NewIntValue(0), nil
}

func (f *StreamSetChunkSizeFunction) GetName() string { return "stream_set_chunk_size" }

func (f *StreamSetChunkSizeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "size", 1, nil, data.Int{}),
	}
}

func (f *StreamSetChunkSizeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, nil),
		node.NewVariable(nil, "size", 1, data.Int{}),
	}
}
