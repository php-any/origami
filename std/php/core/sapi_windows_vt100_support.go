package core

import (
	"os"
	"reflect"
	"runtime"

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

	// 调用 Windows 特定的实现
	return f.callWindowsImpl(file, enableValue)
}

// callWindowsImpl 在 Windows 上调用 Windows API
// 这个函数的具体实现使用条件编译（见 sapi_windows_vt100_support_windows.go）
func (f *SapiWindowsVt100SupportFunction) callWindowsImpl(file *os.File, enableValue data.Value) (data.GetValue, data.Control) {
	// 默认实现（非 Windows 系统不会到达这里，但为了编译通过）
	return data.NewBoolValue(false), nil
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
