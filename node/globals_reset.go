package node

import "github.com/php-any/origami/data"

// SuperglobalArrayProvider 由请求级 / VM 级容器实现，避免 $GLOBALS / $_SESSION 进程级串请求。
type SuperglobalArrayProvider interface {
	EnsureGlobalsArray() *data.ObjectValue
	EnsureSessionArray() *data.ObjectValue
}

func globalsArrayFromContext(ctx data.Context) *data.ObjectValue {
	if ctx != nil {
		if p, ok := ctx.GetVM().(SuperglobalArrayProvider); ok {
			return p.EnsureGlobalsArray()
		}
	}
	if globalsValue == nil {
		globalsValue = data.NewObjectValue()
	}
	return globalsValue
}

func sessionArrayFromContext(ctx data.Context) *data.ObjectValue {
	if ctx != nil {
		if p, ok := ctx.GetVM().(SuperglobalArrayProvider); ok {
			return p.EnsureSessionArray()
		}
	}
	if sessionValue == nil {
		sessionValue = data.NewObjectValue()
	}
	return sessionValue
}

// ResetSuperglobals 清空包级回退缓存。
// 使用 fpm.RequestVM / 已实现 SuperglobalArrayProvider 的 VM 时，请求态不依赖这些包级变量。
func ResetSuperglobals() {
	getValue = nil
	postValue = nil
	serverValue = nil
	requestValue = nil
	cookieValue = nil
	sessionValue = nil
	filesValue = nil
	envValue = nil
	globalsValue = nil
}
