package stream

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"io"
)

type StreamCopyToStreamFunction struct{}

func (*StreamCopyToStreamFunction) GetName() string               { return "stream_copy_to_stream" }
func (*StreamCopyToStreamFunction) GetParams() []data.GetValue    { return streamCopyParams }
func (*StreamCopyToStreamFunction) GetVariables() []data.Variable { return streamCopyVariables }

var streamCopyParams = []data.GetValue{node.NewParameter(nil, "from", 0, nil, nil), node.NewParameter(nil, "to", 1, nil, nil), node.NewParameter(nil, "length", 2, data.NewNullValue(), nil), node.NewParameter(nil, "offset", 3, data.NewIntValue(0), nil)}
var streamCopyVariables = []data.Variable{node.NewVariable(nil, "from", 0, nil), node.NewVariable(nil, "to", 1, nil), node.NewVariable(nil, "length", 2, nil), node.NewVariable(nil, "offset", 3, nil)}

func (*StreamCopyToStreamFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	data.CheckRequest(ctx.GoContext())
	defer data.CheckRequest(ctx.GoContext())
	from, _ := ctx.GetIndexValue(0)
	to, _ := ctx.GetIndexValue(1)
	source, ok1 := from.(*core.ResourceValue)
	target, ok2 := to.(*core.ResourceValue)
	if !ok1 || !ok2 || !source.IsPHPResourceOpen() || !target.IsPHPResourceOpen() {
		return nil, data.NewTypeError(nil, fmt.Errorf("stream_copy_to_stream(): supplied resource is not a valid stream resource"))
	}
	reader, ok1 := source.GetResource().(io.Reader)
	writer, ok2 := target.GetResource().(io.Writer)
	if !ok1 || !ok2 {
		return data.NewBoolValue(false), nil
	}
	offset := 0
	if v, _ := ctx.GetIndexValue(3); v != nil {
		if n, ok := v.(data.AsInt); ok {
			offset, _ = n.AsInt()
		}
	}
	if offset < 0 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("stream_copy_to_stream(): offset must be greater than or equal to 0"), "ValueError")
	}
	if offset > 0 {
		seeker, ok := reader.(io.Seeker)
		if !ok {
			return data.NewBoolValue(false), nil
		}
		if _, err := seeker.Seek(int64(offset), io.SeekStart); err != nil {
			return data.NewBoolValue(false), nil
		}
	}
	length := -1
	if v, _ := ctx.GetIndexValue(2); v != nil {
		_, isNull := v.(*data.NullValue)
		if n, ok := v.(data.AsInt); ok && !isNull {
			length, _ = n.AsInt()
		}
	}
	if length < -1 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("stream_copy_to_stream(): length must be greater than or equal to -1"), "ValueError")
	}
	if length >= 0 {
		reader = io.LimitReader(reader, int64(length))
	}
	n, err := io.Copy(writer, reader)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(int(n)), nil
}
