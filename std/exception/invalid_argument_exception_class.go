package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewInvalidArgumentExceptionClass() *InvalidArgumentExceptionClass {
	source := &Exception{}

	return &InvalidArgumentExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

type InvalidArgumentExceptionClass struct {
	node.Node
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func (s *InvalidArgumentExceptionClass) AsString() string {
	return s.getMessage.source.msg
}

func (s *InvalidArgumentExceptionClass) IsThrow() bool {
	//TODO implement me
	panic("implement me")
}

func (s *InvalidArgumentExceptionClass) GetError() *data.Error {
	//TODO implement me
	panic("implement me")
}

func (s *InvalidArgumentExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *InvalidArgumentExceptionClass) GetName() string {
	return "InvalidArgumentException"
}

func (s *InvalidArgumentExceptionClass) GetExtend() *string {
	extend := "LogicException"
	return &extend
}

func (s *InvalidArgumentExceptionClass) GetImplements() []string {
	return nil
}

func (s *InvalidArgumentExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *InvalidArgumentExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *InvalidArgumentExceptionClass) GetMethod(name string) (data.Method, bool) {
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

func (s *InvalidArgumentExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
		s.getTraceAsString,
	}
}

func (s *InvalidArgumentExceptionClass) GetConstruct() data.Method {
	return s.exception
}
