package std

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIncludeFunction() data.FuncStmt {
	return &IncludeFunction{}
}

type IncludeFunction struct{}

func (f *IncludeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取文件路径参数
	pathValue, _ := ctx.GetIndexValue(0)

	var filePath string
	switch p := pathValue.(type) {
	case data.AsString:
		filePath = p.AsString()
	default:
		return data.NewBoolValue(false), nil
	}

	// 检查路径是否为空
	if filePath == "" {
		return data.NewBoolValue(false), nil
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, data.NewErrorThrow(nil, errors.New(fmt.Sprintf("include 文件失败: %s, 错误: %v\n", filePath, err)))
	}

	// 如果路径是相对路径，则基于当前目录解析
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(currentDir, filePath)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，返回 false
		return nil, data.NewErrorThrow(nil, errors.New(fmt.Sprintf("include 文件失败: %s, 错误: %v\n", filePath, err)))
	}

	// 检查文件是否为目录
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if fileInfo.IsDir() {
		// 如果是目录，返回 false
		return nil, data.NewErrorThrow(nil, errors.New(fmt.Sprintf("include 文件失败: %s, 错误: 无法引入目录", filePath)))
	}

	// 获取 VM 实例
	vm := ctx.GetVM()

	// 尝试加载并执行文件
	v, acl := vm.LoadAndRun(filePath)
	if acl != nil {
		// 如果加载失败，返回 false
		return nil, acl
	}

	return v, nil
}

func (f *IncludeFunction) GetName() string {
	return "include"
}

func (f *IncludeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *IncludeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
