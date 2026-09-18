package data

// CallFrame 表示 PHP 调用栈中的一帧（供 debug_backtrace）
type CallFrame struct {
	File     string
	Line     int
	Function string
	Class    string
	Type     string // "->" 或 "::"
}

// CallStackTracker 可选接口：VM 实现后 debug_backtrace 可返回真实 PHP 调用栈
type CallStackTracker interface {
	PushCallFrame(frame CallFrame)
	PopCallFrame()
	SnapshotCallStack() []CallFrame
}

// CallRecorder 请求内调用深度 + 栈。优先由 Context 实现，避免共享主 VM 串请求。
type CallRecorder interface {
	EnterCall() int
	LeaveCall()
	CallStackTracker
}

// ContextEscaper 标记调用帧被生成器等长期持有，禁止归还 Context 池。
type ContextEscaper interface {
	MarkEscaped()
	IsEscaped() bool
}
