module github.com/php-any/origami/examples/database

go 1.24

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/php-any/origami v0.0.0
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/php-any/origami => ../../
