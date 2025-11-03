package php

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIsFileFunction() data.FuncStmt {
	return &IsFileFunction{}
}

type IsFileFunction struct{}

func (f *IsFileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 检查是否为文件（不是目录）
	isFile := !fileInfo.IsDir()
	return data.NewBoolValue(isFile), nil
}

func (f *IsFileFunction) GetName() string {
	return "is_file"
}

func (f *IsFileFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, nil),
	}
}

func (f *IsFileFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, data.NewBaseType("string")),
	}
}
