package file

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsReadableFunction 实现 is_readable 函数
// 检查文件是否可读
type IsReadableFunction struct{}

func NewIsReadableFunction() data.FuncStmt {
	return &IsReadableFunction{}
}

func (f *IsReadableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 检查文件是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		// 如果文件不存在或其他错误，返回 false
		return data.NewBoolValue(false), nil
	}

	// 检查是否为目录（目录不能直接读取）
	if fileInfo.IsDir() {
		return data.NewBoolValue(false), nil
	}

	// 检查文件是否可读
	file, err := os.Open(path)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	file.Close()

	return data.NewBoolValue(true), nil
}

func (f *IsReadableFunction) GetName() string {
	return "is_readable"
}

func (f *IsReadableFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *IsReadableFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
