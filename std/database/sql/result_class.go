package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewResultClassFrom(source sqlsrc.Result) data.ClassStmt {
	return &ResultClass{
		source:       source,
		lastInsertId: &ResultLastInsertIdMethod{source: source},
		rowsAffected: &ResultRowsAffectedMethod{source: source},
	}
}

type ResultClass struct {
	node.Node
	source       sqlsrc.Result
	lastInsertId data.Method
	rowsAffected data.Method
}

func (s *ResultClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *ResultClass) GetName() string                            { return "Database\\Sql\\Result" }
func (s *ResultClass) GetExtend() *string                         { return nil }
func (s *ResultClass) GetImplements() []string                    { return nil }
func (s *ResultClass) AsString() string                           { return "Result{}" }
func (s *ResultClass) GetSource() any                             { return s.source }
func (s *ResultClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *ResultClass) GetPropertyList() []data.Property           { return nil }
func (s *ResultClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "lastInsertId":
		return s.lastInsertId, true
	case "rowsAffected":
		return s.rowsAffected, true
	}
	return nil, false
}

func (s *ResultClass) GetMethods() []data.Method {
	return []data.Method{
		s.lastInsertId,
		s.rowsAffected,
	}
}

func (s *ResultClass) GetConstruct() data.Method { return nil }
