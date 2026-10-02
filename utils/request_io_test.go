package utils

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
)

func TestRequestFileReadWrite(t *testing.T) {
	for _, size := range []int{0, 10, 512, 32768} {
		path := filepath.Join(t.TempDir(), "file")
		content := bytes.Repeat([]byte("x"), size)
		ctx, cancel := context.WithCancel(context.Background())
		if n, err := WriteFileContext(ctx, path, content, false); err != nil || n != size {
			t.Fatalf("write = %d, %v", n, err)
		}
		if n, err := WriteFileContext(ctx, path, []byte("end"), true); err != nil || n != 3 {
			t.Fatalf("append = %d, %v", n, err)
		}
		read, err := ReadFileContext(ctx, path)
		if err != nil || !bytes.Equal(read, append(content, []byte("end")...)) {
			t.Fatalf("read = %d bytes, %v", len(read), err)
		}
		cancel()
	}
}

func TestCanceledWriteDoesNotTruncate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	func() {
		defer func() {
			if !data.IsRequestCanceled(recover()) {
				t.Fatal("write ignored cancellation")
			}
		}()
		_, _ = WriteFileContext(ctx, path, []byte("replacement"), false)
	}()
	content, err := os.ReadFile(path)
	if err != nil || string(content) != "original" {
		t.Fatalf("canceled write altered file: %q, %v", content, err)
	}
}

func BenchmarkRequestReadFile(b *testing.B) {
	path := filepath.Join(b.TempDir(), "file")
	if err := os.WriteFile(path, bytes.Repeat([]byte("x"), 32768), 0600); err != nil {
		b.Fatal(err)
	}
	for _, cancellable := range []bool{false, true} {
		name := "background"
		ctx := context.Background()
		if cancellable {
			name = "request"
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			defer cancel()
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := ReadFileContext(ctx, path); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
