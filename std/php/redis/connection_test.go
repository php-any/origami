package redis

import (
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/php-any/origami/data"
)

func liveClient(t *testing.T, persistent bool, id string) *clientState {
	t.Helper()
	port, _ := strconv.Atoi(os.Getenv("ORIGAMI_REDIS_TEST_PORT"))
	if port == 0 {
		t.Skip("set ORIGAMI_REDIS_TEST_PORT for live connection tests")
	}
	host := os.Getenv("ORIGAMI_REDIS_TEST_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	s := newState()
	_, ctl := openConnection(nil, s, []data.Value{data.NewStringValue(host), data.NewIntValue(port), data.NewIntValue(2), data.NewStringValue(id), data.NewIntValue(0), data.NewIntValue(2)}, persistent)
	if ctl != nil {
		t.Fatalf("connect: %v", ctl)
	}
	t.Cleanup(s.disconnect)
	return s
}
func liveCommand(t *testing.T, s *clientState, args ...string) any {
	t.Helper()
	v, err := s.exchange(context.Background(), args)
	if err != nil {
		t.Fatal(err)
	}
	if e, ok := v.(serverError); ok {
		t.Fatal(e)
	}
	return v
}
func killClient(t *testing.T, admin, victim *clientState) {
	t.Helper()
	id := liveCommand(t, victim, "CLIENT", "ID").(int64)
	liveCommand(t, admin, "CLIENT", "KILL", "ID", strconv.FormatInt(id, 10))
	deadline := time.Now().Add(time.Second)
	for socketIdle(victim.conn) && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if socketIdle(victim.conn) {
		t.Fatal("closed peer was not detected")
	}
}

func TestRedisIntegrationPersistentPool(t *testing.T) {
	id := "pool-test-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	first := liveClient(t, true, id)
	connectionID := liveCommand(t, first, "CLIENT", "ID")
	first.release() // Same lifecycle path as the client finalizer; no GC timing dependency.
	second := liveClient(t, true, id)
	if got := liveCommand(t, second, "CLIENT", "ID"); got != connectionID {
		t.Fatalf("socket not reused: %v != %v", got, connectionID)
	}
	third := liveClient(t, true, id)
	if liveCommand(t, third, "CLIENT", "ID") == connectionID {
		t.Fatal("live instances share a socket")
	}
	otherID := liveClient(t, true, id+"-other")
	if liveCommand(t, otherID, "CLIENT", "ID") == connectionID {
		t.Fatal("persistent ID isolation")
	}
	second.disconnect() // Explicit close force-closes even a persistent socket.
	fourth := liveClient(t, true, id)
	if liveCommand(t, fourth, "CLIENT", "ID") == connectionID {
		t.Fatal("closed socket returned to pool")
	}
	liveCommand(t, fourth, "WATCH", "origami-pool-unused")
	fourth.watching = true
	fourth.release()
	if fourth.conn != nil {
		t.Fatal("WATCH connection retained")
	}
}

func TestRedisIntegrationReconnectPreservesDatabase(t *testing.T) {
	admin := liveClient(t, false, "")
	s := liveClient(t, false, "")
	liveCommand(t, s, "SELECT", "1")
	s.database = 1
	old := liveCommand(t, s, "CLIENT", "ID")
	killClient(t, admin, s)
	if err := s.ensureConnected(context.Background()); err != nil {
		t.Fatal(err)
	}
	if liveCommand(t, s, "CLIENT", "ID") == old {
		t.Fatal("no reconnect")
	}
	info := liveCommand(t, s, "CLIENT", "INFO").(string)
	if !strings.Contains(info, " db=1 ") {
		t.Fatalf("database was not restored: %s", info)
	}
	s.maxRetries = 0
	killClient(t, admin, s)
	if _, err := s.exchange(context.Background(), []string{"PING"}); err == nil {
		t.Fatal("dead socket succeeded")
	}
	if err := s.ensureConnected(context.Background()); err == nil {
		t.Fatal("disabled retries reconnected")
	}
}

func TestRedisIntegrationReconnectPreservesAuth(t *testing.T) {
	if os.Getenv("ORIGAMI_REDIS_TEST_ACL") != "1" {
		t.Skip("requires disposable ACL-enabled test server")
	}
	admin := liveClient(t, false, "")
	s := liveClient(t, false, "")
	user := "origami-reconnect-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	liveCommand(t, admin, "ACL", "SETUSER", user, "on", ">temporary-test-password", "+@all", "~*")
	defer liveCommand(t, admin, "ACL", "DELUSER", user)
	s.auth = []string{"AUTH", user, "temporary-test-password"}
	liveCommand(t, s, s.auth...)
	killClient(t, admin, s)
	if err := s.ensureConnected(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := liveCommand(t, s, "ACL", "WHOAMI"); got != user {
		t.Fatalf("auth not restored: %v", got)
	}
}

func TestRedisIntegrationWatchNeverReconnects(t *testing.T) {
	admin := liveClient(t, false, "")
	s := liveClient(t, false, "")
	liveCommand(t, s, "WATCH", "origami-watch-unused")
	s.watching = true
	killClient(t, admin, s)
	if err := s.ensureConnected(context.Background()); err == nil {
		t.Fatal("lost WATCH silently reconnected")
	}
	if s.conn != nil {
		t.Fatal("lost WATCH retained socket")
	}
}

func TestBackoffAndCancellation(t *testing.T) {
	s := newState()
	s.backoffBase = 10 * time.Millisecond
	s.backoffCap = 25 * time.Millisecond
	for algorithm := 0; algorithm <= 6; algorithm++ {
		s.backoffAlgorithm = algorithm
		previous := time.Duration(0)
		for i := 0; i < 30; i++ {
			delay := s.backoff(i, &previous)
			if delay < 0 || delay > s.backoffCap {
				t.Fatalf("algorithm %d delay %v", algorithm, delay)
			}
		}
	}
	s.backoffAlgorithm = 4
	previous := time.Duration(0)
	for i, want := range []time.Duration{10, 20, 25} {
		if got := s.backoff(i, &previous); got != want*time.Millisecond {
			t.Fatal(got)
		}
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s.host = "127.0.0.1"
	s.port = listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	s.backoffBase, s.backoffCap, s.backoffAlgorithm = time.Second, time.Second, 6
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	if err := s.ensureConnected(ctx); err != context.DeadlineExceeded {
		t.Fatalf("backoff cancellation: %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("backoff ignored deadline")
	}
}
