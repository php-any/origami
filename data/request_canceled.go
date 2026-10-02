package data

import (
	"context"
	"time"
)

// ErrRequestCanceled 表示 HTTP 请求截止时间已到，必须结束 PHP Handle。
// PHP try 的 recover 必须原样再 panic，不能当成普通 Exception。
var ErrRequestCanceled = RequestCanceled{}

type RequestCanceled struct{}

func (RequestCanceled) Error() string { return "origami request canceled" }

func IsRequestCanceled(v any) bool {
	_, ok := v.(RequestCanceled)
	return ok
}

// CheckRequest is used at blocking I/O boundaries. Cancellation must not become
// a PHP false/empty result that lets the script continue after disconnect.
func CheckRequest(ctx context.Context) {
	if ctx.Err() != nil {
		panic(ErrRequestCanceled)
	}
}

func WaitRequest(ctx context.Context, duration time.Duration) {
	CheckRequest(ctx)
	if ctx.Done() == nil {
		time.Sleep(duration)
		return
	}
	if duration <= 0 {
		return
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		CheckRequest(ctx)
	case <-ctx.Done():
		panic(ErrRequestCanceled)
	}
}
