package file

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FileExistsFunction 实现 file_exists 函数
// 检查文件或目录是否存在
type FileExistsFunction struct{}

func NewFileExistsFunction() data.FuncStmt {
	return &FileExistsFunction{}
}

func (f *FileExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 检查文件或目录是否存在
	_, err := os.Stat(path)
	if err != nil {
		// 如果文件不存在或其他错误，返回 false
		if os.IsNotExist(err) {
			return data.NewBoolValue(false), nil
		}
		// 其他错误也返回 false
		return data.NewBoolValue(false), nil
	}

	// 文件或目录存在
	return data.NewBoolValue(true), nil
}

func (f *FileExistsFunction) GetName() string {
	return "file_exists"
}

func (f *FileExistsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *FileExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
