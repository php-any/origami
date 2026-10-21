package stream

import (
	"github.com/php-any/origami/data"
	"io"
)

type RewindFunction struct{}

func (*RewindFunction) GetName() string               { return "rewind" }
func (*RewindFunction) GetParams() []data.GetValue    { return fseekFunctionGetParams[:1] }
func (*RewindFunction) GetVariables() []data.Variable { return fseekFunctionGetVariables[:1] }
func (*RewindFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	seeker, ok := streamResource(ctx).(io.Seeker)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	_, err := seeker.Seek(0, io.SeekStart)
	return data.NewBoolValue(err == nil), nil
}
