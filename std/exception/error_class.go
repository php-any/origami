package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// phpErrorClass 实现 PHP Error 及其子类（ValueError / TypeError 等）。
type phpErrorClass struct {
	node.Node
	name             string
	extend           string
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getTrace         *ExceptionGetTraceMethod
	getFile          *ExceptionGetFileMethod
	getLine          *ExceptionGetLineMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func newPHPErrorClass(name, extend string) *phpErrorClass {
	source := &Exception{}
	return &phpErrorClass{
		name:             name,
		extend:           extend,
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTrace:         &ExceptionGetTraceMethod{source},
		getFile:          &ExceptionGetFileMethod{source},
		getLine:          &ExceptionGetLineMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

func NewErrorClass() data.ClassStmt      { return newPHPErrorClass("Error", "") }
func NewValueErrorClass() data.ClassStmt { return newPHPErrorClass("ValueError", "Error") }
func NewTypeErrorClass() data.ClassStmt  { return newPHPErrorClass("TypeError", "Error") }
func NewArgumentCountErrorClass() data.ClassStmt {
	return newPHPErrorClass("ArgumentCountError", "TypeError")
}
func NewUnhandledMatchErrorClass() data.ClassStmt {
	return newPHPErrorClass("UnhandledMatchError", "Error")
}
func NewArithmeticErrorClass() data.ClassStmt { return newPHPErrorClass("ArithmeticError", "Error") }
func NewDivisionByZeroErrorClass() data.ClassStmt {
	return newPHPErrorClass("DivisionByZeroError", "ArithmeticError")
}

func (s *phpErrorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *phpErrorClass) GetName() string { return s.name }

func (s *phpErrorClass) GetExtend() *string {
	if s.extend == "" {
		return nil
	}
	e := s.extend
	return &e
}

func (s *phpErrorClass) GetImplements() []string {
	if s.extend == "" {
		return []string{"Throwable"}
	}
	return nil
}

func (s *phpErrorClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *phpErrorClass) GetPropertyList() []data.Property           { return []data.Property{} }

func (s *phpErrorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return s.exception, true
	case "error":
		return s.error, true
	case "getMessage":
		return s.getMessage, true
	case "getTrace":
		return s.getTrace, true
	case "getFile":
		return s.getFile, true
	case "getLine":
		return s.getLine, true
	case "getTraceAsString":
		return s.getTraceAsString, true
	}
	return nil, false
}

func (s *phpErrorClass) GetMethods() []data.Method {
	return []data.Method{s.error, s.getMessage, s.getTrace, s.getFile, s.getLine, s.getTraceAsString}
}

func (s *phpErrorClass) GetConstruct() data.Method { return s.exception }
