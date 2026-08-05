package stream

import (
	"io"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type FseekFunction struct{}

func NewFseekFunction() data.FuncStmt { return &FseekFunction{} }

func (f *FseekFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stream, ok := streamResource(ctx).(interface {
		Seek(int64, int) (int64, error)
	})
	if !ok {
		return data.NewIntValue(-1), nil
	}
	offsetValue, ok := ctx.GetIndexValue(1)
	if !ok {
		return data.NewIntValue(-1), nil
	}
	offsetInt, ok := offsetValue.(data.AsInt)
	if !ok {
		return data.NewIntValue(-1), nil
	}
	offset, err := offsetInt.AsInt()
	if err != nil {
		return data.NewIntValue(-1), nil
	}

	whence := io.SeekStart
	if value, ok := ctx.GetIndexValue(2); ok {
		if asInt, ok := value.(data.AsInt); ok {
			if parsed, err := asInt.AsInt(); err == nil {
				whence = parsed
			}
		}
	}
	if _, err := stream.Seek(int64(offset), whence); err != nil {
		return data.NewIntValue(-1), nil
	}
	return data.NewIntValue(0), nil
}

func (f *FseekFunction) GetName() string { return "fseek" }
func (f *FseekFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "offset", 1, nil, data.Int{}),
		node.NewParameter(nil, "whence", 2, node.NewIntLiteral(nil, "0"), data.Int{}),
	}
}
func (f *FseekFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, nil),
		node.NewVariable(nil, "offset", 1, data.Int{}),
		node.NewVariable(nil, "whence", 2, data.Int{}),
	}
}
