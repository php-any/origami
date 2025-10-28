package php

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewFileGetContentsFunction() data.FuncStmt {
	return &FileGetContentsFunction{}
}

type FileGetContentsFunction struct{}

func (f *FileGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取文件路径参数
	pathValue, _ := ctx.GetIndexValue(0)

	var filePath string
	switch p := pathValue.(type) {
	case data.AsString:
		filePath = p.AsString()
	default:
		return data.NewStringValue(""), nil
	}

	if filePath == "" {
		return data.NewStringValue(""), nil
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(string(bytes)), nil
}

func (f *FileGetContentsFunction) GetName() string { return "file_get_contents" }

func (f *FileGetContentsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *FileGetContentsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
