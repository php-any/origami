package netdata

import (
	"strconv"
	"sync"
	"testing"
)

func TestHTTPRoutesConcurrentAccess(t *testing.T) {
	ClearHTTPRoutes()
	defer ClearHTTPRoutes()

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			AppendHTTPRoute(Route{Method: "GET", Path: "/" + strconv.Itoa(index)})
			_ = HTTPRoutes()
		}(i)
	}
	wg.Wait()

	if got := len(HTTPRoutes()); got != 64 {
		t.Fatalf("route count = %d, want 64", got)
	}
}

func TestHTTPRoutesReturnsSnapshot(t *testing.T) {
	ClearHTTPRoutes()
	defer ClearHTTPRoutes()
	AppendHTTPRoute(Route{Method: "GET", Path: "/original"})

	snapshot := HTTPRoutes()
	snapshot[0].Path = "/changed"
	if got := HTTPRoutes()[0].Path; got != "/original" {
		t.Fatalf("global route changed through snapshot: %q", got)
	}
}
