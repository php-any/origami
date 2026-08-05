package exception

import (
	"runtime"

	"github.com/php-any/origami/data"
)

type traceFrame struct {
	file     string
	line     int
	function string
}

type Exception struct {
	msg   string
	file  string
	line  int
	trace []traceFrame
}

func (e *Exception) Exception(msg string) {
	e.msg = msg
	e.captureTrace()
}

func (e *Exception) captureTrace() {
	const depth = 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	e.trace = nil
	for {
		frame, more := frames.Next()
		if frame.File == "" && !more {
			break
		}
		if frame.File != "" {
			if e.file == "" {
				e.file = frame.File
				e.line = frame.Line
			}
			e.trace = append(e.trace, traceFrame{
				file:     frame.File,
				line:     frame.Line,
				function: frame.Function,
			})
		}
		if !more {
			break
		}
	}
}

func (e *Exception) Error() string {
	return ""
}

func (e *Exception) GetMessage() string {
	return e.msg
}

func (e *Exception) GetFile() string {
	return e.file
}

func (e *Exception) GetLine() int {
	return e.line
}

func (e *Exception) GetTraceAsString() string {
	return "Stack trace:\n  at Exception.constructor()\n  at main()"
}

func (e *Exception) GetTraceValues() []data.Value {
	frames := make([]data.Value, 0, len(e.trace))
	for _, frame := range e.trace {
		obj := data.NewObjectValue()
		obj.SetProperty("file", data.NewStringValue(frame.file))
		obj.SetProperty("line", data.NewIntValue(frame.line))
		if frame.function != "" {
			obj.SetProperty("function", data.NewStringValue(frame.function))
		}
		frames = append(frames, obj)
	}
	return frames
}

// ExceptionMethods 封装一个 Exception 实例及其关联方法，方便其他包复用。
type ExceptionMethods struct {
	ConstructMethod   data.Method
	ErrorMethod       data.Method
	GetMessageMethod  data.Method
	GetCodeMethod     data.Method
	GetPreviousMethod data.Method
	GetTraceMethod    data.Method
	GetFileMethod     data.Method
	GetLineMethod     data.Method
	GetTraceAsString  data.Method
}

// SetMessage 设置内部 Exception 的消息，供外部包创建异常时使用。
func (m *ExceptionMethods) SetMessage(msg string) {
	if construct, ok := m.ConstructMethod.(*ExceptionExceptionMethod); ok {
		construct.source.msg = msg
	}
}

// NewExceptionMethods 创建一组封装好的 Exception 方法，供其他包继承使用。
func NewExceptionMethods() ExceptionMethods {
	src := &Exception{}
	return ExceptionMethods{
		ConstructMethod:   &ExceptionExceptionMethod{source: src},
		ErrorMethod:       &ExceptionErrorMethod{source: src},
		GetMessageMethod:  &ExceptionGetMessageMethod{source: src},
		GetCodeMethod:     &ExceptionGetCodeMethod{source: src},
		GetPreviousMethod: &ExceptionGetPreviousMethod{source: src},
		GetTraceMethod:    &ExceptionGetTraceMethod{source: src},
		GetFileMethod:     &ExceptionGetFileMethod{source: src},
		GetLineMethod:     &ExceptionGetLineMethod{source: src},
		GetTraceAsString:  &ExceptionGetTraceAsStringMethod{source: src},
	}
}
