package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewRowsClass() data.ClassStmt {
	return &RowsClass{
		source:        nil,
		close:         &RowsCloseMethod{source: nil},
		columnTypes:   &RowsColumnTypesMethod{source: nil},
		columns:       &RowsColumnsMethod{source: nil},
		err:           &RowsErrMethod{source: nil},
		next:          &RowsNextMethod{source: nil},
		nextResultSet: &RowsNextResultSetMethod{source: nil},
		scan:          &RowsScanMethod{source: nil},
	}
}

func NewRowsClassFrom(source *sqlsrc.Rows) data.ClassStmt {
	return &RowsClass{
		source:        source,
		close:         &RowsCloseMethod{source: source},
		columnTypes:   &RowsColumnTypesMethod{source: source},
		columns:       &RowsColumnsMethod{source: source},
		err:           &RowsErrMethod{source: source},
		next:          &RowsNextMethod{source: source},
		nextResultSet: &RowsNextResultSetMethod{source: source},
		scan:          &RowsScanMethod{source: source},
	}
}

type RowsClass struct {
	node.Node
	source        *sqlsrc.Rows
	close         data.Method
	columnTypes   data.Method
	columns       data.Method
	err           data.Method
	next          data.Method
	nextResultSet data.Method
	scan          data.Method
}

func (s *RowsClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *RowsClass) GetName() string                            { return "database\\sql\\Rows" }
func (s *RowsClass) GetExtend() *string                         { return nil }
func (s *RowsClass) GetImplements() []string                    { return nil }
func (s *RowsClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *RowsClass) GetProperties() map[string]data.Property    { return nil }
func (s *RowsClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "close":
		return s.close, true
	case "columnTypes":
		return s.columnTypes, true
	case "columns":
		return s.columns, true
	case "err":
		return s.err, true
	case "next":
		return s.next, true
	case "nextResultSet":
		return s.nextResultSet, true
	case "scan":
		return s.scan, true
	}
	return nil, false
}

func (s *RowsClass) GetMethods() []data.Method {
	return []data.Method{
		s.close,
		s.columnTypes,
		s.columns,
		s.err,
		s.next,
		s.nextResultSet,
		s.scan,
	}
}

func (s *RowsClass) GetConstruct() data.Method { return nil }
