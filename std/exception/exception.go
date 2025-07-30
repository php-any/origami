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
