package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewUnexpectedValueExceptionClass() *UnexpectedValueExceptionClass {
	source := &Exception{}

	return &UnexpectedValueExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

type UnexpectedValueExceptionClass struct {
	node.Node
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func (s *UnexpectedValueExceptionClass) AsString() string {
	return s.getMessage.source.msg
}

func (s *UnexpectedValueExceptionClass) IsThrow() bool {
	panic("implement me")
}

func (s *UnexpectedValueExceptionClass) GetError() *data.Error {
	panic("implement me")
}

func (s *UnexpectedValueExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *UnexpectedValueExceptionClass) GetName() string {
	return "UnexpectedValueException"
}

func (s *UnexpectedValueExceptionClass) GetExtend() *string {
	extend := "RuntimeException"
	return &extend
}

func (s *UnexpectedValueExceptionClass) GetImplements() []string {
	return nil
}

func (s *UnexpectedValueExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *UnexpectedValueExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *UnexpectedValueExceptionClass) GetMethod(name string) (data.Method, bool) {
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

func (s *UnexpectedValueExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
		s.getTraceAsString,
	}
}

func (s *UnexpectedValueExceptionClass) GetConstruct() data.Method {
	return s.exception
}
