package migration

import (
	"database/sql"
	"fmt"
	"strings"
)

// execQuerier 供方言查询元数据及执行 DDL（兼容 *sql.DB / *sql.Tx）
type execQuerier interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// Dialect 数据库方言：元数据查询与 DDL 生成
type Dialect interface {
	Name() string
	TableExists(q execQuerier, tableName string) (bool, error)
	TableColumns(q execQuerier, tableName string) (map[string]bool, error)
	BuildCreateTableSQL(schema *TableSchema) string
	BuildAddColumnSQL(tableName string, col ColumnSchema) string
}

func detectDialect(conn *sql.DB) (Dialect, error) {
	if conn == nil {
		return nil, fmt.Errorf("数据库连接为空")
	}
	driverName := strings.ToLower(fmt.Sprintf("%T", conn.Driver()))
	switch {
	case strings.Contains(driverName, "sqlite"):
		return sqliteDialect{}, nil
	case strings.Contains(driverName, "mysql"):
		return mysqlDialect{}, nil
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", driverName)
	}
}

func quoteBacktick(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func quoteDouble(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
