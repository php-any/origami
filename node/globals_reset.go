package node

import "github.com/php-any/origami/data"

// SuperglobalArrayProvider 由请求级 / VM 级容器实现，避免 $GLOBALS / $_SESSION 进程级串请求。
type SuperglobalArrayProvider interface {
	EnsureGlobalsArray() *data.ArrayValue
	EnsureSessionArray() *data.ArrayValue
}

func globalsArrayFromContext(ctx data.Context) *data.ArrayValue {
	if ctx != nil {
		if p, ok := ctx.GetVM().(SuperglobalArrayProvider); ok {
			return p.EnsureGlobalsArray()
		}
	}
	if globalsValue == nil {
		globalsValue = data.NewArrayValueFromSlots(nil)
	}
	return globalsValue
}

func sessionArrayFromContext(ctx data.Context) *data.ArrayValue {
	if ctx != nil {
		if p, ok := ctx.GetVM().(SuperglobalArrayProvider); ok {
			return p.EnsureSessionArray()
		}
	}
	if sessionValue == nil {
		sessionValue = data.NewArrayValueFromSlots(nil)
	}
	return sessionValue
}

// ResetSuperglobals 清空包级回退缓存。
// 使用 fpm.RequestVM / 已实现 SuperglobalArrayProvider 的 VM 时，请求态不依赖这些包级变量。
func ResetSuperglobals() {
	getValue = nil
	postValue = nil
	requestValue = nil
	cookieValue = nil
	sessionValue = nil
	filesValue = nil
	globalsValue = nil
	argvValue = nil
	argvInitialized = false
	argcValue = nil
}
