package database

import (
	"database/sql"

	"github.com/php-any/origami/data"
)

func newDB(genericMap map[string]data.Types) *db {
	return &db{
		model: genericMap["M"],
	}
}

type db struct {
	conn *sql.DB // 由构造函数初始化

	where     string
	whereArgs []data.Value

	// 泛型对应的具体类型
	model data.Types
}
