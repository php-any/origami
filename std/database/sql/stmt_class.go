package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStmtClass() data.ClassStmt {
	return &StmtClass{
		source:          nil,
		close:           &StmtCloseMethod{source: nil},
		exec:            &StmtExecMethod{source: nil},
		execContext:     &StmtExecContextMethod{source: nil},
		query:           &StmtQueryMethod{source: nil},
		queryContext:    &StmtQueryContextMethod{source: nil},
		queryRow:        &StmtQueryRowMethod{source: nil},
		queryRowContext: &StmtQueryRowContextMethod{source: nil},
	}
}

func NewStmtClassFrom(source *sqlsrc.Stmt) data.ClassStmt {
	return &StmtClass{
		source:          source,
		close:           &StmtCloseMethod{source: source},
		exec:            &StmtExecMethod{source: source},
		execContext:     &StmtExecContextMethod{source: source},
		query:           &StmtQueryMethod{source: source},
		queryContext:    &StmtQueryContextMethod{source: source},
		queryRow:        &StmtQueryRowMethod{source: source},
		queryRowContext: &StmtQueryRowContextMethod{source: source},
	}
}

type StmtClass struct {
	node.Node
	source          *sqlsrc.Stmt
	close           data.Method
	exec            data.Method
	execContext     data.Method
	query           data.Method
	queryContext    data.Method
	queryRow        data.Method
	queryRowContext data.Method
}

func (s *StmtClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *StmtClass) GetName() string                            { return "database\\sql\\Stmt" }
func (s *StmtClass) GetExtend() *string                         { return nil }
func (s *StmtClass) GetImplements() []string                    { return nil }
func (s *StmtClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *StmtClass) GetProperties() map[string]data.Property    { return nil }
func (s *StmtClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "close":
		return s.close, true
	case "exec":
		return s.exec, true
	case "execContext":
		return s.execContext, true
	case "query":
		return s.query, true
	case "queryContext":
		return s.queryContext, true
	case "queryRow":
		return s.queryRow, true
	case "queryRowContext":
		return s.queryRowContext, true
	}
	return nil, false
}

func (s *StmtClass) GetMethods() []data.Method {
	return []data.Method{
		s.close,
		s.exec,
		s.execContext,
		s.query,
		s.queryContext,
		s.queryRow,
		s.queryRowContext,
	}
}

func (s *StmtClass) GetConstruct() data.Method { return nil }
