package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewExceptionClass() *ExceptionClass {
	source := &Exception{}

	return &ExceptionClass{
		exception:  &ExceptionExceptionMethod{source},
		error:      &ExceptionErrorMethod{source},
		getMessage: &ExceptionGetMessageMethod{source},
	}
}

type ExceptionClass struct {
	node.Node
	exception  data.Method
	error      data.Method
	getMessage data.Method
}

func (s *ExceptionClass) AsString() string {
	//TODO implement me
	panic("implement me")
}

func (s *ExceptionClass) IsThrow() bool {
	//TODO implement me
	panic("implement me")
}

func (s *ExceptionClass) GetError() *data.Error {
	//TODO implement me
	panic("implement me")
}

func (s *ExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *ExceptionClass) GetName() string {
	return "Exception"
}

func (s *ExceptionClass) GetExtend() *string {
	return nil
}

func (s *ExceptionClass) GetImplements() []string {
	return nil
}

func (s *ExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *ExceptionClass) GetProperties() map[string]data.Property {
	return nil
}

func (s *ExceptionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return s.exception, true
	case "error":
		return s.error, true
	case "getMessage":
		return s.getMessage, true

	}
	return nil, false
}

func (s *ExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
	}
}

func (s *ExceptionClass) GetConstruct() data.Method {
	return s.exception
}
