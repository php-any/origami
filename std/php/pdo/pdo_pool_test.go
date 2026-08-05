package pdo

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSharedPoolBorrowReturn(t *testing.T) {
	dsn := "file:pdo_pool_go_test?mode=memory&cache=shared"
	db, err := getSharedDB("sqlite", dsn)
	if err != nil {
		t.Fatal(err)
	}
	db2, err := getSharedDB("sqlite", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if db != db2 {
		t.Fatal("same DSN should reuse shared *sql.DB pool")
	}

	c1, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	c2, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c1.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS t(id INTEGER PRIMARY KEY, v TEXT)"); err != nil {
		t.Fatal(err)
	}
	if _, err := c1.ExecContext(context.Background(), "INSERT INTO t(v) VALUES('x')"); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := c2.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM t").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n < 1 {
		t.Fatalf("count=%d", n)
	}
	// Close 归还到池，不应关掉共享 DB
	if err := c1.Close(); err != nil {
		t.Fatal(err)
	}
	if err := c2.Close(); err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("pool should stay open after Conn.Close: %v", err)
	}
}
