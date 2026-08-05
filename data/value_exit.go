package data

// NewExitControl 创建程序退出控制流（exit/die）。
func NewExitControl(code int) ExitControl {
	return &ExitValue{code: code}
}

// ExitValue 实现 ExitControl。
type ExitValue struct {
	code int
}

func (e *ExitValue) GetValue(ctx Context) (GetValue, Control) {
	return NewNullValue(), nil
}

func (e *ExitValue) AsString() string {
	return "exit"
}

func (e *ExitValue) IsExit() bool { return true }

func (e *ExitValue) GetCode() int { return e.code }
