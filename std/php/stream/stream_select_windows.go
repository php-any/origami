//go:build windows

package stream

import (
	"time"

	"github.com/php-any/origami/data"
)

// Call 在 Windows 上对文件流按 PHP 惯例视为立即就绪；无描述符时等待超时后返回 0。
func (f *StreamSelectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	read, _ := ctx.GetIndexValue(0)
	write, _ := ctx.GetIndexValue(1)
	except, _ := ctx.GetIndexValue(2)

	ready := countStreamCollection(read) + countStreamCollection(write) + countStreamCollection(except)
	if ready > 0 {
		return data.NewIntValue(ready), nil
	}

	if timeout := streamSelectTimeout(ctx); timeout > 0 {
		time.Sleep(timeout)
	}
	return data.NewIntValue(0), nil
}
