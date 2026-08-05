package pdo

import (
	"database/sql"
	"sync"
)

// 按 driver+DSN 共享 *sql.DB 连接池。
// 单个 PDO 只从池中借出一条 *sql.Conn；close/析构时 Close(conn) 归还到池，不关闭池本身。
var (
	pdoPoolsMu sync.Mutex
	pdoPools   = map[string]*sql.DB{}
)

func poolKey(goDriver, goDSN string) string {
	return goDriver + "\x00" + goDSN
}

func getSharedDB(goDriver, goDSN string) (*sql.DB, error) {
	key := poolKey(goDriver, goDSN)
	pdoPoolsMu.Lock()
	defer pdoPoolsMu.Unlock()
	if db, ok := pdoPools[key]; ok {
		return db, nil
	}
	db, err := sql.Open(goDriver, goDSN)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	pdoPools[key] = db
	return db, nil
}
