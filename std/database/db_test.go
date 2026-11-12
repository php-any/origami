package database

import (
	"database/sql"
	"testing"
)

func TestGetConnectionPrefersExplicitExecutor(t *testing.T) {
	pool := &sql.DB{}
	tx := &sql.Tx{}
	query := &db{conn: pool, exec: tx}
	if query.getConnection() != tx {
		t.Fatal("explicit transaction was not used")
	}
	query.exec = nil
	if query.getConnection() != pool {
		t.Fatal("pool was not used")
	}
	query.setConnectionName("other")
	if query.exec != nil || query.conn != nil {
		t.Fatal("connection name change kept the previous executor")
	}
}
