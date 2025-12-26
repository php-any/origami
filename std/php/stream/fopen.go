package stream

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// FopenFunction 实现 fopen 函数
type FopenFunction struct{}

func NewFopenFunction() data.FuncStmt {
	return &FopenFunction{}
}

func (f *FopenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取文件名参数
	filenameValue, _ := ctx.GetIndexValue(0)
	if filenameValue == nil {
		return data.NewBoolValue(false), nil
	}

	var filename string
	if s, ok := filenameValue.(data.AsString); ok {
		filename = s.AsString()
	} else {
		filename = filenameValue.AsString()
	}

	if filename == "" {
		return data.NewBoolValue(false), nil
	}

	// 获取模式参数
	modeValue, _ := ctx.GetIndexValue(1)
	mode := "r" // 默认只读模式
	if modeValue != nil {
		if s, ok := modeValue.(data.AsString); ok {
			mode = s.AsString()
		} else {
			mode = modeValue.AsString()
		}
	}
	if mode == "" {
		mode = "r"
	}

	// 打开文件
	file, err := os.OpenFile(filename, parseMode(mode), 0644)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// 创建流信息
	streamInfo := NewStreamInfo(file, mode)

	// 创建流资源类，使用文件描述符作为资源ID
	fd := int(file.Fd())
	resourceClass := core.NewResourceClass("stream", streamInfo, fd)

	// 创建流资源对象
	streamResource := core.NewResourceValue(resourceClass, ctx)

	return streamResource, nil
}

// parseMode 解析 PHP 文件打开模式为 Go 的 os.OpenFile 标志
func parseMode(mode string) int {
	flags := 0
	switch mode[0] {
	case 'r':
		flags = os.O_RDONLY
	case 'w':
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case 'a':
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case 'x':
		flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	case 'c':
		flags = os.O_WRONLY | os.O_CREATE
	}

	// 检查是否有 '+' 表示读写模式
	if len(mode) > 1 && mode[1] == '+' {
		switch mode[0] {
		case 'r':
			flags = os.O_RDWR
		case 'w':
			flags = os.O_RDWR | os.O_CREATE | os.O_TRUNC
		case 'a':
			flags = os.O_RDWR | os.O_CREATE | os.O_APPEND
		case 'x':
			flags = os.O_RDWR | os.O_CREATE | os.O_EXCL
		case 'c':
			flags = os.O_RDWR | os.O_CREATE
		}
	}

	return flags
}

func (f *FopenFunction) GetName() string {
	return "fopen"
}

func (f *FopenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
		node.NewParameter(nil, "mode", 1, node.NewStringLiteral(nil, "r"), nil),
	}
}

func (f *FopenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "mode", 1, data.NewBaseType("string")),
	}
}
