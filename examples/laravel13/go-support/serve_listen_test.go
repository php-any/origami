package gosupport

import (
	"testing"
)

func TestListenTCPRejectsOccupiedPort(t *testing.T) {
	first, err := listenTCP("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()

	addr := first.Addr().String()
	second, err := listenTCP(addr)
	if err == nil {
		second.Close()
		t.Fatalf("occupied port %s should not bind again", addr)
	}
	if !isAddrInUse(err) {
		t.Fatalf("expected address-in-use, got %v", err)
	}
}
