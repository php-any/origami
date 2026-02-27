package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewRuntimeExceptionClass() *RuntimeExceptionClass {
	source := &Exception{}

	return &RuntimeExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

type RuntimeExceptionClass struct {
	node.Node
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func (s *RuntimeExceptionClass) AsString() string {
	return s.getMessage.source.msg
}

func (s *RuntimeExceptionClass) IsThrow() bool {
	//TODO implement me
	panic("implement me")
}

func (s *RuntimeExceptionClass) GetError() *data.Error {
	//TODO implement me
	panic("implement me")
}

func (s *RuntimeExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *RuntimeExceptionClass) GetName() string {
	return "RuntimeException"
}

func (s *RuntimeExceptionClass) GetExtend() *string {
	extend := "Exception"
	return &extend
}

func (s *RuntimeExceptionClass) GetImplements() []string {
	return nil
}

func (s *RuntimeExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *RuntimeExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *RuntimeExceptionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return s.exception, true
	case "error":
		return s.error, true
	case "getMessage":
		return s.getMessage, true
	case "getTraceAsString":
		return s.getTraceAsString, true
	}
	return nil, false
}

func (s *RuntimeExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
		s.getTraceAsString,
	}
}

func (s *RuntimeExceptionClass) GetConstruct() data.Method {
	return s.exception
}
