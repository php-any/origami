package stream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/php-any/origami/data"
)

func TestHTTPContentsCancellation(t *testing.T) {
	for _, sendHeaders := range []bool{false, true} {
		name := "headers"
		if sendHeaders {
			name = "body"
		}
		t.Run(name, func(t *testing.T) {
			started := make(chan struct{})
			stopped := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if sendHeaders {
					w.WriteHeader(http.StatusOK)
					w.(http.Flusher).Flush()
				}
				close(started)
				<-r.Context().Done()
				close(stopped)
			}))
			defer server.Close()
			request, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan any, 1)
			go func() {
				defer func() { done <- recover() }()
				_, _ = HTTPGetContentsContext(request, server.URL, nil)
			}()
			select {
			case <-started:
			case <-time.After(time.Second):
				t.Fatal("HTTP request did not start")
			}
			cancel()
			select {
			case result := <-done:
				if !data.IsRequestCanceled(result) {
					t.Fatalf("HTTP cancellation returned normally: %v", result)
				}
			case <-time.After(time.Second):
				t.Fatal("HTTP read did not cancel")
			}
			select {
			case <-stopped:
			case <-time.After(time.Second):
				t.Fatal("HTTP connection was not released")
			}
		})
	}
}
