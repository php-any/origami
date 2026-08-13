package migration

import (
	"strings"
	"testing"
)

func TestSQLiteCreateTableSQL(t *testing.T) {
	d := sqliteDialect{}
	sql := d.BuildCreateTableSQL(&TableSchema{
		TableName: "users",
		Columns: []ColumnSchema{
			{Name: "id", PHPType: "int", Nullable: false, PrimaryKey: true, AutoIncrement: true},
			{Name: "name", PHPType: "string", Length: 100, Nullable: false},
		},
	})

	if !strings.Contains(sql, `"id" INTEGER PRIMARY KEY AUTOINCREMENT`) {
		t.Fatalf("unexpected sqlite id column: %s", sql)
	}
	if !strings.Contains(sql, `"name" VARCHAR(100) NOT NULL`) {
		t.Fatalf("unexpected sqlite name column: %s", sql)
	}
}

func TestMySQLCreateTableSQL(t *testing.T) {
	d := mysqlDialect{}
	sql := d.BuildCreateTableSQL(&TableSchema{
		TableName: "users",
		Columns: []ColumnSchema{
			{Name: "id", PHPType: "int", Nullable: false, PrimaryKey: true, AutoIncrement: true},
			{Name: "price", PHPType: "float", Nullable: false},
			{Name: "active", PHPType: "bool", Nullable: false},
			{Name: "name", PHPType: "string", Length: 100, Nullable: false},
		},
	})

	if !strings.Contains(sql, "`id` INT AUTO_INCREMENT PRIMARY KEY") {
		t.Fatalf("unexpected mysql id column: %s", sql)
	}
	if !strings.Contains(sql, "`price` DOUBLE NOT NULL") {
		t.Fatalf("unexpected mysql price column: %s", sql)
	}
	if !strings.Contains(sql, "`active` TINYINT(1) NOT NULL") {
		t.Fatalf("unexpected mysql active column: %s", sql)
	}
}

func TestDialectAddColumnSQL(t *testing.T) {
	col := ColumnSchema{Name: "email", PHPType: "string", Length: 100, Nullable: false}

	sqliteSQL := sqliteDialect{}.BuildAddColumnSQL("users", col)
	if sqliteSQL != `ALTER TABLE "users" ADD COLUMN "email" VARCHAR(100) NOT NULL` {
		t.Fatalf("unexpected sqlite add column: %s", sqliteSQL)
	}

	mysqlSQL := mysqlDialect{}.BuildAddColumnSQL("users", col)
	if mysqlSQL != "ALTER TABLE `users` ADD COLUMN `email` VARCHAR(100) NOT NULL" {
		t.Fatalf("unexpected mysql add column: %s", mysqlSQL)
	}
}
