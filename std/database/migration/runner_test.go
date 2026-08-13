package migration

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestApplyMigrationsSQLiteTransactional(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	dialect := sqliteDialect{}
	schemas := []*TableSchema{
		{
			TableName: "users",
			ClassName: "App\\UserEntity",
			Columns: []ColumnSchema{
				{Name: "id", PHPType: "int", Nullable: false, PrimaryKey: true, AutoIncrement: true},
				{Name: "name", PHPType: "string", Length: 100, Nullable: false},
			},
		},
		{
			TableName: "posts",
			ClassName: "App\\PostEntity",
			Columns: []ColumnSchema{
				{Name: "id", PHPType: "int", Nullable: false, PrimaryKey: true, AutoIncrement: true},
				{Name: "title", PHPType: "string", Length: 200, Nullable: false},
			},
		},
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	created, altered, err := applyMigrations(tx, dialect, schemas)
	if err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if len(created) != 2 || len(altered) != 0 {
		t.Fatalf("first run created=%d altered=%d", len(created), len(altered))
	}

	// 模拟新增列
	schemas[0].Columns = append(schemas[0].Columns, ColumnSchema{
		Name: "email", PHPType: "string", Length: 100, Nullable: false,
	})

	tx, err = db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	_, altered, err = applyMigrations(tx, dialect, schemas)
	if err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if len(altered) != 1 || altered[0].column != "email" {
		t.Fatalf("second run altered=%v", altered)
	}

	// 幂等
	tx, err = db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	created, altered, err = applyMigrations(tx, dialect, schemas)
	if err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if len(created) != 0 || len(altered) != 0 {
		t.Fatalf("third run created=%d altered=%d", len(created), len(altered))
	}
}

func TestApplyMigrationsRollbackOnFailure(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	dialect := sqliteDialect{}
	schemas := []*TableSchema{
		{
			TableName: "users",
			ClassName: "App\\UserEntity",
			Columns: []ColumnSchema{
				{Name: "id", PHPType: "int", Nullable: false, PrimaryKey: true, AutoIncrement: true},
			},
		},
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = applyMigrations(tx, dialect, schemas)
	if err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}

	if _, err := tx.Exec("THIS IS NOT VALID SQL"); err == nil {
		_ = tx.Commit()
		t.Fatal("expected invalid SQL to fail")
	}
	_ = tx.Rollback()

	exists, err := dialect.TableExists(db, "users")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("rolled back transaction should not leave created table")
	}
}

func TestDetectDialect(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	dialect, err := detectDialect(db)
	if err != nil {
		t.Fatal(err)
	}
	if dialect.Name() != "sqlite" {
		t.Fatalf("expected sqlite dialect, got %s", dialect.Name())
	}
}
