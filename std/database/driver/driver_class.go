package driver

import (
	driversrc "database/sql/driver"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewDriverClassFrom(source driversrc.Driver) data.ClassStmt {
	return &DriverClass{
		source: source,
		open:   &DriverOpenMethod{source: source},
	}
}

type DriverClass struct {
	node.Node
	source driversrc.Driver
	open   data.Method
}

func (s *DriverClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *DriverClass) GetName() string                            { return "database\\sql\\Driver" }
func (s *DriverClass) GetExtend() *string                         { return nil }
func (s *DriverClass) GetImplements() []string                    { return nil }
func (s *DriverClass) AsString() string                           { return "Driver{}" }
func (s *DriverClass) GetSource() any                             { return s.source }
func (s *DriverClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *DriverClass) GetProperties() map[string]data.Property    { return nil }
func (s *DriverClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "open":
		return s.open, true
	}
	return nil, false
}

func (s *DriverClass) GetMethods() []data.Method {
	return []data.Method{
		s.open,
	}
}

func (s *DriverClass) GetConstruct() data.Method { return nil }
