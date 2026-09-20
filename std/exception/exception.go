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
	if e == nil {
		return
	}
	file, line, frames := captureGoTrace(3)
	e.msg = msg
	e.file = file
	e.line = line
	e.trace = frames
}

// captureGoTrace 采集 Go 调用栈。必须先在局部切片上积累再返回，禁止对共享
// Exception.trace 边 nil 边 append：HTTP 并发下 Livewire throw_unless(new Exception)
// 会同时走进构造函数，竞态会 panic（invalid memory address）。
func captureGoTrace(skip int) (file string, line int, frames []traceFrame) {
	if skip < 1 {
		skip = 1
	}
	const depth = 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs)
	if n <= 0 {
		return "", 0, nil
	}
	it := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := it.Next()
		if frame.File == "" && !more {
			break
		}
		if frame.File != "" {
			if file == "" {
				file = frame.File
				line = frame.Line
			}
			frames = append(frames, traceFrame{
				file:     frame.File,
				line:     frame.Line,
				function: frame.Function,
			})
		}
		if !more {
			break
		}
	}
	return file, line, frames
}

func traceFramesToValues(frames []traceFrame) []data.Value {
	out := make([]data.Value, 0, len(frames))
	for _, frame := range frames {
		obj := data.NewObjectValue()
		obj.SetProperty("file", data.NewStringValue(frame.file))
		obj.SetProperty("line", data.NewIntValue(frame.line))
		if frame.function != "" {
			obj.SetProperty("function", data.NewStringValue(frame.function))
		}
		out = append(out, obj)
	}
	return out
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
	if e == nil {
		return nil
	}
	return traceFramesToValues(e.trace)
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
	if construct, ok := m.ConstructMethod.(*ExceptionExceptionMethod); ok && construct != nil && construct.source != nil {
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
