package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBodyLimitBeforeHandlerIncludingChunked(t *testing.T) {
	for _, length := range []int64{5, -1} {
		called := false
		h := limitedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }), map[string]int{"maxRequests": 1, "maxBodyBytes": 4, "requestTimeoutMs": 1000})
		r := httptest.NewRequest("POST", "/", bytes.NewBufferString("12345"))
		r.ContentLength = length
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 413 || called {
			t.Fatalf("length=%d status=%d called=%v", length, w.Code, called)
		}
	}
	h := limitedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, e := io.ReadAll(r.Body)
		if e != nil || string(b) != "1234" {
			t.Errorf("body=%q err=%v", b, e)
		}
	}), map[string]int{"maxRequests": 1, "maxBodyBytes": 4, "requestTimeoutMs": 1000})
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/", bytes.NewBufferString("1234")))
}
func TestConcurrentRequestLimitAndDeadline(t *testing.T) {
	entered := make(chan struct{})
	done := make(chan struct{})
	h := limitedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(entered); <-r.Context().Done() }), map[string]int{"maxRequests": 1, "maxBodyBytes": 4, "requestTimeoutMs": 100})
	go func() { h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil)); close(done) }()
	<-entered
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 503 {
		t.Fatal(w.Code)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("request deadline not propagated")
	}
}
