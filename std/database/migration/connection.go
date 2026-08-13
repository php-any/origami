package migration

import (
	"database/sql"

	"github.com/php-any/origami/data"
	sqlpkg "github.com/php-any/origami/std/database/sql"
)

func sqlDBFromValue(value data.GetValue) (*sql.DB, bool) {
	if dbClass, ok := value.(*data.ClassValue); ok {
		if dbClassStmt, ok := dbClass.Class.(*sqlpkg.DBClass); ok {
			if sqlDB, ok := dbClassStmt.GetSource().(*sql.DB); ok {
				return sqlDB, true
			}
		}
	}
	if anyVal, ok := value.(*data.AnyValue); ok {
		if sqlDB, ok := anyVal.Value.(*sql.DB); ok {
			return sqlDB, true
		}
	}
	return nil, false
}
