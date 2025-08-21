package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewTxOptionsClass() data.ClassStmt {
	return &TxOptionsClass{
		source:        nil,
		propIsolation: node.NewProperty(nil, "Isolation", "public", true, nil),
		propReadOnly:  node.NewProperty(nil, "ReadOnly", "public", true, nil),
	}
}

func NewTxOptionsClassFrom(source *sqlsrc.TxOptions) data.ClassStmt {
	return &TxOptionsClass{
		source:        source,
		propIsolation: node.NewProperty(nil, "Isolation", "public", true, data.NewAnyValue(source.Isolation)),
		propReadOnly:  node.NewProperty(nil, "ReadOnly", "public", true, data.NewAnyValue(source.ReadOnly)),
	}
}

type TxOptionsClass struct {
	node.Node
	source        *sqlsrc.TxOptions
	propIsolation data.Property
	propReadOnly  data.Property
}

func (s *TxOptionsClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *TxOptionsClass) GetName() string         { return "database\\sql\\TxOptions" }
func (s *TxOptionsClass) GetExtend() *string      { return nil }
func (s *TxOptionsClass) GetImplements() []string { return nil }
func (s *TxOptionsClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "Isolation":
		return s.propIsolation, true
	case "ReadOnly":
		return s.propReadOnly, true
	}
	return nil, false
}

func (s *TxOptionsClass) GetProperties() map[string]data.Property {
	return map[string]data.Property{
		"Isolation": s.propIsolation,
		"ReadOnly":  s.propReadOnly,
	}
}

func (s *TxOptionsClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	}
	return nil, false
}

func (s *TxOptionsClass) GetMethods() []data.Method {
	return []data.Method{}
}

func (s *TxOptionsClass) GetConstruct() data.Method { return nil }
