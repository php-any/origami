package php

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIsDirFunction() data.FuncStmt {
	return &IsDirFunction{}
}

type IsDirFunction struct{}

func (f *IsDirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取路径参数
	pathValue, _ := ctx.GetIndexValue(0)

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		return data.NewBoolValue(false), nil
	}

	// 检查路径是否为空
	if path == "" {
		return data.NewBoolValue(false), nil
	}

	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		// 如果文件不存在或其他错误，返回 false
		return data.NewBoolValue(false), nil
	}

	// 检查是否为目录
	isDir := fileInfo.IsDir()
	return data.NewBoolValue(isDir), nil
}
func (f *IsDirFunction) GetName() string {
	return "is_dir"
}

func (f *IsDirFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, nil),
	}
}

func (f *IsDirFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
	}
}
