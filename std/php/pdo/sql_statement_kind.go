package pdo

import "strings"

// statementReturnsRows recognizes top-level SQL tokens, excluding literals and
// comments. DML RETURNING produces a result set just like SELECT.
func statementReturnsRows(sql string) bool {
	rows, _ := classifyStatement(sql)
	return rows
}

// classifyStatement also identifies DML RETURNING for driver result metadata.
func classifyStatement(sql string) (rows, returning bool) {
	depth := 0
	first := ""
	for i := 0; i < len(sql); {
		c := sql[i]
		if i+1 < len(sql) && c == '-' && sql[i+1] == '-' {
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			continue
		}
		if i+1 < len(sql) && c == '/' && sql[i+1] == '*' {
			i += 2
			for i+1 < len(sql) && !(sql[i] == '*' && sql[i+1] == '/') {
				i++
			}
			if i+1 < len(sql) {
				i += 2
			}
			continue
		}
		if c == '\'' || c == '"' || c == '`' || c == '[' {
			end := c
			if c == '[' {
				end = ']'
			}
			i++
			for i < len(sql) {
				if sql[i] == end {
					i++
					if i < len(sql) && sql[i] == end {
						i++
						continue
					}
					break
				}
				i++
			}
			continue
		}
		if c == '(' {
			depth++
			i++
			continue
		}
		if c == ')' {
			depth--
			i++
			continue
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' {
			start := i
			for i < len(sql) && ((sql[i] >= 'a' && sql[i] <= 'z') || (sql[i] >= 'A' && sql[i] <= 'Z') || (sql[i] >= '0' && sql[i] <= '9') || sql[i] == '_') {
				i++
			}
			word := strings.ToUpper(sql[start:i])
			if first == "" {
				first = word
				switch first {
				case "SELECT", "PRAGMA", "EXPLAIN", "SHOW":
					return true, false
				}
			}
			if depth == 0 && first == "WITH" {
				switch word {
				case "INSERT", "UPDATE", "DELETE", "REPLACE", "SELECT":
					first = word
				}
			}
			if depth == 0 && word == "RETURNING" && (first == "INSERT" || first == "UPDATE" || first == "DELETE" || first == "REPLACE") {
				return true, true
			}
			continue
		}
		i++
	}
	return first == "WITH" || first == "SELECT", false
}
