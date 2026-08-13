package requestvm

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php/core"
)

func TestObIsolationBetweenRequests(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)

	var outA, outB strings.Builder
	a := New(base, func(s string) { _, _ = outA.WriteString(s) })
	b := New(base, func(s string) { _, _ = outB.WriteString(s) })

	ctxA := a.CreateContext(nil)
	ctxB := b.CreateContext(nil)

	start := core.NewObStartFunction()
	clean := core.NewObGetCleanFunction()

	call := func(ctx data.Context, fn data.FuncStmt) data.GetValue {
		c := ctx.CreateContext(fn.GetVariables())
		v, ctl := fn.Call(c)
		if ctl != nil {
			t.Fatalf("unexpected control: %v", ctl)
		}
		return v
	}

	call(ctxA, start)
	a.WriteOutput("AAA")
	call(ctxB, start)
	b.WriteOutput("BBB")

	gotA := call(ctxA, clean).(data.AsString).AsString()
	gotB := call(ctxB, clean).(data.AsString).AsString()
	if gotA != "AAA" || gotB != "BBB" {
		t.Fatalf("ob isolation failed: a=%q b=%q", gotA, gotB)
	}

	a.WriteOutput("bodyA")
	b.WriteOutput("bodyB")
	if outA.String() != "bodyA" || outB.String() != "bodyB" {
		t.Fatalf("output isolation failed: a=%q b=%q", outA.String(), outB.String())
	}
}

func TestSharedClassTable(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	a := New(base, data.DefaultOutputWriter)
	b := New(base, data.DefaultOutputWriter)

	type stub struct{ name string }
	// use AddClass via a proxy: register a simple go class through base
	if control := base.AddClass(&stubClass{name: "App\\Shared"}); control != nil {
		t.Fatalf("AddClass: %v", control)
	}
	if _, ok := a.GetClass("App\\Shared"); !ok {
		t.Fatal("request A missing shared class")
	}
	if _, ok := b.GetClass("App\\Shared"); !ok {
		t.Fatal("request B missing shared class")
	}
}

func TestRequestVMOnlyIsolatesHTTPAndOutputState(t *testing.T) {
	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	a := New(base, data.DefaultOutputWriter)
	b := New(base, data.DefaultOutputWriter)

	a.EnsureGlobalZVal("request_id").Value = data.NewStringValue("A")
	if got := b.EnsureGlobalZVal("request_id").Value.AsString(); got != "A" {
		t.Fatalf("PHP global should be shared through base VM: %q", got)
	}

	if depth := a.EnterCall(); depth != 1 {
		t.Fatalf("request A call depth = %d", depth)
	}
	if depth := b.EnterCall(); depth != 1 {
		t.Fatalf("request B call depth = %d, want isolated per request", depth)
	}
	defer a.LeaveCall()
	defer b.LeaveCall()
	a.PushCallFrame(data.CallFrame{Function: "a"})
	b.PushCallFrame(data.CallFrame{Function: "b"})
	defer a.PopCallFrame()
	defer b.PopCallFrame()
	if got := a.SnapshotCallStack(); len(got) != 1 || got[0].Function != "a" {
		t.Fatalf("request A call stack should be isolated: %#v", got)
	}
	if got := b.SnapshotCallStack(); len(got) != 1 || got[0].Function != "b" {
		t.Fatalf("request B call stack should be isolated: %#v", got)
	}

	a.SetExceptionHandler(data.NewStringValue("handler-a"))
	if got := b.GetExceptionHandler().AsString(); got != "handler-a" {
		t.Fatalf("exception handler should be shared through base VM: %q", got)
	}

	reqA := httptest.NewRequest("GET", "http://example.test/a", nil)
	reqB := httptest.NewRequest("GET", "http://example.test/b", nil)
	a.BindHTTP(reqA, httptest.NewRecorder())
	b.BindHTTP(reqB, httptest.NewRecorder())
	if a.HTTPRequest().URL.Path != "/a" || b.HTTPRequest().URL.Path != "/b" {
		t.Fatal("HTTP request binding leaked")
	}
}

type stubClass struct{ name string }

func (c *stubClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
func (c *stubClass) GetFrom() data.From                       { return nil }
func (c *stubClass) GetName() string                          { return c.name }
func (c *stubClass) GetExtend() *string                       { return nil }
func (c *stubClass) GetImplements() []string                  { return nil }
func (c *stubClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *stubClass) GetPropertyList() []data.Property         { return nil }
func (c *stubClass) GetMethod(string) (data.Method, bool)     { return nil, false }
func (c *stubClass) GetMethods() []data.Method                { return nil }
func (c *stubClass) GetConstruct() data.Method                { return nil }
