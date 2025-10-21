package driver

import (
	driversrc "database/sql/driver"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStmtClassFrom(source driversrc.Stmt) data.ClassStmt {
	return &StmtClass{
		source:   source,
		close:    &StmtCloseMethod{source: source},
		exec:     &StmtExecMethod{source: source},
		numInput: &StmtNumInputMethod{source: source},
		query:    &StmtQueryMethod{source: source},
	}
}

type StmtClass struct {
	node.Node
	source   driversrc.Stmt
	close    data.Method
	exec     data.Method
	numInput data.Method
	query    data.Method
}

func (s *StmtClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *StmtClass) GetName() string                            { return "database\\sql\\Stmt" }
func (s *StmtClass) GetExtend() *string                         { return nil }
func (s *StmtClass) GetImplements() []string                    { return nil }
func (s *StmtClass) AsString() string                           { return "Stmt{}" }
func (s *StmtClass) GetSource() any                             { return s.source }
func (s *StmtClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *StmtClass) GetProperties() map[string]data.Property    { return nil }
func (s *StmtClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "close":
		return s.close, true
	case "exec":
		return s.exec, true
	case "numInput":
		return s.numInput, true
	case "query":
		return s.query, true
	}
	return nil, false
}

func (s *StmtClass) GetMethods() []data.Method {
	return []data.Method{
		s.close,
		s.exec,
		s.numInput,
		s.query,
	}
}

func (s *StmtClass) GetConstruct() data.Method { return nil }
