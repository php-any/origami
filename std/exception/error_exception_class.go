package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewErrorExceptionClass() *ErrorExceptionClass {
	source := &Exception{}

	return &ErrorExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

type ErrorExceptionClass struct {
	node.Node
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func (s *ErrorExceptionClass) AsString() string {
	return s.getMessage.source.msg
}

func (s *ErrorExceptionClass) IsThrow() bool {
	panic("implement me")
}

func (s *ErrorExceptionClass) GetError() *data.Error {
	panic("implement me")
}

func (s *ErrorExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	newS := *s
	return data.NewClassValue(&newS, ctx), nil
}

func (s *ErrorExceptionClass) GetName() string {
	return "ErrorException"
}

func (s *ErrorExceptionClass) GetExtend() *string {
	name := "Exception"
	return &name
}

func (s *ErrorExceptionClass) GetImplements() []string {
	return []string{"Throwable"}
}

func (s *ErrorExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *ErrorExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *ErrorExceptionClass) GetMethod(name string) (data.Method, bool) {
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

func (s *ErrorExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
		s.getTraceAsString,
	}
}

func (s *ErrorExceptionClass) GetConstruct() data.Method {
	return s.exception
}
