package php

import (
	"sync"

	"github.com/php-any/origami/data"
)

// PHP error_get_last / error_clear_last 用的最近错误状态。
type lastErrorInfo struct {
	Type    int
	Message string
	File    string
	Line    int
}

var (
	lastErrorMu sync.Mutex
	lastError   *lastErrorInfo
)

func setLastError(typ int, message, file string, line int) {
	lastErrorMu.Lock()
	defer lastErrorMu.Unlock()
	lastError = &lastErrorInfo{
		Type:    typ,
		Message: message,
		File:    file,
		Line:    line,
	}
}

func clearLastError() {
	lastErrorMu.Lock()
	defer lastErrorMu.Unlock()
	lastError = nil
}

func getLastError() *lastErrorInfo {
	lastErrorMu.Lock()
	defer lastErrorMu.Unlock()
	if lastError == nil {
		return nil
	}
	cp := *lastError
	return &cp
}

func lastErrorToArray(info *lastErrorInfo) data.Value {
	return &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("type", data.NewIntValue(info.Type)),
			data.NewNamedZVal("message", data.NewStringValue(info.Message)),
			data.NewNamedZVal("file", data.NewStringValue(info.File)),
			data.NewNamedZVal("line", data.NewIntValue(info.Line)),
		},
	}
}
