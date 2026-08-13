package php

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewFilePutContentsFunction() data.FuncStmt {
	return &FilePutContentsFunction{}
}

type FilePutContentsFunction struct{}

func (f *FilePutContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取文件路径参数
	pathValue, _ := ctx.GetIndexValue(0)
	dataValue, _ := ctx.GetIndexValue(1)

	// 检查第一个参数是否是资源类型
	if pathValue == nil {
		return nil, utils.NewThrowf("FilePutContentsFunction called with no file path")
	}

	var filePath string
	// 检查是否是字符串类型
	if s, ok := pathValue.(data.AsString); ok {
		filePath = s.AsString()
	} else {
		// 如果不是字符串类型（可能是资源），尝试转换为字符串
		// 如果是资源类型，应该返回错误或使用资源的路径
		filePath = pathValue.AsString()
	}

	if filePath == "" {
		return nil, utils.NewThrowf("FilePutContentsFunction called with no file path")
	}

	var content string
	if s, ok := dataValue.(data.AsString); ok {
		content = s.AsString()
	} else {
		content = dataValue.AsString()
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return nil, utils.NewThrowf("FilePutContentsFunction called with file path '%s': %v", filePath, err)
	}
	return data.NewIntValue(len(content)), nil
}

func (f *FilePutContentsFunction) GetName() string { return "file_put_contents" }

func (f *FilePutContentsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, nil),
	}
}

func (f *FilePutContentsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "data", 1, data.NewBaseType("string")),
	}
}
