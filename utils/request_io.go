package utils

import (
	"context"
	"io"
	"os"

	"github.com/php-any/origami/data"
)

// OpenRequestFile registers the owned handle, not a pooled VM Context. Callers
// stop the callback and close the file when their synchronous operation ends.
// An OS open/stat call itself can still block on an unresponsive filesystem.
func OpenRequestFile(ctx context.Context, name string, flags int, perm os.FileMode) (*os.File, func(), error) {
	data.CheckRequest(ctx)
	file, err := os.OpenFile(name, flags, perm)
	if err != nil {
		data.CheckRequest(ctx)
		return nil, nil, err
	}
	if ctx.Done() == nil {
		return file, func() { _ = file.Close() }, nil
	}
	stop := context.AfterFunc(ctx, func() { _ = file.Close() })
	return file, func() { stop(); _ = file.Close() }, nil
}

func ReadFileContext(ctx context.Context, name string) ([]byte, error) {
	if ctx.Done() == nil {
		return os.ReadFile(name)
	}
	file, closeFile, err := OpenRequestFile(ctx, name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer closeFile()
	defer data.CheckRequest(ctx)
	capacity := 512
	if stat, err := file.Stat(); err == nil && stat.Size() >= 0 {
		size := stat.Size()
		if size < int64(int(^uint(0)>>1)) && size+1 > int64(capacity) {
			capacity = int(size) + 1
		}
	}
	buffer := make([]byte, 0, capacity)
	for {
		data.CheckRequest(ctx)
		if len(buffer) == cap(buffer) {
			buffer = append(buffer, 0)[:len(buffer)]
		}
		n, err := file.Read(buffer[len(buffer):cap(buffer)])
		buffer = buffer[:len(buffer)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return buffer, err
		}
	}
}

func WriteFileContext(ctx context.Context, name string, content []byte, appendMode bool) (int, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if appendMode {
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}
	file, closeFile, err := OpenRequestFile(ctx, name, flags, 0644)
	if err != nil {
		return 0, err
	}
	defer closeFile()
	defer data.CheckRequest(ctx)
	n, err := file.Write(content)
	if err == nil {
		err = file.Close()
	}
	return n, err
}
