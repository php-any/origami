package stream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPResponseHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Composer-Fixture") != "present" {
			t.Error("request header missing")
		}
		w.Header().Set("X-Response-Fixture", "received")
		w.Write([]byte("package metadata"))
	}))
	defer server.Close()
	sc := NewStreamContext(map[string]map[string]string{"http": {"header": "X-Composer-Fixture: present"}}, nil)
	var headers []string
	body, ok := HTTPGetContentsWithHeaders(context.Background(), server.URL, sc, func(h []string) { headers = h })
	if !ok || body != "package metadata" || !strings.Contains(strings.Join(headers, "\n"), "X-Response-Fixture: received") {
		t.Fatalf("ok=%v body=%q headers=%v", ok, body, headers)
	}
}
