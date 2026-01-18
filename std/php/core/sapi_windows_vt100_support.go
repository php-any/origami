package core

import (
	"os"
	"reflect"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SapiWindowsVt100SupportFunction 实现 sapi_windows_vt100_support 函数
// sapi_windows_vt100_support(resource $stream, ?bool $enable = null): bool
// 检查或设置 Windows 控制台的 VT100 支持
type SapiWindowsVt100SupportFunction struct{}

func NewSapiWindowsVt100SupportFunction() data.FuncStmt {
	return &SapiWindowsVt100SupportFunction{}
}

func (f *SapiWindowsVt100SupportFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 如果不是 Windows 系统，直接返回 false
	if runtime.GOOS != "windows" {
		return data.NewBoolValue(false), nil
	}

	// 获取第一个参数：流资源
	streamValue, _ := ctx.GetIndexValue(0)
	if streamValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 从资源对象中获取文件
	var file *os.File
	if res, ok := streamValue.(*ResourceValue); ok {
		resource := res.GetResource()
		if resource == nil {
			return data.NewBoolValue(false), nil
		}

		// 使用反射来访问 StreamInfo 的 File 字段，避免导入循环
		rv := reflect.ValueOf(resource)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}

		// 检查是否有 File 字段
		fileField := rv.FieldByName("File")
		if !fileField.IsValid() || fileField.IsNil() {
			return data.NewBoolValue(false), nil
		}

		// 检查是否有 IsClosed 方法
		isClosedMethod := rv.MethodByName("IsClosed")
		if isClosedMethod.IsValid() {
			results := isClosedMethod.Call(nil)
			if len(results) > 0 && results[0].Bool() {
				return data.NewBoolValue(false), nil
			}
		}

		// 获取 File 字段的值
		if fileField.Kind() == reflect.Ptr {
			file = fileField.Interface().(*os.File)
		} else {
			return data.NewBoolValue(false), nil
		}
	} else {
		// 不是 ResourceValue 类型
		return data.NewBoolValue(false), nil
	}

	// 检查文件是否为 nil
	if file == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取第二个参数：enable（可选）
	enableValue, _ := ctx.GetIndexValue(1)

	// Windows 特定的实现
	return f.handleWindowsVT100(file, enableValue)
}

// handleWindowsVT100 处理 Windows VT100 支持
func (f *SapiWindowsVt100SupportFunction) handleWindowsVT100(file *os.File, enableValue data.Value) (data.GetValue, data.Control) {
	// Windows API 常量
	const (
		ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x0004
	)

	// 动态加载 kernel32.dll 中的函数
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	fd := file.Fd()
	handle := uintptr(fd)

	// 获取当前控制台模式
	var mode uint32
	ret, _, err := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		// GetConsoleMode 失败（例如不是控制台流），返回 false
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

	ret, _, err = setConsoleMode.Call(handle, uintptr(newMode))
	if ret == 0 {
		// SetConsoleMode 失败，返回 false
		_ = err // 忽略错误，因为我们已经检查了返回值
		return data.NewBoolValue(false), nil
	}

	// 设置成功，返回 true
	return data.NewBoolValue(true), nil
}

func (f *SapiWindowsVt100SupportFunction) GetName() string {
	return "sapi_windows_vt100_support"
}

func (f *SapiWindowsVt100SupportFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "enable", 1, nil, data.Bool{}),
	}
}

func (f *SapiWindowsVt100SupportFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
		node.NewVariable(nil, "enable", 1, data.Bool{}),
	}
}
