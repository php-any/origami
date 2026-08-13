package file

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FilesizeFunction 实现 filesize 函数
// 获取文件大小（字节数）
type FilesizeFunction struct{}

func NewFilesizeFunction() data.FuncStmt {
	return &FilesizeFunction{}
}

func (f *FilesizeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 检查是否为目录（目录没有大小概念）
	if fileInfo.IsDir() {
		return data.NewBoolValue(false), nil
	}

	// 返回文件大小
	return data.NewIntValue(int(fileInfo.Size())), nil
}

func (f *FilesizeFunction) GetName() string {
	return "filesize"
}

func (f *FilesizeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *FilesizeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
