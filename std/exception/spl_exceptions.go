package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type simpleExceptionClass struct {
	node.Node
	name             string
	extend           string
	source           *Exception
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func newSimpleExceptionClass(name, extend string) *simpleExceptionClass {
	source := &Exception{}
	return &simpleExceptionClass{
		name:             name,
		extend:           extend,
		source:           source,
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

func (s *simpleExceptionClass) AsString() string      { return s.getMessage.source.msg }
func (s *simpleExceptionClass) IsThrow() bool         { return true }
func (s *simpleExceptionClass) GetError() *data.Error { return nil }
func (s *simpleExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}
func (s *simpleExceptionClass) GetName() string                            { return s.name }
func (s *simpleExceptionClass) GetExtend() *string                         { return &s.extend }
func (s *simpleExceptionClass) GetImplements() []string                    { return nil }
func (s *simpleExceptionClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *simpleExceptionClass) GetPropertyList() []data.Property           { return nil }
func (s *simpleExceptionClass) GetConstruct() data.Method                  { return s.exception }
func (s *simpleExceptionClass) GetMethod(name string) (data.Method, bool) {
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
func (s *simpleExceptionClass) GetMethods() []data.Method {
	return []data.Method{s.error, s.getMessage, s.getTraceAsString}
}

func NewUnexpectedValueExceptionClass() data.ClassStmt {
	return newSimpleExceptionClass("UnexpectedValueException", "RuntimeException")
}
func NewUnderflowExceptionClass() data.ClassStmt {
	return newSimpleExceptionClass("UnderflowException", "RuntimeException")
}
func NewOverflowExceptionClass() data.ClassStmt {
	return newSimpleExceptionClass("OverflowException", "RuntimeException")
}
func NewOutOfBoundsExceptionClass() data.ClassStmt {
	return newSimpleExceptionClass("OutOfBoundsException", "RuntimeException")
}
func NewOutOfRangeExceptionClass() data.ClassStmt {
	return newSimpleExceptionClass("OutOfRangeException", "RuntimeException")
}
func NewBadFunctionCallExceptionClass() data.ClassStmt {
	return newSimpleExceptionClass("BadFunctionCallException", "LogicException")
}
func NewDomainExceptionClass() data.ClassStmt {
	return newSimpleExceptionClass("DomainException", "LogicException")
}
