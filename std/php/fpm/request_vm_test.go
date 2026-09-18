package fpm

import (
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestRequestIsolation(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	a := New(base, data.DefaultOutputWriter)
	b := New(base, data.DefaultOutputWriter)

	a.EnsureGlobalZVal("token").Value = data.NewStringValue("A")
	if got := b.EnsureGlobalZVal("token").Value.AsString(); got != "" && got != "A" {
		t.Fatalf("unexpected shared global between requests: %q", got)
	}
	b.EnsureGlobalZVal("token").Value = data.NewStringValue("B")
	if got := a.EnsureGlobalZVal("token").Value.AsString(); got != "A" {
		t.Fatalf("request A global leaked: %q", got)
	}

	a.EnsureGlobalsArray().SetProperty("__me_cache", data.NewStringValue("userA"))
	if b.EnsureGlobalsArray().HasProperty("__me_cache") {
		t.Fatal("$GLOBALS leaked between RequestVMs")
	}
	b.EnsureGlobalsArray().SetProperty("__me_cache", data.NewStringValue("userB"))
	gotMe, _ := a.EnsureGlobalsArray().GetProperty("__me_cache")
	if gotMe.AsString() != "userA" {
		t.Fatalf("request A $GLOBALS leaked: %v", gotMe)
	}

	a.EnsureSessionArray().SetProperty("uid", data.NewIntValue(1))
	if b.EnsureSessionArray().HasProperty("uid") {
		t.Fatal("$_SESSION leaked between RequestVMs")
	}

	a.BindHTTP(httptest.NewRequest("GET", "http://example.test/a", nil), httptest.NewRecorder())
	b.BindHTTP(httptest.NewRequest("GET", "http://example.test/b", nil), httptest.NewRecorder())
	if a.HTTPRequest().URL.Path != "/a" || b.HTTPRequest().URL.Path != "/b" {
		t.Fatal("HTTP binding leaked between requests")
	}

	a.PushCallFrame(data.CallFrame{Function: "a"})
	b.PushCallFrame(data.CallFrame{Function: "b"})
	if got := a.SnapshotCallStack(); len(got) != 1 || got[0].Function != "a" {
		t.Fatalf("call stack A: %#v", got)
	}
	if got := b.SnapshotCallStack(); len(got) != 1 || got[0].Function != "b" {
		t.Fatalf("call stack B: %#v", got)
	}
}

func TestRequestVMFullOutputBufferMethods(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	var out strings.Builder
	vm := New(base, func(s string) { _, _ = out.WriteString(s) })

	// 无缓冲时 OutputBufferContents / Level / Length 均为空
	if s, ok := vm.OutputBufferContents(); ok || s != "" {
		t.Fatalf("no buffer expected: s=%q ok=%v", s, ok)
	}
	if l := vm.OutputBufferLevel(); l != 0 {
		t.Fatalf("level=%d, want 0", l)
	}
	if n, ok := vm.OutputBufferLength(); ok || n != 0 {
		t.Fatalf("length: n=%d ok=%v", n, ok)
	}
	if n, ok := vm.FlushOutputBuffer(); ok || n != "" {
		t.Fatalf("flush empty: n=%q ok=%v", n, ok)
	}
	if vm.CleanCurrentBuffer() {
		t.Fatal("clean empty should be false")
	}

	// 开启一层缓冲
	vm.StartOutputBuffer()
	vm.WriteOutput("hello")
	if s, ok := vm.OutputBufferContents(); !ok || s != "hello" {
		t.Fatalf("contents: s=%q ok=%v", s, ok)
	}
	if l := vm.OutputBufferLevel(); l != 1 {
		t.Fatalf("level=%d, want 1", l)
	}
	if n, ok := vm.OutputBufferLength(); !ok || n != 5 {
		t.Fatalf("length: n=%d ok=%v", n, ok)
	}

	// CleanCurrentBuffer 清空但不结束缓冲
	if !vm.CleanCurrentBuffer() {
		t.Fatal("clean current should succeed")
	}
	if s, ok := vm.OutputBufferContents(); !ok || s != "" {
		t.Fatalf("after clean: s=%q ok=%v", s, ok)
	}
	if l := vm.OutputBufferLevel(); l != 1 {
		t.Fatalf("after clean level=%d, want 1", l)
	}

	// status / handlers
	st := vm.OutputBufferStatus(false)
	if len(st) != 1 || st[0].Level != 1 || st[0].Name != "default output handler" {
		t.Fatalf("status: %#v", st)
	}
	hd := vm.ListOutputHandlers()
	if len(hd) != 1 || hd[0] != "default output handler" {
		t.Fatalf("handlers: %#v", hd)
	}

	// 嵌套缓冲 + status(full=true)
	vm.StartOutputBuffer()
	vm.WriteOutput("world")
	if l := vm.OutputBufferLevel(); l != 2 {
		t.Fatalf("nested level=%d, want 2", l)
	}
	stFull := vm.OutputBufferStatus(true)
	if len(stFull) != 2 {
		t.Fatalf("full status len=%d, want 2", len(stFull))
	}

	// FlushOutputBuffer 弹出并写出到上一层
	if s, ok := vm.FlushOutputBuffer(); !ok || s != "world" {
		t.Fatalf("flush: s=%q ok=%v", s, ok)
	}
	if l := vm.OutputBufferLevel(); l != 1 {
		t.Fatalf("after flush level=%d, want 1", l)
	}

	// CleanOutputBuffer 弹出并返回内容（level 1 此时含 flush 写出的 world）
	if s, ok := vm.CleanOutputBuffer(); !ok || s != "world" {
		t.Fatalf("clean: s=%q ok=%v", s, ok)
	}
	if l := vm.OutputBufferLevel(); l != 0 {
		t.Fatalf("after clean-all level=%d, want 0", l)
	}

	// implicitFlush
	vm.SetImplicitFlush(true)
	if !vm.IsImplicitFlush() {
		t.Fatal("implicit flush should be true")
	}
	vm.SetImplicitFlush(false)
	if vm.IsImplicitFlush() {
		t.Fatal("implicit flush should be false")
	}

	// 接口断言
	var _ data.OutputBufferHost = vm
}

func TestGlobalsArrayViaContext(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	a := New(base, data.DefaultOutputWriter)
	b := New(base, data.DefaultOutputWriter)

	ctxA := a.CreateContext(nil)
	ctxB := b.CreateContext(nil)
	gv := node.NewGlobalsArrayVariable(nil)

	valA, _ := gv.GetValue(ctxA)
	objA := valA.(*data.ObjectValue)
	objA.SetProperty("__me_cache", data.NewStringValue("Alice"))

	valB, _ := gv.GetValue(ctxB)
	objB := valB.(*data.ObjectValue)
	if objB.HasProperty("__me_cache") {
		t.Fatal("$GLOBALS via context leaked across RequestVMs (cross-user session)")
	}
	objB.SetProperty("__me_cache", data.NewStringValue("Bob"))

	gotA, _ := objA.GetProperty("__me_cache")
	if gotA.AsString() != "Alice" {
		t.Fatalf("request A identity corrupted: %q", gotA.AsString())
	}
}

func TestRequestVMCallStacksConcurrent(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	const n = 16
	var wg sync.WaitGroup
	errs := make(chan string, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			vm := New(base, data.DefaultOutputWriter)
			name := "req-" + strconv.Itoa(i)
			vm.PushCallFrame(data.CallFrame{Function: name})
			if d := vm.EnterCall(); d != 1 {
				errs <- "depth"
				vm.LeaveCall()
				return
			}
			got := vm.SnapshotCallStack()
			if len(got) != 1 || got[0].Function != name {
				errs <- "stack " + name
			}
			vm.LeaveCall()
			vm.PopCallFrame()
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatal(e)
	}
}

func TestRequestVMContextOutputUsesSink(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	var outA, outB strings.Builder
	a := New(base, func(s string) { _, _ = outA.WriteString(s) })
	b := New(base, func(s string) { _, _ = outB.WriteString(s) })

	ctxA := a.CreateContext(nil)
	ctxB := b.CreateContext(nil)
	sinkA, ok := ctxA.(data.OutputSink)
	if !ok {
		t.Fatal("context should implement OutputSink")
	}
	sinkB := ctxB.(data.OutputSink)
	hostA := ctxA.(data.OutputBufferHost)
	hostB := ctxB.(data.OutputBufferHost)

	sinkA.WriteOutput("echo-a")
	sinkB.WriteOutput("echo-b")
	if outA.String() != "echo-a" || outB.String() != "echo-b" {
		t.Fatalf("unbuffered echo leaked: A=%q B=%q", outA.String(), outB.String())
	}

	hostA.StartOutputBuffer()
	sinkA.WriteOutput("buf-a")
	hostB.StartOutputBuffer()
	sinkB.WriteOutput("buf-b")
	gotA, okA := hostA.CleanOutputBuffer()
	gotB, okB := hostB.CleanOutputBuffer()
	if !okA || gotA != "buf-a" {
		t.Fatalf("request A buffer = %q ok=%v", gotA, okA)
	}
	if !okB || gotB != "buf-b" {
		t.Fatalf("request B buffer = %q ok=%v", gotB, okB)
	}
	if outA.String() != "echo-a" || outB.String() != "echo-b" {
		t.Fatalf("buffered content leaked to sink: A=%q B=%q", outA.String(), outB.String())
	}
}
