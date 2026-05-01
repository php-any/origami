package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewBadMethodCallExceptionClass() *BadMethodCallExceptionClass {
	source := &Exception{}

	return &BadMethodCallExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

type BadMethodCallExceptionClass struct {
	node.Node
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func (s *BadMethodCallExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *BadMethodCallExceptionClass) GetName() string {
	return "BadMethodCallException"
}

func (s *BadMethodCallExceptionClass) GetExtend() *string {
	extend := "LogicException"
	return &extend
}

func (s *BadMethodCallExceptionClass) GetImplements() []string { return nil }

func (s *BadMethodCallExceptionClass) GetProperty(_ string) (data.Property, bool) { return nil, false }

func (s *BadMethodCallExceptionClass) GetPropertyList() []data.Property { return []data.Property{} }

func (s *BadMethodCallExceptionClass) GetMethod(name string) (data.Method, bool) {
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

func (s *BadMethodCallExceptionClass) GetMethods() []data.Method {
	return []data.Method{s.error, s.getMessage, s.getTraceAsString}
}

func (s *BadMethodCallExceptionClass) GetConstruct() data.Method { return s.exception }
