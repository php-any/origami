package stream

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

type FtruncateFunction struct{}

func (*FtruncateFunction) GetName() string { return "ftruncate" }
func (*FtruncateFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "stream", 0, nil, nil), node.NewParameter(nil, "size", 1, nil, nil)}
}
func (*FtruncateFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "stream", 0, nil), node.NewVariable(nil, "size", 1, nil)}
}
func (*FtruncateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	sizeValue, _ := ctx.GetIndexValue(1)
	size := 0
	if n, ok := sizeValue.(data.AsInt); ok {
		size, _ = n.AsInt()
	}
	if size < 0 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("ftruncate(): size must be greater than or equal to 0"), "ValueError")
	}
	resource, ok := v.(*core.ResourceValue)
	if !ok {
		return nil, data.NewTypeError(nil, fmt.Errorf("ftruncate(): argument must be a stream resource"))
	}
	stream, ok := resource.GetResource().(*StreamInfo)
	if !ok || stream.IsClosed() {
		return nil, data.NewTypeError(nil, fmt.Errorf("ftruncate(): invalid stream resource"))
	}
	file := stream.openFile()
	if file == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(file.Truncate(int64(size)) == nil), nil
}
