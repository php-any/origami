package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ReflectionExceptionClass 实现 ReflectionException 类
// ReflectionException 继承自 Exception，用于反射相关的错误
type ReflectionExceptionClass struct {
	node.Node
	exception        data.Method
	getMessage       *ExceptionGetMessageMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
}

func NewReflectionExceptionClass() *ReflectionExceptionClass {
	source := &Exception{}

	return &ReflectionExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
	}
}

func (s *ReflectionExceptionClass) AsString() string {
	return s.getMessage.source.msg
}

func (s *ReflectionExceptionClass) IsThrow() bool {
	return true
}

func (s *ReflectionExceptionClass) GetError() *data.Error {
	// 返回 nil，实际的错误信息存储在 Exception 对象中
	return nil
}

func (s *ReflectionExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *ReflectionExceptionClass) GetName() string {
	return "ReflectionException"
}

func (s *ReflectionExceptionClass) GetExtend() *string {
	extend := "Exception"
	return &extend
}

func (s *ReflectionExceptionClass) GetImplements() []string {
	// ReflectionException 实现 Throwable 接口（通过继承 Exception）
	return []string{"Throwable"}
}

func (s *ReflectionExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *ReflectionExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *ReflectionExceptionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return s.exception, true
	case "getMessage":
		return s.getMessage, true
	case "getTraceAsString":
		return s.getTraceAsString, true
	}
	return nil, false
}

func (s *ReflectionExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.exception,
		s.getMessage,
		s.getTraceAsString,
	}
}

func (s *ReflectionExceptionClass) GetConstruct() data.Method {
	return s.exception
}
