package system

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DateTimeClass struct {
	node.Node
	format       data.Method
	getTimestamp data.Method
}

func (s *DateTimeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newDateTime()

	return data.NewClassValue(&DateTimeClass{
		format:       &DateTimeFormatMethod{source},
		getTimestamp: &DateTimeGetTimestampMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (s *DateTimeClass) GetName() string {
	return "System\\DateTime"
}

func (s *DateTimeClass) GetExtend() *string {
	return nil
}

func (s *DateTimeClass) GetImplements() []string {
	return nil
}

func (s *DateTimeClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *DateTimeClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *DateTimeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "format":
		return s.format, true
	case "getTimestamp":
		return s.getTimestamp, true
	}
	return nil, false
}

func (s *DateTimeClass) GetMethods() []data.Method {
	return []data.Method{
		s.format,
		s.getTimestamp,
	}
}

func (t *DateTimeClass) GetConstruct() data.Method {
	return nil
}
