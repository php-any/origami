package driver

import (
	driversrc "database/sql/driver"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewConnClassFrom(source driversrc.Conn) data.ClassStmt {
	return &ConnClass{
		source:  source,
		begin:   &ConnBeginMethod{source: source},
		close:   &ConnCloseMethod{source: source},
		prepare: &ConnPrepareMethod{source: source},
	}
}

type ConnClass struct {
	node.Node
	source  driversrc.Conn
	begin   data.Method
	close   data.Method
	prepare data.Method
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
	case "begin":
		return s.begin, true
	case "close":
		return s.close, true
	case "prepare":
		return s.prepare, true
	}
	return nil, false
}

func (s *ConnClass) GetMethods() []data.Method {
	return []data.Method{
		s.begin,
		s.close,
		s.prepare,
	}
}

func (s *ConnClass) GetConstruct() data.Method { return nil }
