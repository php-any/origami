//go:build windows

package stream

import (
	"fmt"
	"github.com/php-any/origami/data"
)

// PHP's Windows file/pipe select reports distinct handles as immediately
// ready and retains their set membership. This does not guarantee a pipe read
// cannot block; resource reads independently observe request cancellation.
func (f *StreamSelectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if ctx.GoContext().Err() != nil {
		panic(data.ErrRequestCanceled)
	}
	if ctl := validateSelectTimeout(ctx); ctl != nil {
		return nil, ctl
	}
	read := data.CowSeparateIndex(ctx, 0)
	write := data.CowSeparateIndex(ctx, 1)
	except := data.CowSeparateIndex(ctx, 2)
	handles := make(map[int]struct{})
	for _, collection := range []data.Value{read, write, except} {
		rangeStreamCollection(collection, func(stream data.Value) bool {
			if fd, ok := streamFileDescriptor(stream); ok {
				handles[fd] = struct{}{}
			}
			return true
		})
	}
	if len(handles) == 0 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("No stream arrays were passed"), "ValueError")
	}
	filterSelectableStreams(read, func(int) bool { return true })
	filterSelectableStreams(write, func(int) bool { return true })
	filterSelectableStreams(except, func(int) bool { return true })
	return data.NewIntValue(len(handles)), nil
}
