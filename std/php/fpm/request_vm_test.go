package fpm

import (
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
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
