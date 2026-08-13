package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewDBClass() data.ClassStmt {
	return &DBClass{
		source:             nil,
		begin:              &DBBeginMethod{source: nil},
		beginTx:            &DBBeginTxMethod{source: nil},
		close:              &DBCloseMethod{source: nil},
		conn:               &DBConnMethod{source: nil},
		driver:             &DBDriverMethod{source: nil},
		exec:               &DBExecMethod{source: nil},
		execContext:        &DBExecContextMethod{source: nil},
		ping:               &DBPingMethod{source: nil},
		pingContext:        &DBPingContextMethod{source: nil},
		prepare:            &DBPrepareMethod{source: nil},
		prepareContext:     &DBPrepareContextMethod{source: nil},
		query:              &DBQueryMethod{source: nil},
		queryContext:       &DBQueryContextMethod{source: nil},
		queryRow:           &DBQueryRowMethod{source: nil},
		queryRowContext:    &DBQueryRowContextMethod{source: nil},
		setConnMaxIdleTime: &DBSetConnMaxIdleTimeMethod{source: nil},
		setConnMaxLifetime: &DBSetConnMaxLifetimeMethod{source: nil},
		setMaxIdleConns:    &DBSetMaxIdleConnsMethod{source: nil},
		setMaxOpenConns:    &DBSetMaxOpenConnsMethod{source: nil},
		stats:              &DBStatsMethod{source: nil},
	}
}

func NewDBClassFrom(source *sqlsrc.DB) data.ClassStmt {
	return &DBClass{
		source:             source,
		begin:              &DBBeginMethod{source: source},
		beginTx:            &DBBeginTxMethod{source: source},
		close:              &DBCloseMethod{source: source},
		conn:               &DBConnMethod{source: source},
		driver:             &DBDriverMethod{source: source},
		exec:               &DBExecMethod{source: source},
		execContext:        &DBExecContextMethod{source: source},
		ping:               &DBPingMethod{source: source},
		pingContext:        &DBPingContextMethod{source: source},
		prepare:            &DBPrepareMethod{source: source},
		prepareContext:     &DBPrepareContextMethod{source: source},
		query:              &DBQueryMethod{source: source},
		queryContext:       &DBQueryContextMethod{source: source},
		queryRow:           &DBQueryRowMethod{source: source},
		queryRowContext:    &DBQueryRowContextMethod{source: source},
		setConnMaxIdleTime: &DBSetConnMaxIdleTimeMethod{source: source},
		setConnMaxLifetime: &DBSetConnMaxLifetimeMethod{source: source},
		setMaxIdleConns:    &DBSetMaxIdleConnsMethod{source: source},
		setMaxOpenConns:    &DBSetMaxOpenConnsMethod{source: source},
		stats:              &DBStatsMethod{source: source},
	}
}

type DBClass struct {
	node.Node
	source             *sqlsrc.DB
	begin              data.Method
	beginTx            data.Method
	close              data.Method
	conn               data.Method
	driver             data.Method
	exec               data.Method
	execContext        data.Method
	ping               data.Method
	pingContext        data.Method
	prepare            data.Method
	prepareContext     data.Method
	query              data.Method
	queryContext       data.Method
	queryRow           data.Method
	queryRowContext    data.Method
	setConnMaxIdleTime data.Method
	setConnMaxLifetime data.Method
	setMaxIdleConns    data.Method
	setMaxOpenConns    data.Method
	stats              data.Method
}

func (s *DBClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *DBClass) GetName() string                            { return "Database\\Sql\\DB" }
func (s *DBClass) GetExtend() *string                         { return nil }
func (s *DBClass) GetImplements() []string                    { return nil }
func (s *DBClass) AsString() string                           { return "DB{}" }
func (s *DBClass) GetSource() any                             { return s.source }
func (s *DBClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *DBClass) GetPropertyList() []data.Property           { return nil }
func (s *DBClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "begin":
		return s.begin, true
	case "beginTx":
		return s.beginTx, true
	case "close":
		return s.close, true
	case "conn":
		return s.conn, true
	case "driver":
		return s.driver, true
	case "exec":
		return s.exec, true
	case "execContext":
		return s.execContext, true
	case "ping":
		return s.ping, true
	case "pingContext":
		return s.pingContext, true
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
	case "setConnMaxIdleTime":
		return s.setConnMaxIdleTime, true
	case "setConnMaxLifetime":
		return s.setConnMaxLifetime, true
	case "setMaxIdleConns":
		return s.setMaxIdleConns, true
	case "setMaxOpenConns":
		return s.setMaxOpenConns, true
	case "stats":
		return s.stats, true
	}
	return nil, false
}

func (s *DBClass) GetMethods() []data.Method {
	return []data.Method{
		s.begin,
		s.beginTx,
		s.close,
		s.conn,
		s.driver,
		s.exec,
		s.execContext,
		s.ping,
		s.pingContext,
		s.prepare,
		s.prepareContext,
		s.query,
		s.queryContext,
		s.queryRow,
		s.queryRowContext,
		s.setConnMaxIdleTime,
		s.setConnMaxLifetime,
		s.setMaxIdleConns,
		s.setMaxOpenConns,
		s.stats,
	}
}

func (s *DBClass) GetConstruct() data.Method { return nil }
