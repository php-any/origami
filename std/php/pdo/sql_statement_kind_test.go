package pdo

import "testing"

func TestStatementReturnsRows(t *testing.T) {
	cases := []struct {
		sql  string
		want bool
	}{
		{"-- note\nSELECT 1", true},
		{"INSERT INTO t(v) VALUES(?) RETURNING id", true},
		{"UPDATE t SET v=? returning id", true},
		{"DELETE FROM t RETURNING *", true},
		{"INSERT INTO t(v) VALUES('RETURNING')", false},
		{"UPDATE t SET v='it''s RETURNING' /* RETURNING id */", false},
		{"CREATE TABLE t(RETURNING_VALUE TEXT)", false},
		{"INSERT INTO t(\"RETURNING\") VALUES(?)", false},
		{"INSERT INTO t(v) VALUES(?) -- RETURNING id", false},
	}
	for _, c := range cases {
		if got := statementReturnsRows(c.sql); got != c.want {
			t.Errorf("%q: got %v want %v", c.sql, got, c.want)
		}
	}
}

func TestCTEReturningClassification(t *testing.T) {
	rows, returning := classifyStatement("WITH source(v) AS (SELECT 'RETURNING') INSERT INTO t(v) SELECT v FROM source RETURNING id")
	if !rows || !returning {
		t.Fatal("CTE DML RETURNING was not classified")
	}
	rows, returning = classifyStatement("WITH source(v) AS (SELECT 'RETURNING') SELECT v FROM source")
	if !rows || returning {
		t.Fatal("CTE SELECT was misclassified")
	}
}
