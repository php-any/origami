package exception

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
