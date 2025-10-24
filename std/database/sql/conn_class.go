package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewConnClass() data.ClassStmt {
	return &ConnClass{
		source:          nil,
		beginTx:         &ConnBeginTxMethod{source: nil},
		close:           &ConnCloseMethod{source: nil},
		execContext:     &ConnExecContextMethod{source: nil},
		pingContext:     &ConnPingContextMethod{source: nil},
		prepareContext:  &ConnPrepareContextMethod{source: nil},
		queryContext:    &ConnQueryContextMethod{source: nil},
		queryRowContext: &ConnQueryRowContextMethod{source: nil},
		raw:             &ConnRawMethod{source: nil},
	}
}

func NewConnClassFrom(source *sqlsrc.Conn) data.ClassStmt {
	return &ConnClass{
		source:          source,
		beginTx:         &ConnBeginTxMethod{source: source},
		close:           &ConnCloseMethod{source: source},
		execContext:     &ConnExecContextMethod{source: source},
		pingContext:     &ConnPingContextMethod{source: source},
		prepareContext:  &ConnPrepareContextMethod{source: source},
		queryContext:    &ConnQueryContextMethod{source: source},
		queryRowContext: &ConnQueryRowContextMethod{source: source},
		raw:             &ConnRawMethod{source: source},
	}
}

type ConnClass struct {
	node.Node
	source          *sqlsrc.Conn
	beginTx         data.Method
	close           data.Method
	execContext     data.Method
	pingContext     data.Method
	prepareContext  data.Method
	queryContext    data.Method
	queryRowContext data.Method
	raw             data.Method
}

func (s *ConnClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *ConnClass) GetName() string                            { return "Database\\Sql\\Conn" }
func (s *ConnClass) GetExtend() *string                         { return nil }
func (s *ConnClass) GetImplements() []string                    { return nil }
func (s *ConnClass) AsString() string                           { return "Conn{}" }
func (s *ConnClass) GetSource() any                             { return s.source }
func (s *ConnClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *ConnClass) GetPropertyList() []data.Property           { return nil }
func (s *ConnClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "beginTx":
		return s.beginTx, true
	case "close":
		return s.close, true
	case "execContext":
		return s.execContext, true
	case "pingContext":
		return s.pingContext, true
	case "prepareContext":
		return s.prepareContext, true
	case "queryContext":
		return s.queryContext, true
	case "queryRowContext":
		return s.queryRowContext, true
	case "raw":
		return s.raw, true
	}
	return nil, false
}

func (s *ConnClass) GetMethods() []data.Method {
	return []data.Method{
		s.beginTx,
		s.close,
		s.execContext,
		s.pingContext,
		s.prepareContext,
		s.queryContext,
		s.queryRowContext,
		s.raw,
	}
}

func (s *ConnClass) GetConstruct() data.Method { return nil }
