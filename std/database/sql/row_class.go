package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewRowClass() data.ClassStmt {
	return &RowClass{
		source: nil,
		err:    &RowErrMethod{source: nil},
		scan:   &RowScanMethod{source: nil},
	}
}

func NewRowClassFrom(source *sqlsrc.Row) data.ClassStmt {
	return &RowClass{
		source: source,
		err:    &RowErrMethod{source: source},
		scan:   &RowScanMethod{source: source},
	}
}

type RowClass struct {
	node.Node
	source *sqlsrc.Row
	err    data.Method
	scan   data.Method
}

func (s *RowClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *RowClass) GetName() string                            { return "Database\\Sql\\Row" }
func (s *RowClass) GetExtend() *string                         { return nil }
func (s *RowClass) GetImplements() []string                    { return nil }
func (s *RowClass) AsString() string                           { return "Row{}" }
func (s *RowClass) GetSource() any                             { return s.source }
func (s *RowClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *RowClass) GetPropertyList() []data.Property           { return nil }
func (s *RowClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "err":
		return s.err, true
	case "scan":
		return s.scan, true
	}
	return nil, false
}

func (s *RowClass) GetMethods() []data.Method {
	return []data.Method{
		s.err,
		s.scan,
	}
}

func (s *RowClass) GetConstruct() data.Method { return nil }
