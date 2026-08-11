package fpm

import (
	"net/http/httptest"
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
