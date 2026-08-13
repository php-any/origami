package migration

import (
	"database/sql"
	"fmt"
	"strings"
)

type sqliteDialect struct{}

func (sqliteDialect) Name() string { return "sqlite" }

func (sqliteDialect) TableExists(q execQuerier, tableName string) (bool, error) {
	var name string
	err := q.QueryRow(
		"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
		tableName,
	).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return name != "", nil
}

func (sqliteDialect) TableColumns(q execQuerier, tableName string) (map[string]bool, error) {
	rows, err := q.Query("PRAGMA table_info(" + quoteDouble(tableName) + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return nil, err
		}
		cols[name] = true
	}
	return cols, rows.Err()
}

func (d sqliteDialect) BuildCreateTableSQL(schema *TableSchema) string {
	parts := make([]string, 0, len(schema.Columns))
	for _, col := range schema.Columns {
		parts = append(parts, d.columnDefinition(col))
	}
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n    %s\n)",
		quoteDouble(schema.TableName),
		strings.Join(parts, ",\n    "),
	)
}

func (d sqliteDialect) BuildAddColumnSQL(tableName string, col ColumnSchema) string {
	return fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN %s",
		quoteDouble(tableName),
		d.columnDefinition(col),
	)
}

func (sqliteDialect) columnDefinition(col ColumnSchema) string {
	def := fmt.Sprintf("%s %s", quoteDouble(col.Name), sqliteSQLType(col))

	if col.PrimaryKey && col.AutoIncrement {
		def += " PRIMARY KEY AUTOINCREMENT"
	} else if col.PrimaryKey {
		def += " PRIMARY KEY"
	}

	if !col.Nullable && !col.PrimaryKey {
		def += " NOT NULL"
	}
	return def
}

func sqliteSQLType(col ColumnSchema) string {
	switch strings.ToLower(col.PHPType) {
	case "int", "integer":
		return "INTEGER"
	case "float", "double", "real":
		return "REAL"
	case "bool", "boolean":
		return "INTEGER"
	case "string":
		length := col.Length
		if length <= 0 {
			length = 255
		}
		if length > 65535 {
			return "TEXT"
		}
		return fmt.Sprintf("VARCHAR(%d)", length)
	default:
		return "TEXT"
	}
}
