package driver

import (
	driversrc "database/sql/driver"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewRowsClassFrom(source driversrc.Rows) data.ClassStmt {
	return &RowsClass{
		source:  source,
		close:   &RowsCloseMethod{source: source},
		columns: &RowsColumnsMethod{source: source},
		next:    &RowsNextMethod{source: source},
	}
}

type RowsClass struct {
	node.Node
	source  driversrc.Rows
	close   data.Method
	columns data.Method
	next    data.Method
}

func (s *RowsClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *RowsClass) GetName() string                            { return "Database\\Sql\\Rows" }
func (s *RowsClass) GetExtend() *string                         { return nil }
func (s *RowsClass) GetImplements() []string                    { return nil }
func (s *RowsClass) AsString() string                           { return "Rows{}" }
func (s *RowsClass) GetSource() any                             { return s.source }
func (s *RowsClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *RowsClass) GetPropertyList() []data.Property           { return nil }
func (s *RowsClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "close":
		return s.close, true
	case "columns":
		return s.columns, true
	case "next":
		return s.next, true
	}
	return nil, false
}

func (s *RowsClass) GetMethods() []data.Method {
	return []data.Method{
		s.close,
		s.columns,
		s.next,
	}
}

func (s *RowsClass) GetConstruct() data.Method { return nil }
