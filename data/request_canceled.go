package data

// ErrRequestCanceled 表示 HTTP 请求截止时间已到，必须结束 PHP Handle。
// PHP try 的 recover 必须原样再 panic，不能当成普通 Exception。
var ErrRequestCanceled = RequestCanceled{}

type RequestCanceled struct{}

func (RequestCanceled) Error() string { return "origami request canceled" }

func IsRequestCanceled(v any) bool {
	_, ok := v.(RequestCanceled)
	return ok
}
