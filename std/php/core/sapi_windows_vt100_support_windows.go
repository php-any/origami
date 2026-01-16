//go:build windows
// +build windows

package core

import (
	"os"

	"github.com/php-any/origami/data"
	"golang.org/x/sys/windows"
)

// callWindowsImpl 在 Windows 上调用 Windows API 的实现
func (f *SapiWindowsVt100SupportFunction) callWindowsImpl(file *os.File, enableValue data.Value) (data.GetValue, data.Control) {
	fd := file.Fd()

	// 转换为 Windows 句柄
	handle := windows.Handle(fd)

	// Windows API 常量
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x0004

	// 获取当前控制台模式
	var mode uint32
	err := windows.GetConsoleMode(handle, &mode)
	if err != nil {
		// 如果获取失败（例如不是控制台流），返回 false
		return data.NewBoolValue(false), nil
	}

	if enableValue == nil {
		// 如果没有提供 enable 参数，返回当前状态
		enabled := (mode & ENABLE_VIRTUAL_TERMINAL_PROCESSING) != 0
		return data.NewBoolValue(enabled), nil
	}

	// 如果提供了 enable 参数，设置状态
	var enable bool
	if boolVal, ok := enableValue.(*data.BoolValue); ok {
		enable = boolVal.Value
	} else {
		// 参数类型错误，返回 false
		return data.NewBoolValue(false), nil
	}

	// 设置新的控制台模式
	var newMode uint32
	if enable {
		newMode = mode | ENABLE_VIRTUAL_TERMINAL_PROCESSING
	} else {
		newMode = mode &^ ENABLE_VIRTUAL_TERMINAL_PROCESSING
	}

	err = windows.SetConsoleMode(handle, newMode)
	if err != nil {
		// 设置失败，返回 false
		return data.NewBoolValue(false), nil
	}

	// 设置成功，返回 true
	return data.NewBoolValue(true), nil
}
