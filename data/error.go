package data

// Error 表示节点执行过程中的错误
type Error struct {
	From     From     // 错误发生的来源
	message  string   // 错误信息
	cause    error    // 原始错误
	children []*Error // 子错误
}

// NewError 创建一个新的错误
func NewError(from From, message string, cause error) *Error {
	return &Error{
		From:     from,
		message:  message,
		cause:    cause,
		children: make([]*Error, 0),
	}
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return e.message
}

// GetFrom 获取错误来源
func (e *Error) GetFrom() From {
	return e.From
}

// GetMessage 获取错误信息
func (e *Error) GetMessage() string {
	return e.message
}

// GetCause 获取原始错误
func (e *Error) GetCause() error {
	return e.cause
}

// GetChildren 获取子错误
func (e *Error) GetChildren() []*Error {
	return e.children
}

// AddChild 添加子错误
func (e *Error) AddChild(child *Error) {
	e.children = append(e.children, child)
}
