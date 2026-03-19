package exception

import "github.com/php-any/origami/data"

type Exception struct {
	msg string
}

func (e *Exception) Exception(msg string) {
	e.msg = msg
}

func (e *Exception) Error() string {
	return ""
}

func (e *Exception) GetMessage() string {
	return e.msg
}

func (e *Exception) GetTraceAsString() string {
	// 简单的堆栈跟踪实现
	return "Stack trace:\n  at Exception.constructor()\n  at main()"
}

// ExceptionMethods 封装一个 Exception 实例及其关联方法，方便其他包复用。
type ExceptionMethods struct {
	ConstructMethod  data.Method
	ErrorMethod      data.Method
	GetMessageMethod data.Method
	GetTraceMethod   data.Method
}

// SetMessage 设置内部 Exception 的消息，供外部包创建异常时使用。
func (m *ExceptionMethods) SetMessage(msg string) {
	// 通过 ConstructMethod 里面的 source 设置消息
	if construct, ok := m.ConstructMethod.(*ExceptionExceptionMethod); ok {
		construct.source.msg = msg
	}
}

// NewExceptionMethods 创建一组封装好的 Exception 方法，供其他包继承使用。
func NewExceptionMethods() ExceptionMethods {
	src := &Exception{}
	return ExceptionMethods{
		ConstructMethod:  &ExceptionExceptionMethod{source: src},
		ErrorMethod:      &ExceptionErrorMethod{source: src},
		GetMessageMethod: &ExceptionGetMessageMethod{source: src},
		GetTraceMethod:   &ExceptionGetTraceAsStringMethod{source: src},
	}
}
