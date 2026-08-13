package php

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/std/php/stream"
	"github.com/php-any/origami/utils"
)

func NewFileGetContentsFunction() data.FuncStmt {
	return &FileGetContentsFunction{}
}

type FileGetContentsFunction struct{}

func (f *FileGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取文件路径参数
	pathValue, _ := ctx.GetIndexValue(0)

	if pathValue == nil {
		return nil, utils.NewThrowf("FileGetContentsFunction called with no file path")
	}

	var filePath string
	// 检查是否是字符串类型
	if s, ok := pathValue.(data.AsString); ok {
		filePath = s.AsString()
	} else {
		// 如果不是字符串类型（可能是资源），尝试转换为字符串
		filePath = pathValue.AsString()
	}

	if filePath == "" {
		return nil, utils.NewThrowf("FileGetContentsFunction called with no file path")
	}

	if filePath == "php://input" || strings.HasPrefix(filePath, "php://input") {
		return data.NewStringValue(core.PhptInputBody()), nil
	}

	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		contextVal, _ := ctx.GetIndexValue(2)
		streamCtx := stream.ContextFromResource(contextVal)
		content, ok := stream.HTTPGetContents(filePath, streamCtx)
		if !ok {
			return data.NewBoolValue(false), nil
		}
		return data.NewStringValue(content), nil
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, utils.NewThrowf("FileGetContentsFunction called with file path '%s': %v", filePath, err)
	}
	return data.NewStringValue(string(bytes)), nil
}

func (f *FileGetContentsFunction) GetName() string { return "file_get_contents" }

func (f *FileGetContentsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
		node.NewParameter(nil, "use_include_path", 1, data.NewBoolValue(false), nil),
		node.NewParameter(nil, "context", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *FileGetContentsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "use_include_path", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "context", 2, data.NewNullableType(data.NewBaseType("resource"))),
	}
}
