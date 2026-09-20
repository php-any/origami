package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewErrorExceptionClass() *ErrorExceptionClass {
	source := &Exception{}

	return &ErrorExceptionClass{
		construct:        &ErrorExceptionConstructMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getCode:          &ExceptionGetCodeMethod{source},
		getPrevious:      &ExceptionGetPreviousMethod{source},
		getSeverity:      &ExceptionGetSeverityMethod{source},
		getTrace:         &ExceptionGetTraceMethod{source},
		getFile:          &ExceptionGetFileMethod{source},
		getLine:          &ExceptionGetLineMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
		propSeverity:     node.NewProperty(nil, "severity", "protected", false, data.NewIntValue(1)),
	}
}

type ErrorExceptionClass struct {
	node.Node
	construct        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getCode          *ExceptionGetCodeMethod
	getPrevious      *ExceptionGetPreviousMethod
	getSeverity      *ExceptionGetSeverityMethod
	getTrace         *ExceptionGetTraceMethod
	getFile          *ExceptionGetFileMethod
	getLine          *ExceptionGetLineMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
	propSeverity     data.Property
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

func (s *ErrorExceptionClass) GetProperty(name string) (data.Property, bool) {
	if name == "severity" {
		return s.propSeverity, true
	}
	return nil, false
}

func (s *ErrorExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{s.propSeverity}
}

func (s *ErrorExceptionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return s.construct, true
	case "error":
		return s.error, true
	case "getMessage":
		return s.getMessage, true
	case "getCode":
		return s.getCode, true
	case "getPrevious":
		return s.getPrevious, true
	case "getSeverity":
		return s.getSeverity, true
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

func (s *ErrorExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
		s.getCode,
		s.getPrevious,
		s.getSeverity,
		s.getTrace,
		s.getFile,
		s.getLine,
		s.getTraceAsString,
	}
}

func (s *ErrorExceptionClass) GetConstruct() data.Method {
	return s.construct
}
