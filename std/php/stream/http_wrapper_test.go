package stream

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPGetContentsPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Fatalf("content-type = %q", ct)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "a=1&b=2" {
			t.Fatalf("body = %q", string(body))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	sc := NewStreamContext(map[string]map[string]string{
		"http": {
			"method":        "POST",
			"header":        "Content-Type: application/x-www-form-urlencoded",
			"content":       "a=1&b=2",
			"ignore_errors": "1",
			"timeout":       "5",
		},
	}, nil)

	body, ok := HTTPGetContents(server.URL, sc)
	if !ok {
		t.Fatal("HTTPGetContents returned false")
	}
	if body != `{"ok":true}` {
		t.Fatalf("body = %q", body)
	}
}

func TestHTTPGetContentsStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("missing"))
	}))
	defer server.Close()

	body, ok := HTTPGetContents(server.URL, nil)
	if ok {
		t.Fatalf("expected failure, got body %q", body)
	}
}

func TestHTTPGetContentsIgnoreErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("missing"))
	}))
	defer server.Close()

	sc := NewStreamContext(map[string]map[string]string{
		"http": {"ignore_errors": "1"},
	}, nil)

	body, ok := HTTPGetContents(server.URL, sc)
	if !ok {
		t.Fatal("HTTPGetContents returned false")
	}
	if body != "missing" {
		t.Fatalf("body = %q", body)
	}
}
