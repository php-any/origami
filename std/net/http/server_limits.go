package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/net/internal/script"
)

func (s *ServerClass) configureMethod() data.Method {
	return &script.Method{Name: "configure", Args: []string{"options"}, Run: func(c data.Context) (data.GetValue, error) {
		defaults := map[string]int{"maxConnections": 256, "maxRequests": 128, "maxBodyBytes": 1048576, "maxHeaderBytes": 16384, "readHeaderTimeoutMs": 5000, "readTimeoutMs": 10000, "writeTimeoutMs": 10000, "idleTimeoutMs": 30000, "requestTimeoutMs": 10000}
		if s.limits != nil {
			defaults = s.limits
		}
		v, e := script.Options(c, 0, defaults)
		if e == nil {
			s.limits = v
		}
		return nil, e
	}}
}
func limitedHandler(next http.Handler, limits map[string]int) http.Handler {
	slots := make(chan struct{}, limits["maxRequests"])
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case slots <- struct{}{}:
			defer func() { <-slots }()
		default:
			http.Error(w, "Service Unavailable", 503)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(limits["requestTimeoutMs"])*time.Millisecond)
		defer cancel()
		r = r.WithContext(ctx)
		if r.ContentLength > int64(limits["maxBodyBytes"]) {
			http.Error(w, "Request Entity Too Large", 413)
			return
		}
		// Bound chunked bodies before language superglobals or body helpers read them.
		body := http.MaxBytesReader(w, r.Body, int64(limits["maxBodyBytes"]))
		defer body.Close()
		b, e := io.ReadAll(body)
		if e != nil {
			var large *http.MaxBytesError
			if errors.As(e, &large) {
				http.Error(w, "Request Entity Too Large", 413)
			} else {
				http.Error(w, "Cannot read request body", 400)
			}
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(b))
		next.ServeHTTP(w, r)
	})
}
