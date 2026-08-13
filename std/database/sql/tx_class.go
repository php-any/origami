package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewTxClass() data.ClassStmt {
	return &TxClass{
		source:          nil,
		commit:          &TxCommitMethod{source: nil},
		exec:            &TxExecMethod{source: nil},
		execContext:     &TxExecContextMethod{source: nil},
		prepare:         &TxPrepareMethod{source: nil},
		prepareContext:  &TxPrepareContextMethod{source: nil},
		query:           &TxQueryMethod{source: nil},
		queryContext:    &TxQueryContextMethod{source: nil},
		queryRow:        &TxQueryRowMethod{source: nil},
		queryRowContext: &TxQueryRowContextMethod{source: nil},
		rollback:        &TxRollbackMethod{source: nil},
		stmt:            &TxStmtMethod{source: nil},
		stmtContext:     &TxStmtContextMethod{source: nil},
	}
}

func NewTxClassFrom(source *sqlsrc.Tx) data.ClassStmt {
	return &TxClass{
		source:          source,
		commit:          &TxCommitMethod{source: source},
		exec:            &TxExecMethod{source: source},
		execContext:     &TxExecContextMethod{source: source},
		prepare:         &TxPrepareMethod{source: source},
		prepareContext:  &TxPrepareContextMethod{source: source},
		query:           &TxQueryMethod{source: source},
		queryContext:    &TxQueryContextMethod{source: source},
		queryRow:        &TxQueryRowMethod{source: source},
		queryRowContext: &TxQueryRowContextMethod{source: source},
		rollback:        &TxRollbackMethod{source: source},
		stmt:            &TxStmtMethod{source: source},
		stmtContext:     &TxStmtContextMethod{source: source},
	}
}

type TxClass struct {
	node.Node
	source          *sqlsrc.Tx
	commit          data.Method
	exec            data.Method
	execContext     data.Method
	prepare         data.Method
	prepareContext  data.Method
	query           data.Method
	queryContext    data.Method
	queryRow        data.Method
	queryRowContext data.Method
	rollback        data.Method
	stmt            data.Method
	stmtContext     data.Method
}

func (s *TxClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *TxClass) GetName() string                            { return "Database\\Sql\\Tx" }
func (s *TxClass) GetExtend() *string                         { return nil }
func (s *TxClass) GetImplements() []string                    { return nil }
func (s *TxClass) AsString() string                           { return "Tx{}" }
func (s *TxClass) GetSource() any                             { return s.source }
func (s *TxClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *TxClass) GetPropertyList() []data.Property           { return nil }
func (s *TxClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "commit":
		return s.commit, true
	case "exec":
		return s.exec, true
	case "execContext":
		return s.execContext, true
	case "prepare":
		return s.prepare, true
	case "prepareContext":
		return s.prepareContext, true
	case "query":
		return s.query, true
	case "queryContext":
		return s.queryContext, true
	case "queryRow":
		return s.queryRow, true
	case "queryRowContext":
		return s.queryRowContext, true
	case "rollback":
		return s.rollback, true
	case "stmt":
		return s.stmt, true
	case "stmtContext":
		return s.stmtContext, true
	}
	return nil, false
}

func (s *TxClass) GetMethods() []data.Method {
	return []data.Method{
		s.commit,
		s.exec,
		s.execContext,
		s.prepare,
		s.prepareContext,
		s.query,
		s.queryContext,
		s.queryRow,
		s.queryRowContext,
		s.rollback,
		s.stmt,
		s.stmtContext,
	}
}

func (s *TxClass) GetConstruct() data.Method { return nil }
