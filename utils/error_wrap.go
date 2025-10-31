package utils

import (
	"fmt"
	"runtime"

	"github.com/php-any/origami/data"
)

// NewThrow 基于调用点构建可跳转的 From，并返回带位置信息的错误控制
// 期望来源是调用 NewThrow 的上上层（跳过 NewThrowWithSkip 与 NewThrow 自身）
func NewThrow(err error) data.Control { return NewThrowWithSkip(2, err) }

// NewThrowf 格式化错误并附带调用点位置
func NewThrowf(format string, args ...any) data.Control {
	return NewThrowWithSkip(2, fmt.Errorf(format, args...))
}

// NewThrowWithSkip 允许自定义 runtime.Caller 的 skip 层级
func NewThrowWithSkip(skip int, err error) data.Control {
	if err == nil {
		return nil
	}
	// 获取调用位置
	// runtime.Caller(skip): 0=本函数, 1=调用者(NewThrowWithSkip的调用者), 2=再上一层...
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return data.NewErrorThrow(nil, err)
	}

	// 组装 From：仅用于点击跳转，偏移量可为 0，列设置为 0
	from := data.NewBaseFromWithPosition(file, 0, 0, line-1, 0, line-1, 0)
	return data.NewErrorThrow(from, err)
}
