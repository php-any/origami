package validation

// Violation 单条校验失败信息。
type Violation struct {
	Field   string
	Message string
}
