package migration

import (
	"fmt"
	"strings"
)

type mysqlDialect struct{}

func (mysqlDialect) Name() string { return "mysql" }

func (mysqlDialect) TableExists(q execQuerier, tableName string) (bool, error) {
	var count int
	err := q.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`,
		tableName,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (mysqlDialect) TableColumns(q execQuerier, tableName string) (map[string]bool, error) {
	rows, err := q.Query(`
		SELECT COLUMN_NAME
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`,
		tableName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		cols[name] = true
	}
	return cols, rows.Err()
}

func (d mysqlDialect) BuildCreateTableSQL(schema *TableSchema) string {
	parts := make([]string, 0, len(schema.Columns))
	for _, col := range schema.Columns {
		parts = append(parts, d.columnDefinition(col))
	}
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n    %s\n)",
		quoteBacktick(schema.TableName),
		strings.Join(parts, ",\n    "),
	)
}

func (d mysqlDialect) BuildAddColumnSQL(tableName string, col ColumnSchema) string {
	return fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN %s",
		quoteBacktick(tableName),
		d.columnDefinition(col),
	)
}

func (mysqlDialect) columnDefinition(col ColumnSchema) string {
	def := fmt.Sprintf("%s %s", quoteBacktick(col.Name), mysqlSQLType(col))

	if col.AutoIncrement {
		def += " AUTO_INCREMENT"
	}

	if col.PrimaryKey {
		def += " PRIMARY KEY"
	}

	if !col.Nullable && !col.PrimaryKey {
		def += " NOT NULL"
	}
	return def
}

func mysqlSQLType(col ColumnSchema) string {
	switch strings.ToLower(col.PHPType) {
	case "int", "integer":
		return "INT"
	case "float", "double", "real":
		return "DOUBLE"
	case "bool", "boolean":
		return "TINYINT(1)"
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
