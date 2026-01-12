package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewLogicExceptionClass() *LogicExceptionClass {
	source := &Exception{}

	return &LogicExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

type LogicExceptionClass struct {
	node.Node
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func (s *LogicExceptionClass) AsString() string {
	return s.getMessage.source.msg
}

func (s *LogicExceptionClass) IsThrow() bool {
	//TODO implement me
	panic("implement me")
}

func (s *LogicExceptionClass) GetError() *data.Error {
	//TODO implement me
	panic("implement me")
}

func (s *LogicExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *LogicExceptionClass) GetName() string {
	return "LogicException"
}

func (s *LogicExceptionClass) GetExtend() *string {
	extend := "Exception"
	return &extend
}

func (s *LogicExceptionClass) GetImplements() []string {
	return nil
}

func (s *LogicExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *LogicExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *LogicExceptionClass) GetMethod(name string) (data.Method, bool) {
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

func (s *LogicExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
		s.getTraceAsString,
	}
}

func (s *LogicExceptionClass) GetConstruct() data.Method {
	return s.exception
}
