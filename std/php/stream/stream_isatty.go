package stream

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// StreamIsattyFunction 实现 stream_isatty 函数
// stream_isatty(resource $stream): bool
// 检查流是否是终端（TTY）
type StreamIsattyFunction struct{}

func NewStreamIsattyFunction() data.FuncStmt {
	return &StreamIsattyFunction{}
}

func (f *StreamIsattyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取流资源
	streamValue, _ := ctx.GetIndexValue(0)
	if streamValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 从资源对象中获取 StreamInfo
	var streamInfo *StreamInfo
	if res, ok := streamValue.(*core.ResourceValue); ok {
		resource := res.GetResource()
		if resource == nil {
			return data.NewBoolValue(false), nil
		}
		if info, ok := resource.(*StreamInfo); ok {
			streamInfo = info
		} else {
			// 不是 StreamInfo 类型，返回 false
			return data.NewBoolValue(false), nil
		}
	} else {
		// 不是 ResourceValue 类型
		return data.NewBoolValue(false), nil
	}

	// 检查流是否已关闭
	if streamInfo.IsClosed() {
		return data.NewBoolValue(false), nil
	}

	// 检查文件是否为 nil
	if streamInfo.File == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查是否是 TTY
	// 使用 os.File 的 Stat() 方法结合文件模式来判断
	fileInfo, err := streamInfo.File.Stat()
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// 检查文件模式是否是字符设备（TTY 通常是字符设备）
	mode := fileInfo.Mode()
	isCharDevice := (mode & os.ModeCharDevice) != 0

	// 如果文件是字符设备，很可能是 TTY
	if isCharDevice {
		return data.NewBoolValue(true), nil
	}

	// 对于标准输入/输出/错误流，检查文件描述符
	// 如果文件模式是普通文件，则不是 TTY
	if mode.IsRegular() {
		return data.NewBoolValue(false), nil
	}

	// 其他情况（管道、设备等）可能不是 TTY
	// 但为了简化，我们主要依赖字符设备检查
	return data.NewBoolValue(false), nil
}

func (f *StreamIsattyFunction) GetName() string {
	return "stream_isatty"
}

func (f *StreamIsattyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
	}
}

func (f *StreamIsattyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
	}
}
