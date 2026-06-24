package http

import (
	httpsrc "net/http"
	"net/http/httptest"
	"testing"

	"github.com/php-any/origami/data"
)

func TestBufferedWriter_SetCookie(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	bw.SetCookie(&httpsrc.Cookie{Name: "sid", Value: "abc", Path: "/"})
	bw.sendHeader()

	got := rec.Header().Get("Set-Cookie")
	if got == "" {
		t.Fatal("Set-Cookie header missing")
	}
}

func TestRequestAttrs_SetAndGet(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	attachRequestAttrs(req)
	defer detachRequestAttrs(req)

	bag := requestAttrs(req)
	bag["user"] = data.NewStringValue("alice")

	if got := bag["user"].AsString(); got != "alice" {
		t.Fatalf("attribute = %q, want alice", got)
	}
	_ = rec
}

func TestApplyMiddlewares_PriorityOrder(t *testing.T) {
	var order []int
	mk := func(id int) MiddlewareFunc {
		return func(next httpsrc.Handler) httpsrc.Handler {
			return httpsrc.HandlerFunc(func(w httpsrc.ResponseWriter, r *httpsrc.Request) {
				order = append(order, id)
				next.ServeHTTP(w, r)
			})
		}
	}

	entries := []middlewareEntry{
		{priority: 10, fn: mk(2)},
		{priority: 0, fn: mk(1)},
	}
	h := applyMiddlewares(httpsrc.HandlerFunc(func(w httpsrc.ResponseWriter, r *httpsrc.Request) {
		order = append(order, 3)
	}), entries)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("order = %v, want [1 2 3]", order)
	}
}
