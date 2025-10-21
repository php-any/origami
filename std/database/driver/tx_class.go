package driver

import (
	driversrc "database/sql/driver"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewTxClassFrom(source driversrc.Tx) data.ClassStmt {
	return &TxClass{
		source:   source,
		commit:   &TxCommitMethod{source: source},
		rollback: &TxRollbackMethod{source: source},
	}
}

type TxClass struct {
	node.Node
	source   driversrc.Tx
	commit   data.Method
	rollback data.Method
}

func (s *TxClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *TxClass) GetName() string                            { return "database\\sql\\Tx" }
func (s *TxClass) GetExtend() *string                         { return nil }
func (s *TxClass) GetImplements() []string                    { return nil }
func (s *TxClass) AsString() string                           { return "Tx{}" }
func (s *TxClass) GetSource() any                             { return s.source }
func (s *TxClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *TxClass) GetProperties() map[string]data.Property    { return nil }
func (s *TxClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "commit":
		return s.commit, true
	case "rollback":
		return s.rollback, true
	}
	return nil, false
}

func (s *TxClass) GetMethods() []data.Method {
	return []data.Method{
		s.commit,
		s.rollback,
	}
}

func (s *TxClass) GetConstruct() data.Method { return nil }
