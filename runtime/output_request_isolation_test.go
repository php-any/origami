package runtime

import (
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/php-any/origami/parser"
)

// TestRequestOutputIsolation 验证共享 base VM 上并发请求的输出缓冲互不串扰。
func TestRequestOutputIsolation(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	const n = 32
	var wg sync.WaitGroup
	errs := make(chan string, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			restore := BeginRequestOutput()
			defer restore()

			base.StartOutputBuffer()
			marker := "req-" + strconv.Itoa(i)
			base.WriteOutput(marker)
			got, ok := base.CleanOutputBuffer()
			if !ok {
				errs <- "clean not ok"
				return
			}
			if got != marker {
				errs <- "got " + got + " want " + marker
			}
			if base.OutputBufferLevel() != 0 {
				errs <- "level not 0 after clean"
			}
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatal(e)
	}
}

func TestMutedRequestStdoutDoesNotUseFallback(t *testing.T) {
	restore := BeginRequestOutput()
	defer restore()
	MuteRequestStdout()

	fallback := false
	st := currentRequestOutput()
	if st == nil {
		t.Fatal("request output missing")
	}
	st.write("page-html", func(string) { fallback = true })
	if fallback {
		t.Fatal("muted request sink must not fall back to stdout")
	}
}

func TestTakeRequestOutputDrainsOctaneStyleBuffer(t *testing.T) {
	restore := BeginRequestOutput()
	defer restore()
	StartRequestOutputBuffer()
	MuteRequestStdout()
	st := currentRequestOutput()
	if st == nil {
		t.Fatal("request output missing")
	}
	st.write("welcome-html", func(string) { t.Fatal("must not fall back to stdout") })
	got := TakeRequestOutput()
	if got != "welcome-html" {
		t.Fatalf("leftover=%q", got)
	}
	if st.level() != 0 {
		t.Fatalf("level=%d after take", st.level())
	}
}

func TestBootContextEchoUsesLiveRequestOutput(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	boot := base.CreateContext(nil).(*Context)
	restore := BeginRequestOutput()
	defer restore()

	st := currentRequestOutput()
	if st == nil {
		t.Fatal("request output missing")
	}
	st.start()
	boot.WriteOutput("from-boot-ctx")
	got, ok := st.contents()
	if !ok || got != "from-boot-ctx" {
		t.Fatalf("boot context echo must use this goroutine request buffer, got %q ok=%v", got, ok)
	}
}

func TestConcurrentContextsDoNotShareOutput(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	const n = 32
	var wg sync.WaitGroup
	errs := make(chan string, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			restore := BeginRequestOutput()
			defer restore()
			ctx := base.CreateContext(nil).(*Context)
			child := ctx.CreateContext(nil).(*Context)
			mine := currentRequestOutput()
			if mine == nil || child.out != mine || ctx.out != mine {
				errs <- "context did not bind this goroutine output"
				return
			}
			ctx.StartOutputBuffer()
			marker := "req-" + strconv.Itoa(i)
			ctx.WriteOutput(marker)
			got, ok := ctx.CleanOutputBuffer()
			if !ok || got != marker {
				errs <- "got " + got
			}
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatal(e)
	}
}

func TestContextCachesOutputState(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	restore := BeginRequestOutput()
	defer restore()

	ctx := base.CreateContext(nil).(*Context)
	if ctx.out == nil {
		t.Fatal("CreateContext should bind request outputState")
	}
	child := ctx.CreateContext(nil).(*Context)
	if child.out != ctx.out {
		t.Fatal("child context should inherit outputState pointer")
	}
	ctx.StartOutputBuffer()
	ctx.WriteOutput("via-ctx")
	got, ok := ctx.CleanOutputBuffer()
	if !ok || got != "via-ctx" {
		t.Fatalf("context buffer = %q ok=%v", got, ok)
	}
}

func TestOutputStateConcurrentWrites(t *testing.T) {
	st := NewOutputState()
	st.start()
	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			st.write("["+strconv.Itoa(i)+"]", nil)
		}()
	}
	wg.Wait()
	got, ok := st.contents()
	if !ok {
		t.Fatal("expected buffer")
	}
	var total int
	for i := 0; i < n; i++ {
		mark := "[" + strconv.Itoa(i) + "]"
		if !strings.Contains(got, mark) {
			t.Fatalf("missing %q in %q", mark, got)
		}
		total += len(mark)
	}
	if len(got) != total {
		t.Fatalf("len=%d want %d got=%q", len(got), total, got)
	}
}
