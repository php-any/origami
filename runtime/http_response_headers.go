package runtime

import "sync"

type httpResponseHeaders struct {
	mu      sync.Mutex
	headers []string
}

func (s *httpResponseHeaders) SetLastHTTPResponseHeaders(headers []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.headers = append([]string(nil), headers...)
}
func (s *httpResponseHeaders) LastHTTPResponseHeaders() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.headers...)
}
