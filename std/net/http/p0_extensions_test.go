package http

import (
	httpsrc "net/http"
	"net/http/httptest"
	"sync"
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
	bag.Lock()
	bag.values["user"] = data.NewStringValue("alice")
	bag.Unlock()

	bag.RLock()
	got := bag.values["user"].AsString()
	bag.RUnlock()
	if got != "alice" {
		t.Fatalf("attribute = %q, want alice", got)
	}
	_ = rec
}

func TestRequestAttrs_ConcurrentAccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	attachRequestAttrs(req)
	defer detachRequestAttrs(req)

	bag := requestAttrs(req)
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := data.NewIntValue(index).AsString()
			bag.Lock()
			bag.values[key] = data.NewIntValue(index)
			bag.Unlock()
			bag.RLock()
			_ = bag.values[key]
			bag.RUnlock()
		}(i)
	}
	wg.Wait()
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
